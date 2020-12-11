package file

import (
	"path"
	"fmt"
	printscrape "github.com/slysterous/print-scrape/internal/printscrape"
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
	err := ioutil.WriteFile(m.composeFilePath(src.Code,src.Type), src.Data, 0644)
	if err != nil {
		return fmt.Errorf("file: could not create file, err: %v", err)
	}
	return nil
}

func (m Manager) composeFilePath(code,ext string) string {
	return path.Join(m.ImageFolder,code+"."+ext)
}

// Purge deletes every file from the file system
func (m Manager) Purge() error {
	// TODO fetch config to get the path in which everything is saved
	return nil
}
