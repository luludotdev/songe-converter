package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func registerStringFlag(p *string, name string, alias string, def string, usage string) {
	flag.StringVar(p, name, def, usage)
	flag.StringVar(p, alias, def, usage+" (short)")
}

func registerBoolFlag(p *bool, name string, alias string, def bool, usage string) {
	flag.BoolVar(p, name, def, usage)
	flag.BoolVar(p, alias, def, usage+" (alias of -"+name+")")
}

func main() {
	var (
		output    string
		allDirs   bool
		keepFiles bool
		dryRun    bool
	)

	registerStringFlag(&output, "output", "o", "", "save converted hashes and errors to file")
	registerBoolFlag(&allDirs, "all-dirs", "a", false, "run on all subfolders of given directory")
	registerBoolFlag(&keepFiles, "keep-orig", "k", false, "do not delete original JSON files")
	registerBoolFlag(&dryRun, "dry-run", "d", false, "don't modify filesystem, only log output")

	if len(os.Args[1:]) == 0 {
		fmt.Print("Songe Converter -- by lolPants\n\nflags:\n")
		flag.PrintDefaults()
		return
	}

	flag.Parse()

	dirs := make([]string, 0)
	if allDirs == true {
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
	} else {
		dirs = flag.Args()
	}

	flags := CommandFlags{keepFiles, dryRun}
	c := make(chan Result, len(dirs))

	if flags.dryRun {
		log.Print("WARNING: Performing a dry run!")
	}

	for _, dir := range dirs {
		go run(dir, flags, c)
	}

	results := make([]Result, 0)
	for i := 0; i < len(dirs); i++ {
		result := <-c
		results = append(results, result)
	}

	if output != "" {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
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
}

// Result Converted Hashes
type Result struct {
	dir string

	oldHash string
	newHash string
	err     error
}
