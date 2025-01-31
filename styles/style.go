package styles

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).Render

var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render

var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render

func Purple(message string, args ...any) {
	fmt.Println(headerStyle(fmt.Sprintf(message, args...)))
}

func Green(message string, args ...any) {
	fmt.Println(successStyle(fmt.Sprintf(message, args...)))
}

func Red(message string, args ...any) {
	fmt.Println(errorStyle(fmt.Sprintf(message, args...)))
}
