package utils

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
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

// RegisterStringFlag Registers a string flag
func RegisterStringFlag(reference *string, shortName string, longName string, description string) {
	registerUsage()

	defaultValue := *reference
	flag.StringVar(reference, longName, *reference, description)
	flag.StringVar(reference, shortName, *reference, description)

	var usage string
	if defaultValue != "" {
		usage = "[string]"
	} else {
		usage = "string"
	}

	f := betterFlag{usage, shortName, longName, description}
	flags = append(flags, f)
}

// RegisterBoolFlag Registers a boolean flag
func RegisterBoolFlag(reference *bool, shortName string, longName string, description string) {
	registerUsage()

	flag.BoolVar(reference, longName, *reference, description)
	flag.BoolVar(reference, shortName, *reference, description)
	usage := ""

	f := betterFlag{usage, shortName, longName, description}
	flags = append(flags, f)
}

// RegisterUintFlag Registers a uint flag
func RegisterUintFlag(reference *uint, shortName string, longName string, description string) {
	registerUsage()

	defaultValue := *reference
	flag.UintVar(reference, longName, *reference, description)
	flag.UintVar(reference, shortName, *reference, description)

	var usage string
	if defaultValue != 0 {
		usage = "[uint]"
	} else {
		usage = "uint"
	}

	f := betterFlag{usage, shortName, longName, description}
	flags = append(flags, f)
}

// SetDetails Set CLI details
func SetDetails(name string, link string) {
	cliName = name
	cliLink = link
}

// SetExample Set CLI example command
func SetExample(example string) {
	cliExample = example
}

// ParseFlags Parses flags and returns extra args
func ParseFlags(args *[]string) {
	if hasParsed == true {
		return
	}

	hasParsed = true
	flag.Parse()

	flagArgs := flag.Args()
	reflect.ValueOf(args).Elem().Set(reflect.ValueOf(flagArgs))
}

func registerUsage() {
	if usageRegistered == true {
		return
	}

	usageRegistered = true
	flag.Usage = func() {
		PrintUsage()
	}

	helpFlag := betterFlag{
		shortName:   "h",
		longName:    "help",
		description: "Prints this help information."}
	flags = append(flags, helpFlag)
}

// PrintUsage prints CLI flags usage
func PrintUsage() {
	fmt.Println(cliName)
	fmt.Println(cliLink)
	fmt.Print("\n")

	var (
		maxShortName int
		maxLongName  int
	)

	for _, f := range flags {
		if len(f.shortName) > maxShortName {
			maxShortName = len(f.shortName)
		}

		combined := len(f.longName + " " + f.usage)
		if combined > maxLongName {
			maxLongName = combined
		}
	}

	fmt.Println("  Flags:")
	sort.Slice(flags, func(i, j int) bool {
		return flags[i].shortName < flags[j].shortName
	})

	for _, f := range flags {
		fmt.Printf("    -%-*v", maxShortName+1, f.shortName)
		fmt.Printf("--%-*v", maxLongName+4, f.longName+" "+f.usage)
		fmt.Println(f.description)
	}

	if cliExample != "" {
		fmt.Println("\n  Example:")
		fmt.Printf("    %v\n", cliExample)
	}
}

// PrintUsageAndExit prints CLI flags usage and exits
func PrintUsageAndExit() {
	PrintUsage()
	os.Exit(1)
}
