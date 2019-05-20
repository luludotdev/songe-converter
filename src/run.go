package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

func run(dir string, flags CommandFlags, c chan Result) {
	info := path.Join(dir, "info.json")
	infoJSON, infoErr := readInfo(info)
	if infoErr != nil && os.IsNotExist(infoErr) {
		log.Print("No info.json found in \"" + dir + "\", skipping!")

		result := Result{dir: dir, oldHash: "", newHash: "", err: errors.New("info.json not found")}
		c <- result
		return
	} else if infoErr != nil {
		log.Print("Something went wrong when converting \"" + dir + "\"")
		log.Print(infoErr)

		result := Result{dir: dir, oldHash: "", newHash: "", err: infoErr}
		c <- result
	} else {
		log.Print("Converting \"" + dir + "\"")
	}

	var newInfoJSON NewInfoJSON
	newInfoJSON.Version = "2.0.0"

	newInfoJSON.SongName = infoJSON.SongName
	newInfoJSON.SongSubName = infoJSON.SongSubName
	newInfoJSON.SongAuthorName = infoJSON.AuthorName
	newInfoJSON.LevelAuthorName = ""

	newInfoJSON.Contributors = make([]Contributor, 0)
	for _, c := range infoJSON.Contributors {
		contributor := Contributor{Role: c.Role, Name: c.Name, IconPath: c.IconPath}
		newInfoJSON.Contributors = append(newInfoJSON.Contributors, contributor)
	}

	newInfoJSON.BeatsPerMinute = infoJSON.BeatsPerMinute
	newInfoJSON.SongTimeOffset = 0

	newInfoJSON.PreviewStartTime = infoJSON.PreviewStartTime
	newInfoJSON.PreviewDuration = infoJSON.PreviewDuration

	newInfoJSON.CoverImageFilename = infoJSON.CoverImagePath

	newInfoJSON.EnvironmentName = infoJSON.EnvironmentName
	newInfoJSON.CustomEnvironment = infoJSON.CustomEnvironment
	newInfoJSON.CustomEnvironmentHash = infoJSON.CustomEnvironmentHash

	allBytes := make([]byte, 0)
	toDelete := make([]string, 0)

	newInfoJSON.DifficultyBeatmapSets = make([]DifficultyBeatmapSet, 0)
	for _, diff := range infoJSON.DifficultyLevels {
		// Read JSON
		json := path.Join(dir, diff.JSONPath)
		toDelete = append(toDelete, json)

		diffJSON, diffErr := readDifficulty(json)
		if diffErr != nil && os.IsNotExist(diffErr) {
			log.Print(diff.JSONPath + " not found in \"" + dir + "\", skipping!")

			result := Result{dir: dir, oldHash: "", newHash: "", err: errors.New(diff.JSONPath + " not found")}
			c <- result
			return
		} else if diffErr != nil {
			log.Print("Something went wrong when reading \"" + json + "\"")
			log.Print(diffErr)

			result := Result{dir: dir, oldHash: "", newHash: "", err: diffErr}
			c <- result
		}

		// New File Name
		diffJSONFileName := strings.Replace(diff.JSONPath, ".json", ".dat", -1)

		var characteristic string
		if infoJSON.OneSaber {
			characteristic = "OneSaber"
		} else if diff.Characteristic == "One Saber" {
			characteristic = "OneSaber"
		} else if diff.Characteristic == "No Arrows" {
			characteristic = "NoArrows"
		} else if diff.Characteristic != "" {
			characteristic = diff.Characteristic
		} else {
			characteristic = "Standard"
		}

		var beatmapSet DifficultyBeatmapSet
		beatmapSetIdx := -1
		for i := range newInfoJSON.DifficultyBeatmapSets {
			if newInfoJSON.DifficultyBeatmapSets[i].BeatmapCharacteristicName == characteristic {
				beatmapSet = newInfoJSON.DifficultyBeatmapSets[i]
				beatmapSetIdx = i
				break
			}
		}

		if beatmapSetIdx == -1 {
			beatmapSet.BeatmapCharacteristicName = characteristic
			beatmapSet.DifficultyBeatmaps = make([]DifficultyBeatmap, 0)

			newInfoJSON.DifficultyBeatmapSets = append(newInfoJSON.DifficultyBeatmapSets, beatmapSet)
			beatmapSetIdx = len(newInfoJSON.DifficultyBeatmapSets) - 1
		}

		var difficulty DifficultyBeatmap
		difficulty.Difficulty = diff.Difficulty
		difficulty.DifficultyRank = getRank(diff.Difficulty)
		difficulty.DifficultyLabel = diff.DifficultyLabel
		difficulty.BeatmapFilename = diffJSONFileName
		difficulty.NoteJumpMovementSpeed = diffJSON.NoteJumpSpeed
		difficulty.NoteJumpStartBeatOffset = diffJSON.NoteJumpStartBeatOffset
		difficulty.EditorOffset = diff.Offset
		difficulty.EditorOldOffset = diff.OldOffset
		difficulty.Warnings = diffJSON.Warnings
		difficulty.Information = diffJSON.Information
		difficulty.Suggestions = diffJSON.Suggestions
		difficulty.Requirements = diffJSON.Requirements

		if difficulty.Warnings == nil {
			difficulty.Warnings = make([]string, 0)
		}

		if difficulty.Information == nil {
			difficulty.Information = make([]string, 0)
		}

		if difficulty.Suggestions == nil {
			difficulty.Suggestions = make([]string, 0)
		}

		if difficulty.Requirements == nil {
			difficulty.Requirements = make([]string, 0)
		}

		difficulty.ColorLeft = diffJSON.ColorLeft
		difficulty.ColorRight = diffJSON.ColorRight

		beatmapSet.DifficultyBeatmaps = append(beatmapSet.DifficultyBeatmaps, difficulty)
		newInfoJSON.DifficultyBeatmapSets[beatmapSetIdx] = beatmapSet

		newInfoJSON.Shuffle = diffJSON.Shuffle
		newInfoJSON.ShufflePeriod = diffJSON.ShufflePeriod
		newInfoJSON.SongFilename = diff.AudioPath

		if diffJSON.BeatsPerMinute != 0 {
			infoJSON.BeatsPerMinute = diffJSON.BeatsPerMinute
		}

		var newDiffJSON NewDifficultyJSON
		newDiffJSON.Version = "2.0.0"

		// Set
		newDiffJSON.BPMChanges = diffJSON.BPMChanges
		newDiffJSON.Events = diffJSON.Events
		newDiffJSON.Notes = diffJSON.Notes
		newDiffJSON.Obstacles = diffJSON.Obstacles
		newDiffJSON.Bookmarks = diffJSON.Bookmarks

		if newDiffJSON.BPMChanges == nil {
			newDiffJSON.BPMChanges = make([]BPMChange, 0)
		}

		if newDiffJSON.Events == nil {
			newDiffJSON.Events = make([]Event, 0)
		}

		if newDiffJSON.Notes == nil {
			newDiffJSON.Notes = make([]Note, 0)
		}

		if newDiffJSON.Obstacles == nil {
			newDiffJSON.Obstacles = make([]Obstacle, 0)
		}

		if newDiffJSON.Bookmarks == nil {
			newDiffJSON.Bookmarks = make([]Bookmark, 0)
		}

		// Save
		diffJSONBytes, _ := JSONMarshal(newDiffJSON)
		allBytes = append(allBytes, diffJSONBytes...)
		diffJSONPath := path.Join(dir, diffJSONFileName)
		if flags.dryRun == false {
			_ = ioutil.WriteFile(diffJSONPath, diffJSONBytes, 0644)
		}
	}

	for _, set := range newInfoJSON.DifficultyBeatmapSets {
		sort.Slice(set.DifficultyBeatmaps, func(i, j int) bool {
			return set.DifficultyBeatmaps[i].DifficultyRank < set.DifficultyBeatmaps[j].DifficultyRank
		})
	}

	infoJSONBytes, _ := JSONMarshalPretty(newInfoJSON)
	allBytes = append(allBytes, infoJSONBytes...)
	infoJSONPath := path.Join(dir, "info.dat")
	if flags.dryRun == false {
		_ = ioutil.WriteFile(infoJSONPath, infoJSONBytes, 0644)
	}

	oldHash, err := calculateOldHash(infoJSON, dir)
	if err != nil {
		log.Print("Something went wrong when converting \"" + dir + "\"")
		log.Print(err)

		result := Result{dir: dir, oldHash: "", newHash: "", err: err}
		c <- result
	}

	if flags.keepFiles == false && flags.dryRun == false {
		err := os.Remove(info)
		if err != nil {
			log.Print("Failed to delete \"" + info + "\"")
		}

		for _, d := range toDelete {
			err := os.Remove(d)
			if err != nil {
				log.Print("Failed to delete \"" + d + "\"")
			}
		}
	}

	newHash := calculateHashSHA1(allBytes)
	result := Result{dir: dir, oldHash: oldHash, newHash: newHash, err: nil}
	c <- result
}
