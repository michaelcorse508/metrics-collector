package jsonfilesaver

import (
	"os"
)

type RawFileSaver struct {
	saveFilePath string
}

// NewRawFileSaver creates the instance of RawFileSaver struct with given file path where to save data.
func NewRawFileSaver(saveFilePath string) *RawFileSaver {
	return &RawFileSaver{
		saveFilePath: saveFilePath,
	}
}

// Save stores data to file set in RawFileSaver.
func (fs *RawFileSaver) Save(data []byte) error {
	file, err := os.OpenFile(fs.saveFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Load loads data from file set in RawFileSaver.
func (fs *RawFileSaver) Load() ([]byte, error) {
	data, err := os.ReadFile(fs.saveFilePath)
	return data, err
}
