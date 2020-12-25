package file

import (
	"fmt"
	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
	"io/ioutil"
	"os"
	"path"
)

// Manager is
type Manager struct {
	ScrapFolder string
}

// NewManager constructs a new file manager.
func NewManager(imageFolder string) *Manager {
	return &Manager{
		ScrapFolder: imageFolder,
	}
}

// SaveFile saves image bytes to a specified file.
func (m Manager) SaveFile(src scrapmon.ScrapedFile) error {
	//Write the bytes to the file
	err := ioutil.WriteFile(m.composeFilePath(src.Code, src.Type), src.Data, 0644)
	if err != nil {
		return fmt.Errorf("file: could not create file, err: %v", err)
	}
	return nil
}

func (m Manager) composeFilePath(code, ext string) string {
	return path.Join(m.ScrapFolder, code+"."+ext)
}

// Purge deletes every file from the file system
func (m Manager) Purge() error {
	dir, err := ioutil.ReadDir(m.ScrapFolder)
	if err !=nil {
		return fmt.Errorf("file: could not read scrap directory, err: %v",err)
	}
	for _, d := range dir {
		err = os.RemoveAll(path.Join([]string{"tmp", d.Name()}...))
		if err!=nil {
			return fmt.Errorf("file: could not delete scrap, err: %v",err)
		}
	}
	return nil
}
