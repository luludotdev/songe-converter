package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	"github.com/ttacon/chalk"
)

var (
	sha1ver = "unknown"
	gitTag  string

	printVer  bool
	printHelp bool

	concurrency uint = 5
	output           = ""
	glob             = ""
	allDirs          = false
	keepFiles        = false
	dryRun           = false
	quiet            = false

	args []string
)

func main() {
	SetDetails("Songe Converter "+gitTag, "https://github.com/lolPants/songe-converter")
	SetExample("./songe-converter -g '**/info.json' ./CustomSongs")

	RegisterBoolFlag(&printVer, "v", "version", "Print version information.")
	RegisterUintFlag(&concurrency, "c", "concurrency", "Max number of jobs allowed to run at a time.")
	RegisterStringFlag(&output, "o", "output", "Save converted hashes and errors to file.")
	RegisterStringFlag(&glob, "g", "glob", "Use a glob to match directories.")
	RegisterBoolFlag(&allDirs, "a", "all-dirs", "Run on all subfolders of given directory.")
	RegisterBoolFlag(&keepFiles, "k", "keep-original", "Do not delete original JSON files")
	RegisterBoolFlag(&dryRun, "d", "dry-run", "Do not modify filesystem, only log output.")
	RegisterBoolFlag(&quiet, "q", "quiet", "Do not print to stdout.")
	ParseFlags(&args)

	if len(os.Args[1:]) == 0 {
		PrintUsageAndExit()
		return
	}

	if printVer == true {
		if gitTag != "" {
			fmt.Println(gitTag)
		}

		fmt.Println(sha1ver)
		return
	}

	if printHelp == true {
		PrintUsageAndExit()
		return
	}

	if concurrency < 1 {
		printError(chalk.Bold.TextStyle("--concurrency") + " cannot be less than 1!")
	}

	if allDirs == true && glob != "" {
		printError(chalk.Bold.TextStyle("--all-dirs") + " and " +
			chalk.Bold.TextStyle("--glob") + " cannot be used together!")
	}

	dirs := make([]string, 0)
	if allDirs == true || glob != "" {
		var dir string

		if len(args) > 0 && args[0] != "" {
			dir = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				printError(err.Error())
			}

			dir = cwd
		}

		if allDirs == true {
			fileInfo, err := ioutil.ReadDir(dir)
			if err != nil {
				printError("Could not list subdirectories of \"" + dir + "\"")
			}

			for _, file := range fileInfo {
				if file.IsDir() {
					subDir := filepath.Join(dir, file.Name())
					dirs = append(dirs, subDir)
				}
			}
		} else if glob != "" {
			pattern := filepath.Join(dir, glob)

			paths, err := doublestar.Glob(pattern)
			if err != nil {
				printError("Error matching glob path!")
			}

			dirs = paths
		}
	} else {
		dirs = args
	}

	flags := CommandFlags{keepFiles, dryRun, quiet}
	c := make(chan Result, len(dirs))

	if flags.dryRun && !flags.quiet {
		printWarning("Performing a dry run!")
	}

	sem := make(chan bool, concurrency)
	for _, dir := range dirs {
		sem <- true
		go func(dir string) {
			defer func() { <-sem }()
			run(dir, flags, c)
		}(dir)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	results := make([]Result, 0)
	for i := 0; i < len(dirs); i++ {
		result := <-c
		results = append(results, result)
	}

	if output != "" {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			f.Close()
			printError(err.Error())
		}

		defer f.Close()

		f.WriteString("HASHES:\n")
		for _, result := range results {
			if result.err != nil {
				continue
			}

			f.WriteString(result.oldHash)
			f.WriteString("\t")
			f.WriteString(result.newHash)
			f.WriteString("\t")
			f.WriteString(result.dir)
			f.WriteString("\n")
		}

		f.WriteString("\nERRORS:\n")
		for _, result := range results {
			if result.err == nil {
				continue
			}

			f.WriteString(result.err.Error())
			f.WriteString("\t")
			f.WriteString(result.dir)
			f.WriteString("\n")
		}
	}
}

// CommandFlags Command Flags
type CommandFlags struct {
	keepFiles bool
	dryRun    bool
	quiet     bool
}

// Result Converted Hashes
type Result struct {
	dir string

	oldHash string
	newHash string
	err     error
}
