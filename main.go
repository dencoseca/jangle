package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styling using lipgloss
var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).Render
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render
	errorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(mainUsage)
		os.Exit(1)
	}

	mainCommand := os.Args[1]
	args := os.Args[2:]

	switch mainCommand {
	case "--help", "help":
		fmt.Println(mainUsage)
	case "set":
		set(args)
	case "get":
		get(args)
	case "update":
		update(args)
	case "ls":
		list()
	case "delete":
		remove(args)
	default:
		fmt.Println(mainUsage)
		os.Exit(1)
	}
}

func set(args []string) {
	if len(args) != 2 {
		fmt.Println(setUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])
	keyValue := args[1]

	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w", keyValue)
	if err := cmd.Run(); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to add key '%s' to the Keychain.", stripPrefix(keyName))))
		os.Exit(1)
	}

	fmt.Println(successStyle(fmt.Sprintf("Successfully added key '%s' to the Keychain.", stripPrefix(keyName))))
}

func get(args []string) {
	if len(args) != 1 {
		fmt.Println(getUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])

	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Key '%s' not found in the Keychain.", stripPrefix(keyName))))
		os.Exit(1)
	}

	fmt.Println(strings.TrimSpace(string(output)))
}

func update(args []string) {
	if len(args) != 2 {
		fmt.Println(updateUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])
	newValue := args[1]

	// Delete the old key (if it exists)
	cmdDelete := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", keyName)
	cmdDelete.Run() // Ignore errors since we are updating it anyway

	// Add a new key
	cmdAdd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w", newValue)
	if err := cmdAdd.Run(); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to update key '%s' in the Keychain.", stripPrefix(keyName))))
		os.Exit(1)
	}

	fmt.Println(successStyle(fmt.Sprintf("Successfully updated key '%s' in the Keychain.", stripPrefix(keyName))))
}

func list() {
	if len(os.Args[2:]) > 0 {
		fmt.Println(listUsage)
		os.Exit(1)
	}

	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error accessing the Keychain: %v", err)))
		os.Exit(1)
	}

	regex := regexp.MustCompile(`"svce"<blob>="jangle_(.*?)"`)
	matches := regex.FindAllStringSubmatch(string(output), -1)

	var keys []string
	for _, match := range matches {
		if len(match) > 1 {
			keys = append(keys, match[1])
		}
	}

	if len(keys) == 0 {
		fmt.Println(errorStyle(fmt.Sprintf("No keys found in the Keychain for user '%s'.", os.Getenv("USER"))))
		os.Exit(1)
	}

	fmt.Println(headerStyle(fmt.Sprintf("Keys in the Keychain for user '%s':", os.Getenv("USER"))))
	for _, key := range keys {
		fmt.Println("  - " + key)
	}
}

func remove(args []string) {
	if len(args) != 1 {
		fmt.Println(deleteUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])

	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", keyName)
	if err := cmd.Run(); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to remove key '%s' from the Keychain.", stripPrefix(keyName))))
		os.Exit(1)
	}

	fmt.Println(successStyle(fmt.Sprintf("Successfully removed key '%s' from the Keychain.", stripPrefix(keyName))))
}
