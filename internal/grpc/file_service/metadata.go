package fileservice

import (
	"encoding/json"
	"os"
	"path/filepath"

)

const metaFile = "files_meta.json"
const unknownHost = "unknown"



func (fs *FileServer) SaveMetaData() error {
	file, err := os.Create(filepath.Join(fs.serviceFilesDir, metaFile))
	if err != nil {

		return err
	}
	err = json.NewEncoder(file).Encode(fs.metaData)
	return err
}

func (fs *FileServer) LoadMetaData() error {
	file, err := os.Open(filepath.Join(fs.serviceFilesDir, metaFile))
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}
	err = json.NewDecoder(file).Decode(fs.metaData)
	return err
}
