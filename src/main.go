package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
)

var sha1ver string

func registerStringFlag(p *string, name string, alias string, def string, usage string) {
	flag.StringVar(p, name, def, usage)
	flag.StringVar(p, alias, def, usage+" (short)")
}

func registerBoolFlag(p *bool, name string, alias string, def bool, usage string) {
	flag.BoolVar(p, name, def, usage)
	flag.BoolVar(p, alias, def, usage+" (alias of -"+name+")")
}

func registerIntFlag(p *int, name string, alias string, def int, usage string) {
	flag.IntVar(p, name, def, usage)
	flag.IntVar(p, alias, def, usage+" (alias of -"+name+")")
}

func main() {
	var (
		printVersion bool
		concurrency  int
		output       string
		glob         string
		allDirs      bool
		keepFiles    bool
		dryRun       bool
		quiet        bool
	)

	registerBoolFlag(&printVersion, "version", "v", false, "print version information")
	registerIntFlag(&concurrency, "concurrency", "c", 5, "max number of jobs allowed to run at a time")
	registerBoolFlag(&allDirs, "all-dirs", "a", false, "run on all subfolders of given directory")
	registerStringFlag(&glob, "glob", "g", "", "run a glob match in a given directory")
	registerBoolFlag(&keepFiles, "keep-orig", "k", false, "do not delete original JSON files")
	registerStringFlag(&output, "output", "o", "", "save converted hashes and errors to file")
	registerBoolFlag(&dryRun, "dry-run", "d", false, "don't modify filesystem, only log output")
	registerBoolFlag(&quiet, "quiet", "q", false, "don't print to stdout")

	if len(os.Args[1:]) == 0 {
		fmt.Print("songe converter -- by lolPants\n\nflags:\n")
		flag.PrintDefaults()
		return
	}

	flag.Parse()

	if printVersion == true {
		if sha1ver == "" {
			sha1ver = "unknown"
		}

		fmt.Println(sha1ver)
		return
	}

	if concurrency < 1 {
		fatalStr("--concurrency cannot be less than 1")
	}

	if allDirs == true && glob != "" {
		fatalStr("--all-dirs and --glob cannot be used together!")
	}

	dirs := make([]string, 0)
	if allDirs == true || glob != "" {
		var dir string

		args := flag.Args()
		if len(args) > 0 && args[0] != "" {
			dir = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				fatal(err)
			}

			dir = cwd
		}

		if allDirs == true {
			fileInfo, err := ioutil.ReadDir(dir)
			if err != nil {
				fatalStr("Could not list subdirectories of \"" + dir + "\"")
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
				fatalStr("Error matching glob path!")
			}

			dirs = paths
		}
	} else {
		dirs = flag.Args()
	}

	flags := CommandFlags{keepFiles, dryRun, quiet}
	c := make(chan Result, len(dirs))

	if flags.dryRun && !flags.quiet {
		log.Print("WARNING: Performing a dry run!")
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
			fatal(err)
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
