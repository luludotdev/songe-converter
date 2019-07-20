package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/briandowns/spinner"
	"github.com/lolPants/songe-converter/flags"
	"github.com/lolPants/songe-converter/log"
	"github.com/lolPants/songe-converter/utils"
	"github.com/ttacon/chalk"
)

var (
	sha1ver = "unknown"
	gitTag  string

	printVer  bool
	printHelp bool

	concurrency = uint(runtime.NumCPU())
	output      = ""
	outputErr   = ""
	glob        = ""
	allDirs     = false
	keepFiles   = false
	dryRun      = false

	args []string
)

func main() {
	flags.SetDetails("Songe Converter "+gitTag, "https://github.com/lolPants/songe-converter")
	flags.SetExample("./songe-converter -g '**/info.json' ./CustomSongs")

	flags.RegisterBoolFlag(&printVer, "v", "version", "Print version information.")
	flags.RegisterUintFlag(&concurrency, "c", "concurrency", "Max number of jobs allowed to run at a time.")
	flags.RegisterStringFlag(&output, "o", "output", "Save converted hashes to file.")
	flags.RegisterStringFlag(&outputErr, "e", "error-out", "Save conversion errors to file.")
	flags.RegisterStringFlag(&glob, "g", "glob", "Use a glob to match directories.")
	flags.RegisterBoolFlag(&allDirs, "a", "all-dirs", "Run on all subfolders of given directory.")
	flags.RegisterBoolFlag(&keepFiles, "k", "keep-original", "Do not delete original JSON files")
	flags.RegisterBoolFlag(&dryRun, "d", "dry-run", "Do not modify filesystem, only log output.")
	flags.Parse(&args)

	if len(os.Args[1:]) == 0 {
		flags.PrintUsageAndExit()
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
		flags.PrintUsageAndExit()
		return
	}

	if concurrency < 1 {
		log.Error(chalk.Bold.TextStyle("--concurrency") + " cannot be less than 1!")
	}

	if allDirs == true && glob != "" {
		log.Error(chalk.Bold.TextStyle("--all-dirs") + " and " +
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
				log.Error(err.Error())
			}

			dir = cwd
		}

		if allDirs == true {
			fileInfo, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Error("Could not list subdirectories of \"" + dir + "\"")
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
				log.Error("Error matching glob path!")
			}

			dirs = paths
		}
	} else {
		dirs = args
	}

	if dryRun {
		log.Warning("Performing a dry run!")
	}

	c := make(chan result, len(dirs))
	complete := make(chan bool, 1)

	var fout *os.File
	var ferr *os.File

	if output != "" {
		fout, err := utils.OpenFileSafe(output)
		if err != nil {
			log.Error(err.Error())
			if fout != nil {
				fout.Close()
			}
		}

		defer fout.Close()
		fout.Truncate(0)
		fout.Seek(0, 0)
	}

	if outputErr != "" {
		ferr, err := utils.OpenFileSafe(outputErr)
		if err != nil {
			log.Error(err.Error())
			if ferr != nil {
				ferr.Close()
			}
		}

		defer ferr.Close()
		ferr.Truncate(0)
		ferr.Seek(0, 0)
	}

	go func() {
		for i := 0; i < len(dirs); i++ {
			r := <-c

			if r.err == nil {
				if fout == nil {
					continue
				}

				fout.WriteString(r.oldHash)
				fout.WriteString("\t")
				fout.WriteString(r.newHash)
				fout.WriteString("\t")
				fout.WriteString(r.dir)
				fout.WriteString("\n")

				fout.Sync()
			} else {
				if ferr == nil {
					continue
				}

				ferr.WriteString(r.err.Error())
				ferr.WriteString("\t")
				ferr.WriteString(r.dir)
				ferr.WriteString("\n")

				ferr.Sync()
			}
		}

		complete <- true
	}()

	sem := make(chan bool, concurrency)
	for _, dir := range dirs {
		sem <- true
		go func(dir string) {
			defer func() { <-sem }()
			convert(dir, c)
		}(dir)
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Suffix = " Converting songes..."
	// s.Start()

	<-complete
	s.Stop()
}
