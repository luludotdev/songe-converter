package converter

import (
	"errors"
	"path/filepath"

	"github.com/TomOnTime/utfutil"
	"github.com/lolPants/songe-converter/types"
	"github.com/lolPants/songe-converter/utils"
)

// ReadDirectoryNew Reads a directory and loads new beatmap format
func ReadDirectoryNew(path string) (*types.NewInfoJSON, error) {
	dirType, _ := BeatmapDirectoryType(path)
	if dirType != New {
		return nil, errors.New("not a new format beatmap")
	}

	hashBytes := make([]byte, 0)

	infoPath := filepath.Join(path, "info.dat")
	bytes, err := utfutil.ReadFile(infoPath, utfutil.UTF8)
	if err != nil {
		return nil, err
	}

	hashBytes = append(hashBytes, bytes...)
	info, err := LoadNewInfo(bytes)
	if err != nil {
		return nil, err
	}

	for i, set := range info.DifficultyBeatmapSets {
		for j, mapDiff := range set.DifficultyBeatmaps {
			diffPath := filepath.Join(path, mapDiff.BeatmapFilename)
			diffBytes, err := utfutil.ReadFile(diffPath, utfutil.UTF8)
			if err != nil {
				return nil, err
			}

			hashBytes = append(hashBytes, diffBytes...)
			diffJSON, err := LoadNewDifficulty(diffBytes)
			if err != nil {
				return nil, err
			}

			info.DifficultyBeatmapSets[i].DifficultyBeatmaps[j].DiffJSON = diffJSON
		}
	}

	info.Hash = utils.CalculateSHA1(hashBytes)
	return info, nil
}

// ReadDirectoryOld Reads a directory and loads old beatmap format
func ReadDirectoryOld(path string) (*types.OldInfoJSON, error) {
	dirType, _ := BeatmapDirectoryType(path)
	if dirType != Old {
		return nil, errors.New("not an old format beatmap")
	}

	hashBytes := make([]byte, 0)

	infoPath := filepath.Join(path, "info.json")
	bytes, err := utfutil.ReadFile(infoPath, utfutil.UTF8)
	if err != nil {
		return nil, err
	}

	info, err := LoadOldInfo(bytes)
	if err != nil {
		return nil, err
	}

	for i, mapDiff := range info.DifficultyLevels {
		diffPath := filepath.Join(path, mapDiff.JSONPath)
		diffBytes, err := utfutil.ReadFile(diffPath, utfutil.UTF8)
		if err != nil {
			return nil, err
		}

		hashBytes = append(hashBytes, diffBytes...)
		diffJSON, err := LoadOldDifficulty(diffBytes)
		if err != nil {
			return nil, err
		}

		info.DifficultyLevels[i].DiffJSON = diffJSON
	}

	info.Hash = utils.CalculateMD5(hashBytes)
	return info, nil
}
