package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
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
			return "", err
		}

		defer file.Close()
		bytes, _ := ioutil.ReadAll(file)
		allBytes = append(allBytes, bytes...)
	}

	return calculateHashMD5(allBytes), nil
}

func calculateHashMD5(bytes []byte) string {
	sum := md5.Sum(bytes)
	return hex.EncodeToString(sum[:])
}

func calculateHashSHA1(bytes []byte) string {
	sum := sha1.Sum(bytes)
	return hex.EncodeToString(sum[:])
}
