package flags

// SetDetails set CLI details line
func SetDetails(name string, link string) {
	cliName = name
	cliLink = link
}

// SetExample set the CLI example
func SetExample(example string) {
	cliExample = example
}
