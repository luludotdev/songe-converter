package converter

import (
	"encoding/json"
	"errors"

	j "github.com/lolPants/songe-converter/json"
	"github.com/lolPants/songe-converter/types"
)

// LoadNewInfo Loads new info struct from byte array
func LoadNewInfo(bytes []byte) (*types.NewInfoJSON, error) {
	valid := json.Valid(bytes)
	if valid == false {
		return nil, errors.New("invalid info.dat")
	}

	var infoJSON types.NewInfoJSON
	err := json.Unmarshal(bytes, &infoJSON)
	if err != nil {
		return nil, err
	}

	return &infoJSON, nil
}

// LoadOldInfo Loads old info struct from byte array
func LoadOldInfo(bytes []byte) (*types.OldInfoJSON, error) {
	valid := j.Valid(bytes)
	if valid == false {
		return nil, errors.New("invalid info.json")
	}

	var infoJSON types.OldInfoJSON
	err := json.Unmarshal(bytes, &infoJSON)
	if err != nil {
		return nil, err
	}

	return &infoJSON, nil
}

// LoadNewDifficulty Loads new difficulty struct from byte array
func LoadNewDifficulty(bytes []byte) (*types.NewDifficultyJSON, error) {
	valid := j.Valid(bytes)
	if valid == false {
		return nil, errors.New("invalid difficulty.dat")
	}

	var diffJSON types.NewDifficultyJSON
	err := json.Unmarshal(bytes, &diffJSON)
	if err != nil {
		return nil, err
	}

	return &diffJSON, nil
}

// LoadOldDifficulty Loads old difficulty struct from byte array
func LoadOldDifficulty(bytes []byte) (*types.OldDifficultyJSON, error) {
	valid := j.Valid(bytes)
	if valid == false {
		return nil, errors.New("invalid difficulty.json")
	}

	var diffJSON types.OldDifficultyJSON
	err := json.Unmarshal(bytes, &diffJSON)
	if err != nil {
		return nil, err
	}

	return &diffJSON, nil
}
