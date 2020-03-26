package converter

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"jackbaron.com/songe-converter/v2/pkg/directory"
	"jackbaron.com/songe-converter/v2/pkg/types"
	"jackbaron.com/songe-converter/v2/pkg/utils"
)

// NewToOld converts the new format back to the old format
func NewToOld(new *types.NewInfoJSON) (*types.OldInfoJSON, error) {
	var old types.OldInfoJSON

	old.SongName = new.SongName
	old.SongSubName = new.SongAuthorName
	old.AuthorName = new.LevelAuthorName

	old.Contributors = make([]types.OldContributor, 0)
	for _, c := range new.CustomData.Contributors {
		contributor := types.OldContributor{Role: c.Role, Name: c.Name, IconPath: c.IconPath}
		old.Contributors = append(old.Contributors, contributor)
	}

	old.BeatsPerMinute = new.BeatsPerMinute
	old.PreviewStartTime = new.PreviewStartTime
	old.PreviewDuration = new.PreviewDuration
	old.CoverImagePath = new.CoverImageFilename
	old.NewAudioPath = new.SongFilename

	old.EnvironmentName = new.EnvironmentName
	old.CustomEnvironment = new.CustomData.CustomEnvironment
	old.CustomEnvironmentHash = new.CustomData.CustomEnvironmentHash

	allBytes := make([]byte, 0)

	for _, set := range new.DifficultyBeatmapSets {
		for _, diff := range set.DifficultyBeatmaps {
			var level types.DifficultyLevel

			level.Difficulty = diff.Difficulty
			level.DifficultyRank = utils.GetOldDifficultyRank(diff.DifficultyRank)
			level.AudioPath = strings.Replace(new.SongFilename, ".egg", ".ogg", -1)
			old.AudioPath = level.AudioPath
			level.JSONPath = strings.Replace(diff.BeatmapFilename, ".dat", ".json", -1)
			level.Offset = diff.CustomData.EditorOffset
			level.OldOffset = diff.CustomData.EditorOldOffset
			level.DifficultyLabel = diff.CustomData.DifficultyLabel

			if set.BeatmapCharacteristicName == "OneSaber" {
				level.Characteristic = "One Saber"
			} else if set.BeatmapCharacteristicName == "NoArrows" {
				level.Characteristic = "No Arrows"
			} else {
				level.Characteristic = set.BeatmapCharacteristicName
			}

			var diffJSON types.OldDifficultyJSON
			diffJSON.Version = "1.5.0"

			diffJSON.ColorLeft = diff.CustomData.ColorLeft
			diffJSON.ColorRight = diff.CustomData.ColorRight

			diffJSON.Shuffle = new.Shuffle
			diffJSON.ShufflePeriod = new.ShufflePeriod
			diffJSON.BeatsPerMinute = new.BeatsPerMinute
			diffJSON.NoteJumpSpeed = diff.NoteJumpMovementSpeed
			diffJSON.NoteJumpStartBeatOffset = diff.NoteJumpStartBeatOffset

			diffJSON.Warnings = diff.CustomData.Warnings
			diffJSON.Information = diff.CustomData.Information
			diffJSON.Suggestions = diff.CustomData.Suggestions
			diffJSON.Requirements = diff.CustomData.Requirements

			if diffJSON.Warnings == nil {
				diffJSON.Warnings = make([]string, 0)
			}

			if diffJSON.Information == nil {
				diffJSON.Information = make([]string, 0)
			}

			if diffJSON.Suggestions == nil {
				diffJSON.Suggestions = make([]string, 0)
			}

			if diffJSON.Requirements == nil {
				diffJSON.Requirements = make([]string, 0)
			}

			diffJSON.BPMChanges = diff.DiffJSON.BPMChanges
			diffJSON.Events = diff.DiffJSON.Events
			diffJSON.Notes = diff.DiffJSON.Notes
			diffJSON.Obstacles = diff.DiffJSON.Obstacles
			diffJSON.Bookmarks = diff.DiffJSON.Bookmarks

			if diffJSON.BPMChanges == nil {
				diffJSON.BPMChanges = make([]types.BPMChange, 0)
			}

			if diffJSON.Events == nil {
				diffJSON.Events = make([]types.Event, 0)
			}

			if diffJSON.Notes == nil {
				diffJSON.Notes = make([]types.Note, 0)
			}

			if diffJSON.Obstacles == nil {
				diffJSON.Obstacles = make([]types.Obstacle, 0)
			}

			if diffJSON.Bookmarks == nil {
				diffJSON.Bookmarks = make([]types.Bookmark, 0)
			}

			level.DiffJSON = &diffJSON
			diffBytes, err := diffJSON.Bytes()
			if err != nil {
				return nil, err
			}

			allBytes = append(allBytes, diffBytes...)
			old.DifficultyLevels = append(old.DifficultyLevels, level)
		}
	}

	old.Hash = utils.CalculateMD5(allBytes)
	return &old, nil
}

