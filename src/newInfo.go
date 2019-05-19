package main

// NewInfoJSON New Info JSON
type NewInfoJSON struct {
	Version string `json:"_version"`

	SongName        string `json:"_songName"`
	SongSubName     string `json:"_songSubName"`
	SongAuthorName  string `json:"_songAuthorName"`
	LevelAuthorName string `json:"_levelAuthorName"`

	Contributors []Contributor `json:"_contributors"`

	BeatsPerMinute float64 `json:"_beatsPerMinute"`
	SongTimeOffset float64 `json:"_songTimeOffset"`
	Shuffle        float64 `json:"_shuffle"`
	ShufflePeriod  float64 `json:"_shufflePeriod"`

	PreviewStartTime float64 `json:"_previewStartTime"`
	PreviewDuration  float64 `json:"_previewDuration"`

	SongFilename       string `json:"_songFilename"`
	CoverImageFilename string `json:"_coverImageFilename"`

	EnvironmentName       string `json:"_environmentName"`
	CustomEnvironment     string `json:"_customEnvironment"`
	CustomEnvironmentHash string `json:"_customEnvironmentHash"`

	DifficultyBeatmapSets []DifficultyBeatmapSet `json:"_difficultyBeatmapSets"`
}

// Contributor New Info JSON Contributors
type Contributor struct {
	Role     string `json:"_role"`
	Name     string `json:"_name"`
	IconPath string `json:"_iconPath"`
}

// BeatmapColor Beatmap Lighting Color
type BeatmapColor struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

// DifficultyBeatmap Beatmap Difficulty Info
type DifficultyBeatmap struct {
	Difficulty      string `json:"_difficulty"`
	DifficultyRank  int    `json:"_difficultyRank"`
	DifficultyLabel string `json:"_difficultyLabel"`
	BeatmapFilename string `json:"_beatmapFilename"`

	NoteJumpMovementSpeed   float64 `json:"_noteJumpMovementSpeed"`
	NoteJumpStartBeatOffset int     `json:"_noteJumpStartBeatOffset"`

	EditorOffset    int `json:"_editorOffset"`
	EditorOldOffset int `json:"_editorOldOffset"`

	ColorLeft  *BeatmapColor `json:"_colorLeft"`
	ColorRight *BeatmapColor `json:"_colorRight"`

	Warnings     []string `json:"_warnings"`
	Information  []string `json:"_information"`
	Suggestions  []string `json:"_suggestions"`
	Requirements []string `json:"_requirements"`
}

// DifficultyBeatmapSet Set of DifficultyBeatmap structs
type DifficultyBeatmapSet struct {
	BeatmapCharacteristicName string              `json:"_beatmapCharacteristicName"`
	DifficultyBeatmaps        []DifficultyBeatmap `json:"_difficultyBeatmaps"`
}
