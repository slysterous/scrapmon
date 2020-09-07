package file

import (
	"fmt"
	"io/ioutil"
)

// Manager is
type Manager struct {
}

// NewManager constructs a new file manager.
func NewManager() *Manager {
	return &Manager{}
}

// SaveFile saves image bytes to a specified file.
func (m Manager) SaveFile(src *[]byte, path string) error {

	//Write the bytes to the file
	err := ioutil.WriteFile(path, *src, 0644)
	if err != nil {
		return fmt.Errorf("file: could not create file, err: %v", err)
	}

	return nil
}

// Purge deletes every file from the file system
func (m Manager) Purge() error {
	return nil
}
