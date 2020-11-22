package domain

import (
	"fmt"
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
	Code string
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
	SaveFile(src ScrapedImage) error
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

// PurgeCommand is what happens when the command is executed.
func (cm CommandManager) PurgeCommand() error {
	err := cm.Storage.Purge()
	if err != nil {
		return fmt.Errorf("could not purge storage, err: %v", err)
	}
	return nil
}
