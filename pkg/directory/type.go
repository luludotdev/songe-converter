package directory

// Type can be "New", "Old", "Both", or "None"
type Type uint

const (
	// New beatmap format
	New Type = 0
	// Old beatmap format
	Old Type = 1
	// Both beatmap formats
	Both Type = 2
	// None No beatmap formats
	None Type = 3
)

func (t Type) String() string {
	names := [...]string{"new", "old", "both", "none"}
	if t < New || t > None {
		return "unknown"
	}

	return names[t]
}

// Valid Directory type can be used to convert automatically
func (t Type) Valid() bool {
	if t == New || t == Old {
		return true
	}

	return false
}
