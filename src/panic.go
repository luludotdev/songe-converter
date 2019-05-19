package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func fatal(err error) {
	log.Print(err)
	fmt.Println("------------------------")
	fmt.Println("Press any key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}
