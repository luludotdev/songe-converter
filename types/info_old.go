package types

import (
	"jackbaron.com/songe-converter/v2/json"
)

// OldInfoJSON is the old beatmap info file
type OldInfoJSON struct {
	Hash         string `json:"-"`
	AudioPath    string `json:"-"`
	NewAudioPath string `json:"-"`

	SongName    string `json:"songName"`
	SongSubName string `json:"songSubName"`
	AuthorName  string `json:"authorName"`

	Contributors []OldContributor `json:"contributors"`

	BeatsPerMinute        float64 `json:"beatsPerMinute"`
	PreviewStartTime      float64 `json:"previewStartTime"`
	PreviewDuration       float64 `json:"previewDuration"`
	CoverImagePath        string  `json:"coverImagePath"`
	EnvironmentName       string  `json:"environmentName"`
	OneSaber              bool    `json:"oneSaber"`
	CustomEnvironment     string  `json:"customEnvironment"`
	CustomEnvironmentHash string  `json:"customEnvironmentHash"`

	DifficultyLevels []DifficultyLevel `json:"difficultyLevels"`
}

// Bytes Convert to byte array
func (i OldInfoJSON) Bytes() ([]byte, error) {
	return json.MarshalPretty(i)
}

// OldContributor Old contributor info
type OldContributor struct {
	Role     string `json:"role"`
	Name     string `json:"name"`
	IconPath string `json:"iconPath"`
}

// DifficultyLevel Old difficulty level
type DifficultyLevel struct {
	Difficulty      string `json:"difficulty"`
	DifficultyRank  int    `json:"difficultyRank"`
	AudioPath       string `json:"audioPath"`
	JSONPath        string `json:"jsonPath"`
	Offset          int    `json:"offset"`
	OldOffset       int    `json:"oldOffset"`
	ChromaToggle    string `json:"chromaToggle"`
	CustomColors    bool   `json:"customColors"`
	Characteristic  string `json:"characteristic"`
	DifficultyLabel string `json:"difficultyLabel"`

	DiffJSON *OldDifficultyJSON `json:"-"`
}
