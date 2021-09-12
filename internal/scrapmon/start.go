package scrapmon

import (
	"context"
	"fmt"
	customNumber "github.com/slysterous/custom-number"
	"sync"
	"time"
)

// StartCommand is what happens when the command is executed.
func (ccm ConcurrentCommandManager) StartCommand(fromCode string, iterations int, workerNumber int) error {
	// mark the start time
	start := time.Now()

	// if no code was provided, then we resume from the last created code or from the beginning.
	if fromCode == "" {
		lastCode, err := ccm.Storage.Dm.GetLatestCreatedScrapCode()
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
	ccm.Logger.Infof("Starting from Code: %s\n", index.String())

	// produce codes in a channel and expose a produceMoreCodes channel to enable
	// a feedback loop
	codes, produceMoreCodes := ccm.CodeAuthority.Produce(ctx, index, iterations, workerNumber)

	// filter out codes if they already exist in the DB.
	filteredCodes, filterErrors := ccm.CodeAuthority.Filter(ctx, ccm.Storage, codes, produceMoreCodes, workerNumber)
	errcList = append(errcList, filterErrors)

	// generate image entries on db and mark them as pending
	pendingFiles, pendingErrors := generatependingFiles(ctx, ccm.Storage, filteredCodes)
	errcList = append(errcList, pendingErrors)

	// initialize an empty pool of workers
	downloadWorkers := make([]<-chan ScrapedFile, workerNumber)
	downloadWorkerErrors := make(<-chan error, 1)

	// start workers
	for i := 0; i < workerNumber; i++ {
		downloadWorkers[i], downloadWorkerErrors = ccm.FileScrapper.DownloadFiles(ctx, ccm.Storage, pendingFiles, produceMoreCodes)
		errcList = append(errcList, downloadWorkerErrors)
	}

	// fan-in download workers
	downloadedImages := mergeDownloads(ctx, downloadWorkers...)

	// initialize an empty pool of workers
	saveWorkers := make([]<-chan Scrap, workerNumber)
	saveWorkersErrors := make(<-chan error, 1)

	// start workers
	for i := 0; i < workerNumber; i++ {
		saveWorkers[i], saveWorkersErrors = ccm.FileScrapper.SaveFiles(ccm.Storage, ctx, downloadedImages)
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

func generatependingFiles(
	ctx context.Context,
	storage Storage,
	filteredCodes <-chan string,
) (<-chan Scrap, <-chan error) {

	pendingFiles := make(chan Scrap, 10)
	errc := make(chan error, 1)

	go func() {
		defer close(pendingFiles)
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
			case pendingFiles <- pendingImage:
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE")
				return
			}
		}
	}()
	return pendingFiles, errc
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
