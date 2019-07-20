package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func calculateOldHash(infoJSON OldInfoJSON, dir string) (string, error) {
	allBytes := make([]byte, 0)
	for _, diff := range infoJSON.DifficultyLevels {
		path := filepath.Join(dir, diff.JSONPath)
		file, err := os.Open(path)

		if err != nil {
			file.Close()
			continue
		}

		defer file.Close()
		bytes, _ := ioutil.ReadAll(file)
		allBytes = append(allBytes, bytes...)
	}

	return calculateHashMD5(allBytes), nil
}
