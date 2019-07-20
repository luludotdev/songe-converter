package main

import (
	"encoding/json"
	"errors"

	"github.com/TomOnTime/utfutil"
)

func readDifficulty(path string) (OldDifficultyJSON, error) {
	bytes, err := utfutil.ReadFile(path, utfutil.UTF8)
	if err != nil {
		return OldDifficultyJSON{}, err
	}

	valid := IsJSON(bytes)
	if valid == false {
		invalidError := errors.New("Invalid difficulty file found at \"" + path + "\"")
		return OldDifficultyJSON{}, invalidError
	}

	var diffJSON OldDifficultyJSON
	json.Unmarshal(bytes, &diffJSON)

	return diffJSON, nil
}
