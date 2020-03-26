package types

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
