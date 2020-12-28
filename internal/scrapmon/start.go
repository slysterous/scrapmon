package scrapmon

import (
	"context"
	"fmt"
	customNumber "github.com/slysterous/custom-number"
	"sync"
	"time"
)

// imgur.com/abcdef.png
// StartCommand is what happens when the command is executed.
func (cm CommandManager) StartCommand(fromCode string, iterations int, workerNumber int) error {
	// mark the start time
	start := time.Now()

	// if no code was provided, then we resume from the last created code or from the beginning.
	if fromCode == "" {
		lastCode, err := cm.Storage.Dm.GetLatestCreatedScrapCode()
		if err != nil {
			return fmt.Errorf("could not get latest image code, err: %v", err)
		}
		if lastCode == nil {
			fromCode = "0"
		} else {
			fromCode = *lastCode
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errcList []<-chan error

	// create an index with the code to start from
	index := createResumeCodeNumber(&fromCode)
	fmt.Printf("Starting from Code: %s\n", index.String())

	// produce codes in a channel and expose a produceMoreCodes channel to enable
	// a feedback loop
	codes, produceMoreCodes := produceCodes(ctx, index, iterations, workerNumber)

	// filter out codes if they already exist in the DB.
	filteredCodes, filterErrors := filterCodes(ctx, cm.Storage, codes, produceMoreCodes, workerNumber)
	errcList = append(errcList, filterErrors)

	// generate image entries on db and mark them as pending
	pendingImages, pendingErrors := generatePendingImages(ctx, cm.Storage, filteredCodes)
	errcList = append(errcList, pendingErrors)

	// initialize an empty pool of workers
	downloadWorkers := make([]<-chan ScrapedFile, workerNumber)
	downloadWorkerErrors := make(<-chan error, 1)

	// start workers
	for i := 0; i < workerNumber; i++ {
		downloadWorkers[i], downloadWorkerErrors = downloadImages(ctx, cm.Storage, cm.Scrapper, pendingImages, produceMoreCodes)
		errcList = append(errcList, downloadWorkerErrors)
	}

	// fan-in download workers
	downloadedImages := mergeDownloads(ctx, downloadWorkers...)

	// initialize an empty pool of workers
	saveWorkers := make([]<-chan Scrap, workerNumber)
	saveWorkersErrors := make(<-chan error, 1)

	// start workers
	for i := 0; i < workerNumber; i++ {
		saveWorkers[i], saveWorkersErrors = saveImages(cm.Storage, ctx, downloadedImages)
		errcList = append(errcList, saveWorkersErrors)
	}

	downloadCount := 0
	for range mergeSaves(ctx, saveWorkers...) {

		downloadCount++
		fmt.Printf("DOWNLOADED AN IMAGE, TOTAL: %d\n", downloadCount)
		if downloadCount >= iterations {
			fmt.Printf("WE SHOULD FINISH NOW!\n")
			//we dont need more codes
			cancel()
			break

		}
	}
	result := waitForPipeline(errcList...)
	duration := time.Since(start)
	// 		// // Formatted string, such as "2h3m0.5s" or "4.503μs"
	// 		fmt.Printf("Total Duration: %s ",duration)
	fmt.Printf("OPERATION COMPLETED: TIME SPENT: %s\n", duration)
	return result
}

// WaitForPipeline waits for results from all error channels.
// It returns early on the first error.
func waitForPipeline(errs ...<-chan error) error {
	errc := mergeErrors(errs...)
	fmt.Print("Waiting for Pipeline to finish!")
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

// MergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
func mergeErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	// We must ensure that the output channel has the capacity to
	// hold as many errors
	// as there are error channels.
	// This will ensure that it never blocks, even
	// if WaitForPipeline returns early.
	out := make(chan error, len(cs))
	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls
	// wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	// Start a goroutine to close out once all the output goroutines
	// are done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// produceCodes generates and feeds the pipeline with codes.
func produceCodes(
	ctx context.Context,
	index customNumber.Number,
	iterations int,
	channelSize int,
) (<-chan string, chan struct{}) {
	produceMoreCodes := make(chan struct{}, iterations+1)
	codes := make(chan string, iterations+1)
	iterationsCounter := 0
	fmt.Printf("PRODUCING CODES")
	go func() {
		defer close(codes)
		defer close(produceMoreCodes)

		for {
			if iterationsCounter < iterations {
				produceMoreCodes <- struct{}{}
				iterationsCounter++
			}
			fmt.Println("IM HERE")

			select {
			case <-produceMoreCodes:
				fmt.Printf("iterationsCounter: %d  iterations: %d\n", iterationsCounter, iterations)
				//fmt.Printf("PRODUCING CODE: %s \n", index.SmartString())
			codesFor:
				for {
					select {
					case codes <- index.String():
						index.Increment()
						break codesFor
					case <-time.After(1 * time.Second):
						fmt.Println("DEADLOCK TIMEOUT PRODUCE CODES")
						break
					}
				}
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE on produce codes")
				return
			}
		}
	}()
	return codes, produceMoreCodes
}

func filterCodes(
	ctx context.Context,
	storage Storage,
	codes <-chan string,
	produceMoreCodes chan<- struct{},
	channelSize int,
) (<-chan string, <-chan error) {

	filteredCodes := make(chan string, channelSize)
	errc := make(chan error, 1)

	go func() {
		defer close(filteredCodes)
		defer close(errc)

		for code := range codes {
			exists, err := storage.Dm.CodeAlreadyExists(code)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			if exists {

			for1:
				for {
					select {
					case produceMoreCodes <- struct{}{}:
						break for1
					case <-time.After(1 * time.Second):
						fmt.Println("DEADLOCK TIMEOUT PRODUCE MORE CODES")
						break
					}
				}
				fmt.Printf("Image %s already exists, asking for another code.\n", code)
				continue
			}
			fmt.Printf("Image %s does not exist, image will be downloaded.\n", code)

			select {
			case filteredCodes <- code:
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			}
		}
	}()
	return filteredCodes, errc
}

func generatePendingImages(
	ctx context.Context,
	storage Storage,
	filteredCodes <-chan string,
) (<-chan Scrap, <-chan error) {

	pendingImages := make(chan Scrap, 10)
	errc := make(chan error, 1)

	go func() {
		defer close(pendingImages)
		defer close(errc)

		for code := range filteredCodes {
			pendingImage := Scrap{
				RefCode:       code,
				Status:        StatusPending,
				CodeCreatedAt: time.Now(),
			}
			fmt.Printf("Creating an entry on DB for: %s\n", code)

			_, err := storage.Dm.CreateScrap(pendingImage)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}

			select {
			case pendingImages <- pendingImage:
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			}
		}
	}()
	return pendingImages, errc
}

func downloadImages(
	ctx context.Context,
	storage Storage,
	scrapper Scrapper,
	pendingImages <-chan Scrap,
	produceMoreCodes chan<- struct{},
) (<-chan ScrapedFile, <-chan error) {

	imagesToSave := make(chan ScrapedFile, 10)
	errc := make(chan error, 1)

	go func() {
		defer close(imagesToSave)
		defer close(errc)

		for image := range pendingImages {
			ScrapedFile, err := scrapper.ScrapeByCode(image.RefCode)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			//If the image was not found then we need a new code
			if ScrapedFile.Data == nil && err == nil {
				fmt.Printf("Image %s was not found, requesting a new one! \n", image.RefCode)
				err = storage.Dm.UpdateScrapStatusByCode(image.RefCode, StatusNotFound)
				if err != nil {
					errc <- err
					return
				}
				produceMoreCodes <- struct{}{}
				continue
			}
			err = storage.Dm.UpdateScrapStatusByCode(image.RefCode, StatusOngoing)
			if err != nil {
				errc <- err
				return
			}

			select {
			case imagesToSave <- ScrapedFile:
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			}
		}

	}()
	return imagesToSave, errc
}

func mergeDownloads(ctx context.Context, channels ...<-chan ScrapedFile) <-chan ScrapedFile {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	downloadedImages := make(chan ScrapedFile)
	multiplex := func(c <-chan ScrapedFile) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			case downloadedImages <- i:
			}
		}
	}
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		defer close(downloadedImages)
		wg.Wait()
	}()
	return downloadedImages
}

