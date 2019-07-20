package main

import (
	"encoding/json"
	"errors"

	"github.com/TomOnTime/utfutil"
)

func readInfo(path string) (OldInfoJSON, error) {
	bytes, err := utfutil.ReadFile(path, utfutil.UTF8)
	if err != nil {
		return OldInfoJSON{}, err
	}

	valid := IsJSON(bytes)
	if valid == false {
		invalidError := errors.New("Invalid info.json")
		return OldInfoJSON{}, invalidError
	}

	var infoJSON OldInfoJSON
	json.Unmarshal(bytes, &infoJSON)

	return infoJSON, nil
}
