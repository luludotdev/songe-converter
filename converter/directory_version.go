package converter

// DirectoryVersion can be "New", "Old", "Both", or "None"
type DirectoryVersion uint

const (
	// New beatmap format
	New DirectoryVersion = 0
	// Old beatmap format
	Old DirectoryVersion = 1
	// Both beatmap formats
	Both DirectoryVersion = 2
	// None No beatmap formats
	None DirectoryVersion = 3
)

func (version DirectoryVersion) String() string {
	names := [...]string{"new", "old", "both", "none"}
	if version < New || version > None {
		return "unknown"
	}

	return names[version]
}

// Valid Directory Version can be used to convert automatically
func (version DirectoryVersion) Valid() bool {
	if version == New || version == Old {
		return true
	}

	return false
}
