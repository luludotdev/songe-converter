package directory

import (
	"errors"
	"path/filepath"

	"jackbaron.com/songe-converter/v2/utils"
)

// ReadType reads a directory and returns its type
func ReadType(path string) (Type, error) {
	exists, err := utils.DirectoryExists(path)
	if err != nil {
		return None, err
	} else if exists == false {
		return None, errors.New("directory does not exist")
	}

	datPath := filepath.Join(path, "info.dat")
	datExists, _ := utils.FileExists(datPath)

	jsonPath := filepath.Join(path, "info.json")
	jsonExists, _ := utils.FileExists(jsonPath)

	if datExists == true && jsonExists == true {
		return Both, nil
	} else if datExists == true && jsonExists == false {
		return New, nil
	} else if datExists == false && jsonExists == true {
		return Old, nil
	}

	return None, nil
}
