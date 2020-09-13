package domain

import (
	"fmt"
	customNumber "github.com/slysterous/print-scrape/pkg/customnumber"
	"strings"
	"time"
)

// CustomNumberDigitValues defines the allowed digits of the custom arithmetic system to be used
//var CustomNumberDigitValues = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var CustomNumberDigitValues = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

// ScreenShotStatus describes the status of a ScreenShot.
type ScreenShotStatus string

// Possible ScreenShot Status  values.
const (
	StatusPending  ScreenShotStatus = "pending"
	StatusOngoing  ScreenShotStatus = "ongoing"
	StatusSuccess  ScreenShotStatus = "success"
	StatusFailure  ScreenShotStatus = "failure"
	StatusNotFound ScreenShotStatus = "notfound"
)

// Config represents the applications configuration parameters
type Config struct {
	Env                     string
	DatabaseUser            string
	DatabasePassword        string
	DatabaseHost            string
	DatabasePort            string
	DatabaseName            string
	HTTPClientTimeout       int
	MaxDBConnections        int
	TorHost                 string
	TorPort                 string
	ScreenShotStorageFolder string
}

// Storage defines the different types of storage.
type Storage struct {
	Fm FileManager
	Dm DatabaseManager
}

// CommandManager handles commands.
type CommandManager struct {
	Storage  Storage
	Scrapper ImageScrapper
}

// ScreenShot defines a scrapped ScreenShot.
type ScreenShot struct {
	ID            int64
	RefCode       string
	CodeCreatedAt time.Time
	FileURI       string
	Status        ScreenShotStatus
}

// Service describes i don't know
type Service struct {
	storage  Storage
	scrapper ImageScrapper
}

type ScrapedImage struct {
	Data []byte
	Type string
}

// CommandFunction defines a function that contains the logic of a command.
type CommandFunction func() error

// // CommandHandler defines the cli client interactions.
// type CommandHandler interface {
// 	HandleStartCommand(ctx context.Context,fn CommandFunction) error
// }

// Purger defines the purging behaviour.
type Purger interface {
	Purge() error
}

// DatabaseManager defines the storage management behaviour.
type DatabaseManager interface {
	CreateScreenShot(ss ScreenShot) (int, error)
	UpdateScreenShotStatusByCode(code string, status ScreenShotStatus) error
	UpdateScreenShotByCode(ss ScreenShot) error
	GetLatestCreatedScreenShotCode() (*string, error)
	CodeAlreadyExists(code string) (bool, error)
	Purger
}

// FileManager defins the file management behaviour.
type FileManager interface {
	SaveFile(src ScrapedImage, path string) error
	Purger
}

// ImageScrapper defines the scrapping behaviour.
type ImageScrapper interface {
	ScrapeImageByCode(code string) (ScrapedImage, error)
}

// Purge will clear all data saved in files and database
func (s *Storage) Purge() error {
	err := s.Dm.Purge()
	if err != nil {
		return err
	}
	err = s.Fm.Purge()
	if err != nil {
		return err
	}
	return nil
}

// IsScreenShotURLValid checks if a screenshot url is valid to be processed.
func IsScreenShotURLValid(url string) bool {
	fmt.Println(url)
	return strings.Contains(url, "https://")
}

// StartCommand is what happens when the command is executed.
func (cm CommandManager) StartCommand(fromCode string, iterations int) error {
	//start:=time.Now()

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
	produceMoreCodes := make(chan struct{}, 10)
	done := make(chan struct{}, 1)
	defer close(done)
	defer close(produceMoreCodes)

	index := createResumeCodeNumber(&fromCode)

	codes := produceCodes(done, produceMoreCodes, index, iterations)

	filteredCodes,_:=filterCodes(cm.Storage,done,codes,produceMoreCodes)

	pendingImages,_:=generatePendingImages(cm.Storage,done,filteredCodes)

	downloadWorkers:=make([]<-chan ScrapedImage,10)
	downloadWorkersErrors:=make([]<-chan error,1)

	for i:=0; i<10;i++{
		downloadWorkers[i],downloadWorkersErrors[i]=downloadImages(cm.Storage,cm.Scrapper,done,pendingImages,produceMoreCodes)
	}



	// numCodes := 0
	// malakiesTouTolh := 0
	// for code := range codes {
	// 	//assume http call failure!
	// 	//httpcall()
	// 	//if err
	// 	malakiesTouTolh++
	// 	fmt.Printf("CODE%d: %s\n", numCodes, code)
	// 	if malakiesTouTolh%10 == 0 {
	// 		produceMoreCodes <- struct{}{}
	// 		continue
	// 	}
	// 	// if numCodes%4 == 0 {
	// 	// 	produceMoreCodes <- struct{}{}
	// 	// 	//numCodes--
	// 	// 	continue
	// 	// }
	// 	//time.Sleep(1 * time.Second)

	// 	numCodes++
	// 	if numCodes >= iterations {

	// 		done <- struct{}{}
	// 		break
	// 	}

	// }
	return nil
}

