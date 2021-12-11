package cmd

import (
	"fmt"

	"github.com/gookit/color"
)

var errorColor = color.FgRed
var warnColor = color.FgYellow
var infoColor = color.FgCyan
var passColor = color.FgCyan
var successColor = color.FgCyan

// pass := color.FgCyan.Render
var success = color.FgCyan.Render

func printErrorMessage(message string) {
	printMessage(errorColor.Render("ERROR"), message)
}

func printSuccessMessage(message string) {
	printMessage(successColor.Render("SUCCESS"), message)
}

func printMessage(status, message string) {
	fmt.Printf("%-16s %s\n", status, message)
}
