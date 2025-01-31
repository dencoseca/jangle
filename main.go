package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).Render
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render
	errorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render
	janglercPath = os.Getenv("HOME") + "/.janglerc"
)

// main serves as the entry point for the program, handling command-line
// arguments to execute various keychain operations.
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

// set adds a key-value pair to the macOS Keychain using the provided arguments
// or exits with an error if arguments are invalid.
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

	// export the secret in .janglerc
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)\n", stripPrefix(keyName), stripPrefix(keyName))

	file, err := os.OpenFile(janglercPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to write to '%s'.", janglercPath)))
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.WriteString(exportLine); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to write to '%s'.", janglercPath)))
		os.Exit(1)
	}

	// Set the key in the shell environment
	envVar := stripPrefix(keyName)
	if err := os.Setenv(envVar, keyValue); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to set '%s' in the environment variables: %v", envVar, err)))
	} else {
		fmt.Println(successStyle(fmt.Sprintf("Successfully set '%s' in the environment variables.", envVar)))
	}

	fmt.Println(successStyle(fmt.Sprintf("Successfully added key '%s'.", stripPrefix(keyName))))
	fmt.Println(successStyle("Source your terminal configuration or restart your shell to use the environment variable."))
}

// get retrieves a secret value by name from the macOS Keychain or exits with an
// error if the key is not found.
func get(args []string) {
	if len(args) != 1 {
		fmt.Println(getUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])

	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Key '%s' not found.", stripPrefix(keyName))))
		os.Exit(1)
	}

	fmt.Println(strings.TrimSpace(string(output)))
}

// update updates an existing key-value pair in the macOS Keychain or adds a new
// one if the key does not exist.
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

	fmt.Println(successStyle(fmt.Sprintf("Successfully updated key '%s'.", stripPrefix(keyName))))
}

// the list retrieves and displays all keys related to "jangle" stored in the macOS Keychain for the current user.
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
		fmt.Println(errorStyle(fmt.Sprintf("No keys found for user '%s'.", os.Getenv("USER"))))
		os.Exit(1)
	}

	fmt.Println(headerStyle(fmt.Sprintf("Keys for user '%s':\n", os.Getenv("USER"))))
	for _, key := range keys {
		fmt.Println("- " + key)
	}
}

// remove deletes a secret from the macOS Keychain based on the provided name or
// exits with an error if arguments are invalid.
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

	// Remove the corresponding export line from .janglerc
	file, err := os.OpenFile(janglercPath, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(errorStyle(fmt.Sprintf("Warning: File '%s' does not exist, no export to remove.", janglercPath)))
		} else {
			fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to open '%s'.", janglercPath)))
			os.Exit(1)
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)", stripPrefix(keyName), stripPrefix(keyName))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != exportLine {
			lines = append(lines, line)
		}
	}

	if scannerErr := scanner.Err(); scannerErr != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error reading from '%s': %v", janglercPath, scannerErr)))
		os.Exit(1)
	}

	// Rewrite the file with the updated content
	if err := os.WriteFile(janglercPath, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error writing to '%s': %v", janglercPath, err)))
		os.Exit(1)
	}

	// remove the key from the shell environment
	envVar := stripPrefix(keyName)
	if err := os.Unsetenv(envVar); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Warning: Failed to remove '%s' from the environment variables: %v", envVar, err)))
	}

	fmt.Println(successStyle(fmt.Sprintf("Successfully removed key '%s'.", stripPrefix(keyName))))
	fmt.Println(successStyle(fmt.Sprintf("To remove the environment variable restart your terminal or run: unset %s", envVar)))

}
