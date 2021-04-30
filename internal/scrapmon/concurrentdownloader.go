package scrapmon

import (
	"context"
	"fmt"

	"github.com/slysterous/scrapmon/internal/log"
)

//go:generate mockgen -destination mock/codeproducer.go -package scrapmon_mock . ConcurrentDownloader

// ConcurrentDownloader describes the actions of a file downloader.
type ConcurrentDownloader interface {
	DownloadFiles(
		ctx context.Context,
		storage Storage,
		scrapper Scrapper,
		pendingFiles <-chan Scrap,
		produceMoreCodes chan<- struct{},
	) (<-chan ScrapedFile, <-chan error)
	SaveFiles(storage Storage,
		ctx context.Context,
		downloadedImages <-chan ScrapedFile) (
		<-chan Scrap, <-chan error)
}

// ConcurrentScrapper is responsible for code creation and handling.
type ConcurrentScrapper struct {
	Logger log.Logger
}

func (cd ConcurrentScrapper) DownloadFiles(
	ctx context.Context,
	storage Storage,
	scrapper Scrapper,
	pendingFiles <-chan Scrap,
	produceMoreCodes chan<- struct{},
) (<-chan ScrapedFile, <-chan error) {

	imagesToSave := make(chan ScrapedFile, 10)
	errc := make(chan error, 1)

	go func() {
		defer close(imagesToSave)
		defer close(errc)

		for image := range pendingFiles {
			ScrapedFile, err := scrapper.ScrapeByCode(image.RefCode, "png")
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			//If the image was not found then we need a new code.
			if ScrapedFile.Data == nil {
				cd.Logger.Infof("File %s was not found, requesting a new one!\n", image.RefCode)
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
				cd.Logger.Debugf("Finished downloading Files!\n")
				return
			}
		}

	}()
	return imagesToSave, errc
}

func (cd ConcurrentScrapper) SaveFiles(
	storage Storage,
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
