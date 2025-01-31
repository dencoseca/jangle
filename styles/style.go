package styles

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render

var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render

// Green outputs a formatted success message in bold green text to the console
// using the provided message and arguments.
func Green(message string, args ...any) {
	fmt.Println(successStyle(fmt.Sprintf(message, args...)))
}

// Red formats a message with error styling and prints it to the console using
// fmt.Printf-like arguments.
func Red(message string, args ...any) {
	fmt.Println(errorStyle(fmt.Sprintf(message, args...)))
}
