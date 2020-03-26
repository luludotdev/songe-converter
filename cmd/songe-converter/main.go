package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/briandowns/spinner"
	"github.com/lolPants/flaggs"
	"github.com/ttacon/chalk"

	"jackbaron.com/songe-converter/v2/pkg/converter"
	"jackbaron.com/songe-converter/v2/pkg/log"
	"jackbaron.com/songe-converter/v2/pkg/utils"
)

var (
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
	flaggs.SetDetails("Songe Converter "+gitTag, "https://github.com/lolPants/songe-converter")
	flaggs.SetExample("./songe-converter -g '**/info.json' ./CustomSongs")

	flaggs.RegisterBoolFlag(&printVer, "v", "version", "Print version information.")
	flaggs.RegisterUintFlag(&concurrency, "c", "concurrency", "Max number of jobs allowed to run at a time.")
	flaggs.RegisterStringFlag(&output, "o", "output", "Save converted hashes to file.")
	flaggs.RegisterStringFlag(&outputErr, "e", "error-out", "Save conversion errors to file.")
	flaggs.RegisterStringFlag(&glob, "g", "glob", "Use a glob to match directories.")
	flaggs.RegisterBoolFlag(&allDirs, "a", "all-dirs", "Run on all subfolders of given directory.")
	flaggs.RegisterBoolFlag(&keepFiles, "k", "keep-original", "Do not delete original JSON files")
	flaggs.RegisterBoolFlag(&dryRun, "d", "dry-run", "Do not modify filesystem, only log output.")
	flaggs.Parse(&args)

	if len(os.Args[1:]) == 0 {
		flaggs.PrintUsageAndExit()
		return
	}

	if printVer == true {
		printVersionInfo()
		os.Exit(0)
	}

	if printHelp == true {
		flaggs.PrintUsageAndExit()
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

	c := make(chan converter.Result, len(dirs))
	complete := make(chan bool, 1)

	fout, err := utils.OpenFileSafe(output)
	if err != nil && output != "" {
		if fout != nil {
			fout.Close()
		}

		log.Error("failed to open output file")
	}

	defer fout.Close()
	fout.Truncate(0)
	fout.Seek(0, 0)

	ferr, err := utils.OpenFileSafe(outputErr)
	if err != nil && outputErr != "" {
		if ferr != nil {
			ferr.Close()
		}

		log.Error("failed to open error file")
	}

	defer ferr.Close()
	ferr.Truncate(0)
	ferr.Seek(0, 0)

	go func() {
		for i := 0; i < len(dirs); i++ {
			r := <-c

			if r.Error == nil {
				if fout == nil {
					continue
				}

				fout.WriteString(r.OldHash)
				fout.WriteString("\t")
				fout.WriteString(r.NewHash)
				fout.WriteString("\t")
				fout.WriteString(r.Directory)
				fout.WriteString("\n")

				fout.Sync()
			} else {
				if ferr == nil {
					continue
				}

				ferr.WriteString(r.Error.Error())
				ferr.WriteString("\t")
				ferr.WriteString(r.Directory)
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

			result := converter.DirOldToNew(dir, dryRun, keepFiles)
			c <- result
		}(dir)
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Suffix = " Converting songes..."
	s.Start()

	<-complete
	s.Stop()
}
