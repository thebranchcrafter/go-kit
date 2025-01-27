package utils

import (
	"fmt"
	"github.com/fatih/color"
)

// PrintGreen prints a message in green.
func PrintGreen(message string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green(message))
}

// PrintRed prints a message in red.
func PrintRed(message string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Println(red(message))
}

// PrintYellow prints a message in yellow.
func PrintYellow(message string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(yellow(message))
}

// PrintCyan prints a message in cyan.
func PrintCyan(message string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Println(cyan(message))
}
