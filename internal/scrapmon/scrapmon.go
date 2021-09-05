package scrapmon

import (
	"fmt"
	"time"
)

//go:generate mockgen -destination mock/scrapmon.go -package scrapmon_mock . DatabaseManager,FileManager,Scrapper

// CustomNumberDigitValues defines the allowed digits of the custom arithmetic system to be used
//var CustomNumberDigitValues = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var CustomNumberDigitValues = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

// ScrapStatus describes the status of a Scrap.
type ScrapStatus string

// Possible Scrap Status  values.
const (
	StatusPending  ScrapStatus = "pending"
	StatusOngoing  ScrapStatus = "ongoing"
	StatusSuccess  ScrapStatus = "success"
	StatusFailure  ScrapStatus = "failure"
	StatusNotFound ScrapStatus = "notfound"
)

// Config represents the applications configuration parameters
type Config struct {
	Env                string
	DatabaseUser       string
	DatabasePassword   string
	DatabaseHost       string
	DatabasePort       string
	DatabaseName       string
	HTTPClientTimeout  int
	MaxDBConnections   int
	TorHost            string
	TorPort            string
	ScrapStorageFolder string
}

// Storage defines the different types of storage.
type Storage struct {
	Fm FileManager
	Dm DatabaseManager
}

// ConcurrentCommandManager handles commands.
type ConcurrentCommandManager struct {
	Storage       Storage
	Scrapper      Scrapper
	CodeAuthority ConcurrentCodeProducer
	FileScrapper  ConcurrentDownloader
}

// Scrap defines a scrapped Scrap.
type Scrap struct {
	ID            int64
	RefCode       string
	CodeCreatedAt time.Time
	FileURI       string
	Status        ScrapStatus
}

// ScrapedFile describes the scraped file and its properties.
type ScrapedFile struct {
	Code string
	Data []byte
	Type string
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Purger defines the purging behaviour.
type Purger interface {
	Purge() error
}

// DatabaseManager defines the storage management behaviour.
type DatabaseManager interface {
	CreateScrap(s Scrap) (int, error)
	UpdateScrapStatusByCode(code string, status ScrapStatus) error
	UpdateScrapByCode(s Scrap) error
	GetLatestCreatedScrapCode() (*string, error)
	CodeAlreadyExists(code string) (bool, error)
	Purger
}

// FileManager defines the file management behaviour.
type FileManager interface {
	SaveFile(src ScrapedFile) error
	Purger
}

// Scrapper defines the scrapping behaviour.
type Scrapper interface {
	ScrapeByCode(code, ext string) (ScrapedFile, error)
}

// StartLogic describes how a StartLogic function should be described.
type StartLogic func(fromCode string, iterations int, workerNumber int) error

// PurgeLogic describes how a PurgeLogic function should be described.
type PurgeLogic func() error

//Purge will clear all data saved in files and database
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

//PurgeCommand is what happens when the command is executed.
func (ccm ConcurrentCommandManager) PurgeCommand() error {
	err := ccm.Storage.Purge()
	if err != nil {
		return fmt.Errorf("could not purge storage, err: %v", err)
	}
	return nil
}
