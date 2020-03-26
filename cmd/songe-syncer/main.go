package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"

	"jackbaron.com/songe-converter/v2/directory"
	"jackbaron.com/songe-converter/v2/flags"
	"jackbaron.com/songe-converter/v2/utils"
)

var (
	dir       string
	outputDir string

	printHelp bool
	args      []string
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Println("Drag a beatmap folder onto this binary to watch it")
		fmt.Println("Beatmaps can be converted to and from the new format")

		exit(1)
	}

	flags.SetDetails("Songe Syncer "+gitTag, "https://github.com/lolPants/songe-converter")
	flags.SetExample("./songe-syncer")
	flags.RegisterBoolFlag(&printVer, "v", "version", "print version information")
	flags.RegisterStringFlag(&outputDir, "o", "output", "output directory")
	flags.Parse(&args)

	if printVer == true {
		printVersionInfo()
		os.Exit(0)
	}

	if printHelp == true {
		flags.PrintUsageAndExit()
		return
	}

	dir = args[0]
	dirType, _ := directory.ReadType(dir)
	if dirType != directory.Old {
		fmt.Println("This folder does not contain an old beatmap!")
		exit(1)
	}

	var cacheFile string
	var cachedPath string
	dataDir, err := utils.DataDir("songe-syncer")
	if err == nil && dataDir != "" {
		if exists, _ := utils.DirectoryExists(dataDir); exists == false {
			os.MkdirAll(dataDir, 0755)
		}

		cacheFile = filepath.Join(dataDir, "path.cache")
		bytes, err := ioutil.ReadFile(cacheFile)

		if err == nil {
			cachedPath = string(bytes)
		}
	}

	if outputDir == "" {
		fmt.Println("Input path to CustomLevels / CustomWIPLevels folder:")
		if cachedPath != "" {
			fmt.Printf("Detected previous path: \"%s\"\n", cachedPath)
			fmt.Print("Press ENTER to use this path\n\n")
		}

		fmt.Print("> ")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			outputDir = scanner.Text()
			if outputDir == "" && cachedPath != "" {
				outputDir = cachedPath
			}
		}
	}

	outputDir = utils.StripQuotes(outputDir)
	outputExists, err := utils.DirectoryExists(outputDir)
	if err != nil || outputExists == false {
		fmt.Println("Output directory does not exist!")
		exit(1)
	}

	if outputDir == dir {
		fmt.Println("Output directory cannot be the same as the input directory!")
		exit(0)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Could not watch this folder for changes!")
		exit(1)
	}

	if cacheFile != "" {
		err = ioutil.WriteFile(cacheFile, []byte(outputDir), 0644)
		if err != nil {
			log.Print(err)
		}
	}

	defer watcher.Close()
	done := make(chan bool)
	change := make(chan bool)

	debounced := debounce.New(100 * time.Millisecond)
	go func() {
		for {
			<-change
			debounced(processDir)
		}
	}()

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}

				change <- true

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("watch error:", err)
			}
		}
	}()

	err = watcher.Add(args[0])
	if err != nil {
		fmt.Println("Could not watch this folder for changes!")
		exit(1)
	}

	debounced(processDir)
	<-done
}