func generatePendingImages(storage Storage,
	done <-chan struct{},
	filteredCodes <-chan string) (
	<-chan ScreenShot,<-chan error) {

	pendingImages:= make(chan ScreenShot,10)
	errc:=make(chan error,1)

	go func() {
		defer close(pendingImages)
		defer close(errc)

		for code := range filteredCodes {
			select {
			case <-done:
				return
			default:

				pendingImage:= ScreenShot{
					RefCode:code,
					Status: StatusPending,
					CodeCreatedAt: time.Now(),
				}
				_,err:=storage.Dm.CreateScreenShot(pendingImage)
				if err != nil {
					// Handle an error that occurs during the goroutine.
					errc<-err
					return
				}
				pendingImages <- pendingImage

			}
		}
	}()
	return pendingImages,errc
}
	

func downloadImages(storage Storage,
	scrapper ImageScrapper,
	done <-chan struct{},
	pendingImages <-chan ScreenShot,
	produceMoreCodes chan<- struct{}) (
	<-chan ScrapedImage, <-chan error)  {
	
	imagesToSave:= make(chan ScrapedImage,10)

	errc := make(chan error,1)

	go func() {
		defer close(imagesToSave)
		defer close(errc)

		for image :=range pendingImages {
			select {
			case <-done:
				return
			default:
				scrapedImage,err:=scrapper.ScrapeImageByCode(image.RefCode)
				if err != nil {
					// Handle an error that occurs during the goroutine.
					errc<-err
					return
				}
				imagesToSave <- scrapedImage
			
			}
		}

	}()
		return imagesToSave,errc
	}


	

	func filterCodes(storage Storage,
	done <-chan struct{},
	codes <-chan string,
	produceMoreCodes chan<- struct{}) (
	<-chan string,<-chan error) {

	usefulCodes := make(chan string, 10)
	errc := make(chan error,1)

	go func() {
		defer close(usefulCodes)
		defer close(errc)
		
		for code := range codes {
			select {
			case <-done:
				return
			default:
				exists, err := storage.Dm.CodeAlreadyExists(code)
				if err != nil {
					// Handle an error that occurs during the goroutine.
					errc<-err
					return
				}
				if exists {
					produceMoreCodes <- struct{}{}
					break
				}
				usefulCodes <- code
			}
		}
	}()
	return usefulCodes,errc
}


// produceCodes generates and feeds the pipeline with codes.
func produceCodes(done <-chan struct{}, produceMoreCodes <-chan struct{}, index customNumber.Number, iterations int) <-chan string {
	codes := make(chan string, 10)
	completedCodes := 0

	go func() {
		defer close(codes)
		for {
			select {
			case <-produceMoreCodes:
				completedCodes--
			case <-done:
				return
			default:
			}
			if completedCodes < iterations {
				fmt.Printf("PRODUCING CODE: %s \n", index.SmartString())
				codes <- index.SmartString()
				index.Increment()
				completedCodes++
			}
		}
	}()
	return codes
}

// // produceCodes feeds
// func produceCodes(done <-chan struct{},index customNumber.Number) (<-chan string) {
// 	codes:= make(chan string,2)
// 	//keep producing codes until
// 	go func() {
// 		defer close(codes)
// 		for {
// 			//time.Sleep(10 * time.Millisecond)
// 			runtime.Gosched()
// 			select {
// 			case <-done:
// 				fmt.Println("CLOSED")
// 				return
// 			default:
// 			}
// 			select {
// 			case <-done:
// 				fmt.Println("CLOSED")
// 				return
// 			case codes <- index.SmartString():
// 				fmt.Printf("Producing code: %s\n",index.SmartString())
// 				index.Increment()
// 			}
// 		}
// 	}()
// 	return codes
// }

