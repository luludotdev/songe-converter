package utils

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

// Error prints error string and exits
func Error(msg string) {
	fmt.Println(chalk.Red.Color("[ERROR]") + " " + msg)
	os.Exit(1)
}

// Warning prints warning string
func Warning(msg string) {
	fmt.Println(chalk.Yellow.Color("[WARNING]") + " " + msg)
}
