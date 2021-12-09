package cmd

import (
	"fmt"

	"github.com/gookit/color"
)

var warn = color.FgYellow.Render
var info = color.FgCyan.Render

// pass := color.FgCyan.Render
var success = color.FgCyan.Render

func printErrorMessage(message string) {
	error := color.FgRed.Render
	printMessage(error("ERROR"), message)
}

func printSuccessMessage(message string) {
	printMessage(success("SUCCESS"), message)
}

func printMessage(status, message string) {
	fmt.Printf("%-16s %s\n", status, message)
}
