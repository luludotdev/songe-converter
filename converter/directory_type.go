package converter

// DirectoryType can be "New", "Old", "Both", or "None"
type DirectoryType uint

const (
	// New beatmap format
	New DirectoryType = 0
	// Old beatmap format
	Old DirectoryType = 1
	// Both beatmap formats
	Both DirectoryType = 2
	// None No beatmap formats
	None DirectoryType = 3
)

func (dirType DirectoryType) String() string {
	names := [...]string{"new", "old", "both", "none"}
	if dirType < New || dirType > None {
		return "unknown"
	}

	return names[dirType]
}

// Valid Directory type can be used to convert automatically
func (dirType DirectoryType) Valid() bool {
	if dirType == New || dirType == Old {
		return true
	}

	return false
}
