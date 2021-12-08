package cmd

import (
	"fmt"

	"github.com/gookit/color"
)

func printErrorMessage(message string) {
	error := color.FgRed.Render
	printMessage(error("ERROR"), message)
}

func printSuccessMessage(message string) {
	success := color.FgCyan.Render
	printMessage(success("SUCCESS"), message)
}

func printMessage(status, message string) {
	fmt.Printf("%-7s %s\n", status, message)
}
