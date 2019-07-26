package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lolPants/songe-converter/converter"
	"github.com/lolPants/songe-converter/directory"
	"github.com/otiai10/copy"
)

func main() {
	exit := func(code int) {
		fmt.Println("\nPress [ENTER] to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		os.Exit(code)
	}

	if len(os.Args[1:]) == 0 {
		fmt.Println("Drag a beatmap folder onto this binary to convert it")
		fmt.Println("Beatmaps can be converted to and from the new format")

		exit(1)
	}

	dir := os.Args[1]
	dirType, _ := directory.ReadType(dir)

	if dirType == directory.None {
		fmt.Println("This directory does not contain a beatmap!")

		exit(1)
	} else if dirType == directory.Both {
		fmt.Println("This directory contains both beatmap formats!")
		fmt.Println("Please delete one so I know which to convert")

		exit(1)
	}

	backup := fmt.Sprintf("%v [BACKUP]", dir)
	err := copy.Copy(dir, backup)
	if err != nil {
		fmt.Println("Could not make a backup of the beatmap!")
		exit(1)
	}

	if dirType == directory.Old {
		converter.DirOldToNew(dir, false, false)
	} else if dirType == directory.New {
		converter.DirNewToOld(dir, false, false)
	}

	fmt.Println("Conversion complete!")
	exit(0)
}
