package types

import "github.com/lolPants/songe-converter/json"

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

// Bytes Convert to byte array
func (i OldDifficultyJSON) Bytes() ([]byte, error) {
	return json.Marshal(i)
}
