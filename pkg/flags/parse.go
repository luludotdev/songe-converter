package flags

import (
	"flag"
	"reflect"
)

type betterFlag struct {
	usage       string
	shortName   string
	longName    string
	description string
}

var (
	usageRegistered = false
	hasParsed       = false

	cliName    string
	cliLink    string
	cliExample string

	flags []betterFlag
)

// Parse parses flags and returns extra args
func Parse(args *[]string) {
	if hasParsed == true {
		return
	}

	hasParsed = true
	flag.Parse()

	flagArgs := flag.Args()
	reflect.ValueOf(args).Elem().Set(reflect.ValueOf(flagArgs))
}
