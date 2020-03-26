package types

import "jackbaron.com/songe-converter/v2/pkg/json"

// NewDifficultyJSON is the new beatmap difficulty file
type NewDifficultyJSON struct {
	Version string `json:"_version"`

	Events    []Event    `json:"_events"`
	Notes     []Note     `json:"_notes"`
	Obstacles []Obstacle `json:"_obstacles"`

	CustomData NewDifficultyCustomData `json:"_customData,omitempty"`
}

// NewDifficultyCustomData is the custom data of new beatmap difficulties
type NewDifficultyCustomData struct {
	BPMChanges []BPMChange `json:"_BPMChanges,omitempty"`
	Bookmarks  []Bookmark  `json:"_bookmarks,omitempty"`

	Time int `json:"_time,omitempty"`
}

// Bytes Convert to byte array
func (i NewDifficultyJSON) Bytes() ([]byte, error) {
	return json.Marshal(i)
}
