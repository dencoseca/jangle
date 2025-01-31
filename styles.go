package main

import "github.com/charmbracelet/lipgloss"

var headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).Render

var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render

var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render
