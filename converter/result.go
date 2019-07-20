package converter

// Result conversion result
type Result struct {
	OldHash string
	NewHash string

	Directory string
	Error     error
}
