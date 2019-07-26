package types

import "github.com/lolPants/songe-converter/json"

// NewDifficultyJSON is the new beatmap difficulty file
type NewDifficultyJSON struct {
	Version string `json:"_version"`

	BPMChanges []BPMChange `json:"_BPMChanges"`
	Events     []Event     `json:"_events"`
	Notes      []Note      `json:"_notes"`
	Obstacles  []Obstacle  `json:"_obstacles"`
	Bookmarks  []Bookmark  `json:"_bookmarks"`
}

// Bytes Convert to byte array
func (i NewDifficultyJSON) Bytes() ([]byte, error) {
	return json.Marshal(i)
}
