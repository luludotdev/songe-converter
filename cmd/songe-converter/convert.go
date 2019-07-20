package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lolPants/songe-converter/converter"
	"github.com/lolPants/songe-converter/directory"
	"github.com/lolPants/songe-converter/utils"
)

func convert(dir string, c chan result) {
	fail := func(err string) {
		e := errors.New(err)
		res := result{
			dir:     dir,
			oldHash: "",
			newHash: "",
			err:     e,
		}

		c <- res
	}

	base := filepath.Base(dir)
	if base == "info.json" {
		dir = filepath.Dir(dir)
	}

	dirType, _ := directory.ReadType(dir)
	if dirType != directory.Old {
		fail("\"" + dir + "\" does not contain an old format beatmap")
		return
	}

	old, err := converter.ReadDirectoryOld(dir)
	if err != nil {
		fail("could not load beatmap at \"" + dir + "\"")
		return
	}

	new, err := converter.OldToNew(old)
	if err != nil {
		fail("failed to convert beatmap at \"" + dir + "\"")
		return
	}

	if dryRun == false {
		newPath := filepath.Join(dir, "info.dat")
		newBytes, err := new.Bytes()
		if err != nil {
			fail("could not serialize \"" + newPath + "\"")
			return
		}

		ioutil.WriteFile(newPath, newBytes, 0644)
		if keepFiles == false {
			infoPath := filepath.Join(dir, "info.json")
			err := os.Remove(infoPath)

			if err != nil {
				fail("could not delete \"" + infoPath + "\"")
				return
			}
		}

		oldAudioPath := filepath.Join(dir, new.OldSongFilename)
		newAudioPath := filepath.Join(dir, new.SongFilename)

		audioChanged := oldAudioPath != newAudioPath
		if audioChanged == true && keepFiles == true {
			_, err := utils.CopyFile(oldAudioPath, newAudioPath)
			if err != nil {
				fail("could not copy \"" + oldAudioPath + "\"")
				return
			}
		} else if audioChanged == true && keepFiles == false {
			err := os.Rename(oldAudioPath, newAudioPath)
			if err != nil {
				fail("could not rename \"" + oldAudioPath + "\"")
				return
			}
		}

		for _, set := range new.DifficultyBeatmapSets {
			for _, diff := range set.DifficultyBeatmaps {
				diffPath := filepath.Join(dir, diff.BeatmapFilename)
				diffBytes, err := new.Bytes()
				if err != nil {
					fail("could not serialize \"" + diffPath + "\"")
					return
				}

				ioutil.WriteFile(diffPath, diffBytes, 0644)
				if keepFiles == false {
					oldDiffName := strings.Replace(diff.BeatmapFilename, ".dat", ".json", -1)
					oldDiffPath := filepath.Join(dir, oldDiffName)
					err := os.Remove(oldDiffPath)

					if err != nil {
						fail("could not delete \"" + oldDiffPath + "\"")
						return
					}
				}
			}
		}
	}

	res := result{
		dir:     dir,
		oldHash: old.Hash,
		newHash: new.Hash,
		err:     nil,
	}

	c <- res
}
