package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lolPants/songe-converter/converter"
	"github.com/lolPants/songe-converter/utils"
	"github.com/otiai10/copy"
)

func processDir() {
	var name string
	_, dirName := filepath.Split(dir)
	if dirName != "" {
		name = dirName
	} else {
		t := strings.TrimSuffix(dir, "/")
		t = strings.TrimSuffix(t, "\\")
		_, dirName = filepath.Split(t)

		if dirName != "" {
			name = dirName
		} else {
			log.Println("Failed to clean the synced directory!")
			exit(1)

			return
		}
	}

	synced := filepath.Join(outputDir, name)
	exists, _ := utils.DirectoryExists(synced)
	if exists {
		err := os.RemoveAll(synced)
		if err != nil {
			log.Println("Failed to clean the synced directory!")
			log.Println("error:", err)
			exit(1)

			return
		}
	}

	err := copy.Copy(dir, synced)
	if err != nil {
		log.Println("Failed to create the synced directory!")
		log.Println("error:", err)
		exit(1)

		return
	}

	r := converter.DirOldToNew(synced, false, false)
	if r.Error != nil {
		log.Println("Failed to convert the beatmap!")
		log.Println("error:", r.Error)

		return
	}

	log.Println("Converted successfully!")
}
