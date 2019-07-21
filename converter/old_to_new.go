package converter

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lolPants/songe-converter/directory"
	"github.com/lolPants/songe-converter/types"
	"github.com/lolPants/songe-converter/utils"
)

// OldToNew converts the old format to the new format
func OldToNew(old *types.OldInfoJSON) (*types.NewInfoJSON, error) {
	var new types.NewInfoJSON

	new.Version = "2.0.0"
	new.SongName = old.SongName
	new.SongSubName = ""
	new.SongAuthorName = old.SongSubName
	new.LevelAuthorName = old.AuthorName

	new.CustomData.Contributors = make([]types.Contributor, 0)
	for _, c := range old.Contributors {
		contributor := types.Contributor{Role: c.Role, Name: c.Name, IconPath: c.IconPath}
		new.CustomData.Contributors = append(new.CustomData.Contributors, contributor)
	}

	new.BeatsPerMinute = old.BeatsPerMinute
	new.SongTimeOffset = 0

	new.PreviewStartTime = old.PreviewStartTime
	new.PreviewDuration = old.PreviewDuration

	new.CoverImageFilename = old.CoverImagePath

	new.EnvironmentName = old.EnvironmentName
	new.CustomData.CustomEnvironment = old.CustomEnvironment
	new.CustomData.CustomEnvironmentHash = old.CustomEnvironmentHash

	for _, diff := range old.DifficultyLevels {
		newFileName := strings.Replace(diff.JSONPath, ".json", ".dat", -1)

		var characteristic string
		if old.OneSaber {
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

		var beatmapSet types.DifficultyBeatmapSet
		beatmapSetIdx := -1

		for i := range new.DifficultyBeatmapSets {
			if new.DifficultyBeatmapSets[i].BeatmapCharacteristicName == characteristic {
				beatmapSet = new.DifficultyBeatmapSets[i]
				beatmapSetIdx = i
				break
			}
		}

		if beatmapSetIdx == -1 {
			beatmapSet.BeatmapCharacteristicName = characteristic
			beatmapSet.DifficultyBeatmaps = make([]types.DifficultyBeatmap, 0)

			new.DifficultyBeatmapSets = append(new.DifficultyBeatmapSets, beatmapSet)
			beatmapSetIdx = len(new.DifficultyBeatmapSets) - 1
		}

		var difficulty types.DifficultyBeatmap
		difficulty.Difficulty = diff.Difficulty
		difficulty.DifficultyRank = utils.GetDifficultyRank(diff.Difficulty)
		difficulty.CustomData.DifficultyLabel = diff.DifficultyLabel
		difficulty.BeatmapFilename = newFileName
		difficulty.NoteJumpMovementSpeed = diff.DiffJSON.NoteJumpSpeed
		difficulty.NoteJumpStartBeatOffset = diff.DiffJSON.NoteJumpStartBeatOffset
		difficulty.CustomData.EditorOffset = diff.Offset
		difficulty.CustomData.EditorOldOffset = diff.OldOffset
		difficulty.CustomData.Warnings = diff.DiffJSON.Warnings
		difficulty.CustomData.Information = diff.DiffJSON.Information
		difficulty.CustomData.Suggestions = diff.DiffJSON.Suggestions
		difficulty.CustomData.Requirements = diff.DiffJSON.Requirements

		if difficulty.CustomData.Warnings == nil {
			difficulty.CustomData.Warnings = make([]string, 0)
		}

		if difficulty.CustomData.Information == nil {
			difficulty.CustomData.Information = make([]string, 0)
		}

		if difficulty.CustomData.Suggestions == nil {
			difficulty.CustomData.Suggestions = make([]string, 0)
		}

		if difficulty.CustomData.Requirements == nil {
			difficulty.CustomData.Requirements = make([]string, 0)
		}

		needsMapExt := utils.NeedsMappingExtensions(diff.DiffJSON)
		hasMapExt := utils.StringInSlice("Mapping Extensions", difficulty.CustomData.Requirements)
		if needsMapExt == true && hasMapExt == false {
			difficulty.CustomData.Requirements = append(difficulty.CustomData.Requirements, "Mapping Extensions")
		}

		difficulty.CustomData.ColorLeft = diff.DiffJSON.ColorLeft
		difficulty.CustomData.ColorRight = diff.DiffJSON.ColorRight

		new.Shuffle = diff.DiffJSON.Shuffle
		new.ShufflePeriod = diff.DiffJSON.ShufflePeriod
		new.SongFilename = diff.AudioPath

		if diff.DiffJSON.BeatsPerMinute != 0 {
			new.BeatsPerMinute = diff.DiffJSON.BeatsPerMinute
		}

		var newDiffJSON types.NewDifficultyJSON
		newDiffJSON.Version = "2.0.0"

		newDiffJSON.BPMChanges = diff.DiffJSON.BPMChanges
		newDiffJSON.Events = diff.DiffJSON.Events
		newDiffJSON.Notes = diff.DiffJSON.Notes
		newDiffJSON.Obstacles = diff.DiffJSON.Obstacles
		newDiffJSON.Bookmarks = diff.DiffJSON.Bookmarks

		if newDiffJSON.BPMChanges == nil {
			newDiffJSON.BPMChanges = make([]types.BPMChange, 0)
		}

		if newDiffJSON.Events == nil {
			newDiffJSON.Events = make([]types.Event, 0)
		}

		if newDiffJSON.Notes == nil {
			newDiffJSON.Notes = make([]types.Note, 0)
		}

		if newDiffJSON.Obstacles == nil {
			newDiffJSON.Obstacles = make([]types.Obstacle, 0)
		}

		if newDiffJSON.Bookmarks == nil {
			newDiffJSON.Bookmarks = make([]types.Bookmark, 0)
		}

		difficulty.DiffJSON = &newDiffJSON
		beatmapSet.DifficultyBeatmaps = append(beatmapSet.DifficultyBeatmaps, difficulty)
		new.DifficultyBeatmapSets[beatmapSetIdx] = beatmapSet
	}

	for _, set := range new.DifficultyBeatmapSets {
		sort.Slice(set.DifficultyBeatmaps, func(i, j int) bool {
			return set.DifficultyBeatmaps[i].DifficultyRank < set.DifficultyBeatmaps[j].DifficultyRank
		})
	}

	new.OldSongFilename = new.SongFilename
	new.SongFilename = strings.Replace(new.SongFilename, ".ogg", ".egg", -1)

	allBytes := make([]byte, 0)
	infoBytes, err := new.Bytes()
	if err != nil {
		return nil, err
	}

	allBytes = append(allBytes, infoBytes...)

	for _, set := range new.DifficultyBeatmapSets {
		for _, d := range set.DifficultyBeatmaps {
			diffBytes, err := d.DiffJSON.Bytes()
			if err != nil {
				return nil, err
			}

			allBytes = append(allBytes, diffBytes...)
		}
	}

	new.Hash = utils.CalculateSHA1(allBytes)
	return &new, nil
}

// DirOldToNew converts an old directory into a new directory
func DirOldToNew(dir string, dryRun bool, keepFiles bool) Result {
	r := Result{Directory: dir}

	base := filepath.Base(dir)
	if base == "info.json" {
		dir = filepath.Dir(dir)
	}

	dirType, _ := directory.ReadType(dir)
	if dirType != directory.Old {
		r.Error = errors.New("\"" + dir + "\" does not contain an old format beatmap")
		return r
	}

	old, err := ReadDirectoryOld(dir)
	if err != nil {
		r.Error = errors.New("could not load beatmap at \"" + dir + "\"")
		return r
	}

	new, err := OldToNew(old)
	if err != nil {
		r.Error = errors.New("failed to convert beatmap at \"" + dir + "\"")
		return r
	}

	r.OldHash = old.Hash
	r.NewHash = new.Hash

	if dryRun == false {
		newPath := filepath.Join(dir, "info.dat")
		newBytes, err := new.Bytes()
		if err != nil {
			r.Error = errors.New("could not serialize \"" + newPath + "\"")
			return r
		}

		ioutil.WriteFile(newPath, newBytes, 0644)
		if keepFiles == false {
			infoPath := filepath.Join(dir, "info.json")
			err := os.Remove(infoPath)

			if err != nil {
				r.Error = errors.New("could not delete \"" + infoPath + "\"")
				return r
			}
		}

		oldAudioPath := filepath.Join(dir, new.OldSongFilename)
		newAudioPath := filepath.Join(dir, new.SongFilename)

		audioChanged := oldAudioPath != newAudioPath
		if audioChanged == true && keepFiles == true {
			_, err := utils.CopyFile(oldAudioPath, newAudioPath)
			if err != nil {
				r.Error = errors.New("could not copy \"" + oldAudioPath + "\"")
				return r
			}
		} else if audioChanged == true && keepFiles == false {
			err := os.Rename(oldAudioPath, newAudioPath)
			if err != nil {
				r.Error = errors.New("could not rename \"" + oldAudioPath + "\"")
				return r
			}
		}

		for _, set := range new.DifficultyBeatmapSets {
			for _, diff := range set.DifficultyBeatmaps {
				diffPath := filepath.Join(dir, diff.BeatmapFilename)
				diffBytes, err := diff.DiffJSON.Bytes()
				if err != nil {
					r.Error = errors.New("could not serialize \"" + diffPath + "\"")
					return r
				}

				ioutil.WriteFile(diffPath, diffBytes, 0644)
				if keepFiles == false {
					oldDiffName := strings.Replace(diff.BeatmapFilename, ".dat", ".json", -1)
					oldDiffPath := filepath.Join(dir, oldDiffName)
					err := os.Remove(oldDiffPath)

					if err != nil {
						r.Error = errors.New("could not delete \"" + oldDiffPath + "\"")
						return r
					}
				}
			}
		}
	}

	return r
}