// DirNewToOld converts a new directory into an old directory
func DirNewToOld(dir string, dryRun bool, keepFiles bool) Result {
	r := Result{Directory: dir}

	base := filepath.Base(dir)
	if base == "info.json" {
		dir = filepath.Dir(dir)
	}

	dirType, _ := directory.ReadType(dir)
	if dirType != directory.New {
		r.Error = errors.New("\"" + dir + "\" does not contain a new format beatmap")
		return r
	}

	new, err := ReadDirectoryNew(dir)
	if err != nil {
		r.Error = errors.New("could not load beatmap at \"" + dir + "\"")
		return r
	}

	old, err := NewToOld(new)
	if err != nil {
		r.Error = errors.New("failed to convert beatmap at \"" + dir + "\"")
		return r
	}

	r.NewHash = new.Hash
	r.OldHash = old.Hash

	if dryRun == false {
		oldPath := filepath.Join(dir, "info.json")
		oldBytes, err := old.Bytes()
		if err != nil {
			r.Error = errors.New("could not serialize \"" + oldPath + "\"")
			return r
		}

		ioutil.WriteFile(oldPath, oldBytes, 0644)
		if keepFiles == false {
			infoPath := filepath.Join(dir, "Info.dat")
			err := os.Remove(infoPath)

			if err != nil {
				r.Error = errors.New("could not delete \"" + infoPath + "\"")
				return r
			}
		}

		newAudioPath := filepath.Join(dir, old.NewAudioPath)
		oldAudioPath := filepath.Join(dir, old.AudioPath)

		audioChanged := oldAudioPath != newAudioPath
		if audioChanged == true && keepFiles == true {
			_, err := utils.CopyFile(newAudioPath, oldAudioPath)
			if err != nil {
				r.Error = errors.New("could not copy \"" + newAudioPath + "\"")
				return r
			}
		} else if audioChanged == true && keepFiles == false {
			err := os.Rename(newAudioPath, oldAudioPath)
			if err != nil {
				r.Error = errors.New("could not rename \"" + newAudioPath + "\"")
				return r
			}
		}

		for _, diff := range old.DifficultyLevels {
			diffPath := filepath.Join(dir, diff.JSONPath)
			diffBytes, err := diff.DiffJSON.Bytes()
			if err != nil {
				r.Error = errors.New("could not serialize \"" + diffPath + "\"")
				return r
			}

			ioutil.WriteFile(diffPath, diffBytes, 0644)
			if keepFiles == false {
				newDiffName := strings.Replace(diff.JSONPath, ".json", ".dat", -1)
				newDiffPath := filepath.Join(dir, newDiffName)
				err := os.Remove(newDiffPath)

				if err != nil {
					r.Error = errors.New("could not delete \"" + newDiffPath + "\"")
					return r
				}
			}
		}
	}

	return r
}
