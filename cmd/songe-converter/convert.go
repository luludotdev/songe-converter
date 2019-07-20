package main

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/lolPants/songe-converter/converter"
	"github.com/lolPants/songe-converter/directory"
)

func convert(dir string, c chan result) {
	fail := func(err string) {
		log.Print(err)

		e := errors.New(err)
		res := result{
			dir:     dir,
			oldHash: "",
			newHash: "",
			err:     e,
		}

		c <- res
	}

	base := filepath.Base(dir)
	if base == "info.json" {
		dir = filepath.Dir(dir)
	}

	dirType, _ := directory.ReadType(dir)
	if dirType != directory.Old {
		fail("\"" + dir + "\" does not contain an old format beatmap")
		return
	}

	old, err := converter.ReadDirectoryOld(dir)
	if err != nil {
		fail("could not load beatmap at \"" + dir + "\"")
		return
	}

	new, err := converter.OldToNew(old)
	if err != nil {
		fail("failed to convert beatmap at \"" + dir + "\"")
		return
	}

	res := result{
		dir:     dir,
		oldHash: old.Hash,
		newHash: new.Hash,
		err:     nil,
	}

	c <- res
}
