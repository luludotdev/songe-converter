package main

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

func printError(msg string) {
	fmt.Println(chalk.Red.Color("[ERROR]") + " " + msg)
	os.Exit(1)
}

func printWarning(msg string) {
	fmt.Println(chalk.Yellow.Color("[WARNING]") + " " + msg)
}
