package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
)

func calculateOldHash(infoJSON OldInfoJSON, dir string) (string, error) {
	allBytes := make([]byte, 0)
	for _, diff := range infoJSON.DifficultyLevels {
		path := path.Join(dir, diff.JSONPath)
		file, err := os.Open(path)

		if err != nil {
			return "", err
		}

		defer file.Close()
		bytes, _ := ioutil.ReadAll(file)
		allBytes = append(allBytes, bytes...)
	}

	return calculateHash(allBytes), nil
}

func calculateHash(bytes []byte) string {
	sum := md5.Sum(bytes)
	return hex.EncodeToString(sum[:])
}