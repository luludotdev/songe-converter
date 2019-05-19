package main

import (
	"log"
	"os"
)

func fatal(err error) {
	log.Print(err)
	os.Exit(1)
}

func fatalStr(msg string) {
	log.Print(msg)
	os.Exit(1)
}
