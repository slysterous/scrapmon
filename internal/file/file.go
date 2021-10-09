package file

import (
	"fmt"
	"os"
	"path"

	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
)

//go:generate mockgen -destination mock/file.go -package file_mock . Writer,Purger

// Manager is the file manager.
type Manager struct {
	ScrapFolder string
	Writer      Writer
	Purger      Purger
}

// Writer is the file writer
type Writer interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// Purger is removing files
type Purger interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	RemoveAll(path string) error
}

// NewManager constructs a new file manager.
func NewManager(imageFolder string, writer Writer, purger Purger) *Manager {
	return &Manager{
		ScrapFolder: imageFolder,
		Writer:      writer,
		Purger:      purger,
	}
}

// SaveFile saves image bytes to a specified file.
func (m Manager) SaveFile(src scrapmon.ScrapedFile) error {
	//Write the bytes to the file
	err := m.Writer.WriteFile(m.composeFilePath(src.Code, src.Type), src.Data, 0644)
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
	dir, err := m.Purger.ReadDir(m.ScrapFolder)
	if err != nil {
		return fmt.Errorf("file: could not read scrap directory, err: %v", err)
	}
	for _, d := range dir {
		a := path.Join([]string{"tmp", d.Name()}...)
		err = m.Purger.RemoveAll(a)
		if err != nil {
			return fmt.Errorf("file: could not delete scrap, err: %v", err)
		}
	}
	return nil
}
