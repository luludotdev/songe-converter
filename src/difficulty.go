package main

import (
	"encoding/json"
	"errors"

	"github.com/TomOnTime/utfutil"
)

// OldDifficultyJSON is the old beatmap difficulty file
type OldDifficultyJSON struct {
	Version string `json:"_version"`

	BeatsPerMinute float64 `json:"_beatsPerMinute"`
	BeatsPerBar    int     `json:"_beatsPerBar"`

	NoteJumpSpeed           float64 `json:"_noteJumpSpeed"`
	NoteJumpStartBeatOffset int     `json:"_noteJumpStartBeatOffset"`

	Shuffle       float64 `json:"_shuffle"`
	ShufflePeriod float64 `json:"_shufflePeriod"`

	ColorLeft  *BeatmapColor `json:"_colorLeft,omitempty"`
	ColorRight *BeatmapColor `json:"_colorRight,omitempty"`

	Time int `json:"_time"`

	Warnings     []string `json:"_warnings"`
	Information  []string `json:"_information"`
	Suggestions  []string `json:"_suggestions"`
	Requirements []string `json:"_requirements"`

	BPMChanges []BPMChange `json:"_BPMChanges"`
	Events     []Event     `json:"_events"`
	Notes      []Note      `json:"_notes"`
	Obstacles  []Obstacle  `json:"_obstacles"`
	Bookmarks  []Bookmark  `json:"_bookmarks"`
}

// NewDifficultyJSON is the new beatmap difficulty file
type NewDifficultyJSON struct {
	Version string `json:"_version"`

	BPMChanges []BPMChange `json:"_BPMChanges"`
	Events     []Event     `json:"_events"`
	Notes      []Note      `json:"_notes"`
	Obstacles  []Obstacle  `json:"_obstacles"`
	Bookmarks  []Bookmark  `json:"_bookmarks"`
}

// BPMChange MM BPM Change
type BPMChange struct {
	BPM             float64 `json:"_BPM"`
	Time            float64 `json:"_time"`
	BeatsPerBar     int     `json:"_beatsPerBar"`
	MetronomeOffset int     `json:"_metronomeOffset"`
}

// Event Beatmap Event
type Event struct {
	Time  float64 `json:"_time"`
	Type  int     `json:"_type"`
	Value int     `json:"_value"`
}

// Note Beatmap Note
type Note struct {
	Time         float64 `json:"_time"`
	LineIndex    int     `json:"_lineIndex"`
	LineLayer    int     `json:"_lineLayer"`
	Type         int     `json:"_type"`
	CutDirection int     `json:"_cutDirection"`
}

// Obstacle Beatmap Obstacle
type Obstacle struct {
	Time      float64 `json:"_time"`
	LineIndex int     `json:"_lineIndex"`
	Type      int     `json:"_type"`
	Duration  float64 `json:"_duration"`
	Width     int     `json:"_width"`
}

// Bookmark MM Bookmark
type Bookmark struct {
	Time float64 `json:"_time"`
	Name string  `json:"_name"`
}

func readDifficulty(path string) (OldDifficultyJSON, error) {
	bytes, err := utfutil.ReadFile(path, utfutil.UTF8)
	if err != nil {
		return OldDifficultyJSON{}, err
	}

	valid := IsJSON(bytes)
	if valid == false {
		invalidError := errors.New("Invalid difficulty file found at \"" + path + "\"")
		return OldDifficultyJSON{}, invalidError
	}

	var diffJSON OldDifficultyJSON
	json.Unmarshal(bytes, &diffJSON)

	return diffJSON, nil
}

func getRank(difficulty string) int {
	if difficulty == "Easy" {
		return 1
	} else if difficulty == "Normal" {
		return 3
	} else if difficulty == "Hard" {
		return 5
	} else if difficulty == "Expert" {
		return 7
	} else if difficulty == "ExpertPlus" {
		return 9
	} else {
		return 3
	}
}
