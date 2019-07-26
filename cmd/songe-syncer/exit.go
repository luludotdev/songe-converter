package main

import (
	"bufio"
	"fmt"
	"os"
)

func exit(code int) {
	fmt.Println("\nPress [ENTER] to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	os.Exit(code)
}
