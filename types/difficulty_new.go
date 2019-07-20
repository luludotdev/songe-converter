package types

// NewDifficultyJSON is the new beatmap difficulty file
type NewDifficultyJSON struct {
	Version string `json:"_version"`

	BPMChanges []BPMChange `json:"_BPMChanges"`
	Events     []Event     `json:"_events"`
	Notes      []Note      `json:"_notes"`
	Obstacles  []Obstacle  `json:"_obstacles"`
	Bookmarks  []Bookmark  `json:"_bookmarks"`
}