func saveImages(storage Storage,
	ctx context.Context,
	downloadedImages <-chan ScrapedFile) (
	<-chan Scrap, <-chan error) {

	savedImages := make(chan Scrap, 10)
	errc := make(chan error, 1)

	go func() {
		defer close(savedImages)
		defer close(errc)

		for image := range downloadedImages {
			err := storage.Fm.SaveFile(image)
			if err != nil {
				errc <- err
				return
			}

			ss := Scrap{
				RefCode: image.Code,
				Status:  StatusSuccess,
				FileURI: "SOMEWHERE" + image.Code + "." + image.Type,
			}
			err = storage.Dm.UpdateScrapByCode(ss)
			if err != nil {
				errc <- err
				return
			}
			select {
			case savedImages <- ss:
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			}
		}
	}()
	return savedImages, errc
}

func mergeSaves(ctx context.Context, channels ...<-chan Scrap) <-chan Scrap {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	savedImages := make(chan Scrap)
	multiplex := func(c <-chan Scrap) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			case savedImages <- i:
			}
		}
	}
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		defer close(savedImages)
		wg.Wait()
	}()
	return savedImages
}

func createResumeCodeNumber(code *string) customNumber.Number {
	// if no code was found
	// or if were starting from 0 then start from the beginning.
	if code == nil || *code == "0" {
		return customNumber.NewNumber(CustomNumberDigitValues, "0")
	}

	number := customNumber.NewNumber(CustomNumberDigitValues, *code)
	number.Increment()
	return number
}
