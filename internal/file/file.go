package file

import (
	"fmt"
	"io"
	"os"
)

// Manager is
type Manager struct {
}

//NewManager constructs a new file manager.
func NewManager() *Manager {
	return &Manager{}
}

func (m Manager) SaveImage(src io.Reader, path string) (err error) {
	//Create a empty file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		errc := file.Close()
		if errc != nil {
			err = fmt.Errorf("error while closeing file: %v", err)
		}
	}()
	//Write the bytes to the fiel
	_, err = io.Copy(file, src)
	if err != nil {
		return err
	}
	return nil
}
