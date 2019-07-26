package flags

import (
	"fmt"
	"os"
	"sort"
)

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

// PrintUsageAndExit prints CLI flags usage and exits with code 1
func PrintUsageAndExit() {
	PrintUsage()
	os.Exit(1)
}

// PrintUsageAndExitWithCode prints CLI flags usage and exits with code
func PrintUsageAndExitWithCode(code int) {
	PrintUsage()
	os.Exit(code)
}