// // StartCommand is what happens when the command is executed.
// func (cm CommandManager) StartCommand(fromCode string, iterations int) error {

// 	start:=time.Now()

// 	imageCount:= 0

// 	//if no code was provided, then we resume from the last created code or from the beginning.
// 	if fromCode == "" {
// 		lastCode, err := cm.Storage.Dm.GetLatestCreatedScreenShotCode()
// 		if err != nil {
// 			return fmt.Errorf("could not get latest image code, err: %v", err)
// 		}
// 		if lastCode == nil {
// 			fromCode = "0"
// 		} else {
// 			fromCode = *lastCode
// 		}
// 	}

// 	index := createResumeCodeNumber(&fromCode)

// 	//iterate untill we reach the last possible image or run out of iterations.
// 	for index.String() != "ZZZZZZZZ" && ((imageCount < iterations) || iterations==-1) {
// 		fmt.Printf("ITERATIONS LEFT: %v \n", iterations - imageCount)

// 		existsAlready, err := cm.Storage.Dm.CodeAlreadyExists(index.SmartString())
// 		if err != nil {
// 			return fmt.Errorf("could not get image, err: %v", err)
// 		}

// 		if existsAlready {
// 			index.Increment()
// 			continue
// 		}

// 		screenShot := ScreenShot{
// 			CodeCreatedAt: time.Now(),
// 			RefCode:       index.SmartString(),
// 			FileURI:       "",
// 		}

// 		// start saving item to db with downloadStatus pending
// 		_, err = cm.Storage.Dm.CreateScreenShot(screenShot)
// 		if err != nil {
// 			return fmt.Errorf("could not save screenshot, err: %v", err)
// 		}

// 		// download image
// 		imageTime:=time.Now()
// 		imagedata, imageType, err := cm.Scrapper.ScrapeImageByCode(screenShot.RefCode)
// 		if err != nil {
// 			fmt.Printf("could not download image stream, err: %v", err)
// 			err = cm.Storage.Dm.UpdateScreenShotStatusByCode(screenShot.RefCode, StatusFailure)
// 			if err != nil {
// 				return fmt.Errorf("could not update screenshot status to Failure, err: %v", err)
// 			}

// 			index.Increment()
// 			continue
// 		}

// 		if imagedata == nil {
// 			err = cm.Storage.Dm.UpdateScreenShotStatusByCode(screenShot.RefCode, StatusFailure)
// 			if err != nil {
// 				return fmt.Errorf("could not update screenshot status to Failure, err: %v", err)
// 			}
// 			index.Increment()
// 			continue
// 		}

// 		err = cm.Storage.Dm.UpdateScreenShotStatusByCode(screenShot.RefCode, StatusOngoing)
// 		if err != nil {
// 			return fmt.Errorf("could not update screenshot status to ongoing, err: %v", err)
// 		}

// 		fileURI := "/media/slysterous/HDD Vault/imgur-images/" + screenShot.RefCode + "." + *imageType

// 		err = cm.Storage.Fm.SaveFile(imagedata, fileURI)
// 		if err != nil {
// 			return fmt.Errorf("could not save image to filesystem, err: %v", err)
// 		}

// 		screenShot.FileURI = fileURI

// 		screenShot.Status = StatusSuccess

// 		err = cm.Storage.Dm.UpdateScreenShotByCode(screenShot)

// 		index.Increment()
// 		imageCount++
// 		// Code to measure
// 		duration := time.Since(imageTime)
// 		// // Formatted string, such as "2h3m0.5s" or "4.503μs"
// 		fmt.Printf("DURATION: %s ",duration)
// 	}
// 		// Code to measure
// 		duration := time.Since(start)
// 		// // Formatted string, such as "2h3m0.5s" or "4.503μs"
// 		fmt.Printf("Total Duration: %s ",duration)
// 	return nil
// }

// PurgeCommand is what happens when the command is executed.
func (cm CommandManager) PurgeCommand() error {
	err := cm.Storage.Purge()
	if err != nil {
		return fmt.Errorf("could not purge storage, err: %v", err)
	}
	return nil
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
