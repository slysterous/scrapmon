package domain

// Config represents the applications configuration parameters
type Config struct {
	Env               string
	DatabaseUser      string
	DatabasePassword  string
	DatabaseHost      string
	DatabasePort      string
	DatabaseName      string
	HTTPClientTimeout int
	MaxDBConnections  int
	TorHost           string
	TorPort           string
}

// Storage defines the different types of storage.
type Storage struct {
	Fm FileManager
	Dm DatabaseManager
}

// Screenshot defines a scrapped screenshot.
type Screenshot struct {
	RefCode string
	FileURI string
}

// Purger defines the purging behaviour.
type Purger interface {
	Purge() error
}

// DatabaseManager defines the storage management behaviour.
type DatabaseManager interface {
	GetScreenshotByCode()
	SaveScreenshotRef()
	Purger
}

// FileManager defins the file management behaviour.
type FileManager interface {
	SaveScreenshot()
	GetScreenshot()
	Purger
}

// ScreenshotScrapper defines the scrapping bevahiour.
type ScreenshotScrapper interface {
	ScrapeScreenshotByCode()
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

// PRNT.SCR SITE to scrap

// DB to save ref codes and file urls

// Filesystem to save files
