package domain

import (
	"context"
	"fmt"
	customNumber "github.com/slysterous/print-scrape/pkg/customnumber"
	"sync"
	"time"
)

// StartCommand is what happens when the command is executed.
func (cm CommandManager) StartCommand(fromCode string, iterations int) error {
	start:=time.Now()

	//imageCount:= 0
	//if no code was provided, then we resume from the last created code or from the beginning.
	if fromCode == "" {
		lastCode, err := cm.Storage.Dm.GetLatestCreatedScreenShotCode()
		if err != nil {
			return fmt.Errorf("could not get latest image code, err: %v", err)
		}
		if lastCode == nil {
			fromCode = "0"
		} else {
			fromCode = *lastCode
		}
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var errcList []<-chan error

	index := createResumeCodeNumber(&fromCode)

	fmt.Printf("Starting from Code: %s\n",index.SmartString())
	
	//concurrency starts here
	//this will produce codes until we reach the iterations or until ctx.done is called.
	codes,produceMoreCodes,completedCodes := produceCodes(ctx, index, iterations)


	// filter out codes if they already exist in the DB.
	filteredCodes, filterErrors := filterCodes(ctx, cm.Storage, codes, produceMoreCodes)
	errcList = append(errcList, filterErrors)

	pendingImages, pendingErrors := generatePendingImages( ctx,cm.Storage,filteredCodes)
	errcList = append(errcList, pendingErrors)

	downloadWorkers := make([]<-chan ScrapedImage, 10)
	downloadWorkerErrors := make(<-chan error, 1)

	for i := 0; i < 10; i++ {
		fmt.Printf("INITIALIZING DOWNLOAD WORKER %d/10 \n", i+1)
		downloadWorkers[i], downloadWorkerErrors = downloadImages(ctx, cm.Storage, cm.Scrapper, pendingImages, produceMoreCodes)
		errcList = append(errcList, downloadWorkerErrors)
	}

	downloadedImages := mergeDownloads(ctx, downloadWorkers...)

	saveWorkers := make([]<-chan ScreenShot, 10)
	saveWorkersErrors := make(<-chan error, 1)

	for i := 0; i < 10; i++ {
		saveWorkers[i], saveWorkersErrors = saveImages(cm.Storage, ctx, downloadedImages)
		errcList = append(errcList, saveWorkersErrors)
	}

	downloadCount := 0
	for range mergeSaves(ctx, saveWorkers...) {
	
		completedCodes <- struct{}{}
		downloadCount++
		fmt.Printf("DOWNLOADED AN IMAGE, TOTAL: %d\n", downloadCount)
		if downloadCount>=10{
			fmt.Printf("WE SHOULD FINISH NOW!")
			//we dont need more codes

			break

		}
	}
	result:=waitForPipeline(errcList...)
	duration := time.Since(start)
	// 		// // Formatted string, such as "2h3m0.5s" or "4.503μs"
	// 		fmt.Printf("Total Duration: %s ",duration)
	fmt.Printf("OPERATION COMPLETED: TIME SPENT: %s",duration)
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
) (<-chan string,chan struct{},chan struct{}) {
	produceMoreCodes:=make(chan struct{},10)
	completedCodes:=make(chan struct{},10)
	codes := make(chan string, 10)
	iterationsCounter := 0
	fmt.Printf("PRODUCING CODES")
	go func() {
		defer close(codes)
		defer close(produceMoreCodes)
		defer close(completedCodes)

		for {
			
			select {
			case <-produceMoreCodes:
				iterationsCounter--
			case <-completedCodes:
				iterationsCounter++
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE on produce codes")
				return
			case codes <- index.SmartString():
				if iterationsCounter < iterations {
					fmt.Printf("PRODUCING CODE: %s \n", index.SmartString())
					
					index.Increment()
				}else{
					return
				}
			}
		}
	}()
	return codes,produceMoreCodes,completedCodes
}

func filterCodes(
	ctx context.Context,
	storage Storage,
	codes <-chan string,
	produceMoreCodes chan<- struct{},
) (<-chan string, <-chan error) {

	filteredCodes := make(chan string, 10)
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
				produceMoreCodes <- struct{}{}
				fmt.Printf("Image %s already exists, asking for another code.\n", code)
				continue
			}
			fmt.Printf("Image %s does not exist, moving on.\n", code)

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
) (<-chan ScreenShot, <-chan error) {

	pendingImages := make(chan ScreenShot, 10)
	errc := make(chan error, 1)

	go func() {
		defer close(pendingImages)
		defer close(errc)

		for code := range filteredCodes {
			pendingImage := ScreenShot{
				RefCode:       code,
				Status:        StatusPending,
				CodeCreatedAt: time.Now(),
			}
			fmt.Printf("Creating an entry on DB for: %s\n", code)
			
			_, err := storage.Dm.CreateScreenShot(pendingImage)
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
	scrapper ImageScrapper,
	pendingImages <-chan ScreenShot,
	produceMoreCodes chan<- struct{},
) (<-chan ScrapedImage, <-chan error) {

	imagesToSave := make(chan ScrapedImage, 10)
	errc := make(chan error, 1)

	go func() {
		defer close(imagesToSave)
		defer close(errc)

		for image := range pendingImages {
			scrapedImage, err := scrapper.ScrapeImageByCode(image.RefCode)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			//If the image was not found then we need a new code
			if scrapedImage.Data == nil && err == nil {
				fmt.Printf("Image %s was not found, requesting a new one! \n", image.RefCode)
				err = storage.Dm.UpdateScreenShotStatusByCode(image.RefCode, StatusNotFound)
				if err != nil {
					errc <- err
					return
				}
				produceMoreCodes <- struct{}{}
				continue
			}
			err = storage.Dm.UpdateScreenShotStatusByCode(image.RefCode, StatusOngoing)
			if err != nil {
				errc <- err
				return
			}

			select {
			case imagesToSave <- scrapedImage:
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			}
		}

	}()
	return imagesToSave, errc
}

func mergeDownloads(ctx context.Context, channels ...<-chan ScrapedImage) <-chan ScrapedImage {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	downloadedImages := make(chan ScrapedImage)
	multiplex := func(c <-chan ScrapedImage) {
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
	downloadedImages <-chan ScrapedImage) (
	<-chan ScreenShot, <-chan error) {

	savedImages := make(chan ScreenShot, 10)
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

			ss := ScreenShot{
				RefCode: image.Code,
				Status:  StatusSuccess,
				FileURI: "SOMEWHERE" + image.Code + "." + image.Type,
			}
			err = storage.Dm.UpdateScreenShotByCode(ss)
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

func mergeSaves(ctx context.Context, channels ...<-chan ScreenShot) <-chan ScreenShot {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	savedImages := make(chan ScreenShot)
	multiplex := func(c <-chan ScreenShot) {
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