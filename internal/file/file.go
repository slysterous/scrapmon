package file

import (
	"fmt"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	"io/ioutil"
)

// Manager is
type Manager struct {
	ImageFolder string
}

// NewManager constructs a new file manager.
func NewManager(imageFolder string) *Manager {
	return &Manager{
		ImageFolder: imageFolder,
	}
}

// SaveFile saves image bytes to a specified file.
func (m Manager) SaveFile(src printscrape.ScrapedImage) error {

	//Write the bytes to the file
	err := ioutil.WriteFile(m.ImageFolder+"/"+src.Code+"."+src.Type, src.Data, 0644)
	if err != nil {
		return fmt.Errorf("file: could not create file, err: %v", err)
	}

	return nil
}

// Purge deletes every file from the file system
func (m Manager) Purge() error {
	return nil
}
