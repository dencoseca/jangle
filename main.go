package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// main executes the command-line interface for managing Mac OS Keychain secrets,
// handling various commands and their arguments.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jangle <command> [args...]")
		os.Exit(1)
	}

	mainCommand := os.Args[1]

	switch mainCommand {
	case "--help", "help":
		fmt.Println(mainUsage)
	case "set":
		set(os.Args[2:])
	case "get":
		get(os.Args[2:])
	case "update":
		update(os.Args[2:])
	case "ls":
		list()
	case "delete":
		delete(os.Args[2:])
	default:
		fmt.Println(mainUsage)
		os.Exit(1)
	}
}

// set stores a secret with a specified name and value in the macOS Keychain,
// ensuring a required input format.
func set(args []string) {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "help") {
		fmt.Println(setUsage)
		return
	}

	if len(args) != 2 {
		fmt.Println(setUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])
	keyValue := args[1]

	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w", keyValue)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Failed to add key '%s' to the Keychain.\n", stripPrefix(keyName))
		os.Exit(1)
	}

	fmt.Printf("Successfully added key '%s' to the Keychain.\n", stripPrefix(keyName))
}

// get retrieves a secret value by its name from the macOS Keychain and prints it
// to the console. It validates the input arguments and handles errors if the
// secret is not found.
func get(args []string) {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "help") {
		fmt.Println(getUsage)
		return
	}

	if len(args) != 1 {
		fmt.Println(getUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])

	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: Key '%s' not found in the Keychain.\n", stripPrefix(keyName))
		os.Exit(1)
	}

	fmt.Println(strings.TrimSpace(string(output)))
}

// update modifies the value of an existing secret in the macOS Keychain or
// creates it if it does not already exist.
func update(args []string) {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "help") {
		fmt.Println(updateUsage)
		return
	}

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
		fmt.Printf("Error: Failed to update key '%s' in the Keychain.\n", stripPrefix(keyName))
		os.Exit(1)
	}

	fmt.Printf("Successfully updated key '%s' in the Keychain.\n", stripPrefix(keyName))
}

// list retrieves and lists all keys stored in the macOS Keychain with the prefix
// `jangle_` for the current user.
func list() {
	if len(os.Args[2:]) > 0 {
		fmt.Println(listUsage)
		os.Exit(1)
	}

	// Execute the security command to dump all keychain items
	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error accessing the Keychain: %v\n", err)
		os.Exit(1)
	}

	// Regex pattern to match and extract keys with the prefix `jangle_`
	regex := regexp.MustCompile(`"svce"<blob>="jangle_(.*?)"`)
	matches := regex.FindAllStringSubmatch(string(output), -1)

	// Extract matched keys
	var keys []string
	for _, match := range matches {
		if len(match) > 1 { // Ensure there's a captured group
			keys = append(keys, match[1])
		}
	}

	// Handle case where no keys were found with the `jangle_` prefix
	if len(keys) == 0 {
		fmt.Printf("No keys found in the Keychain for user '%s'.\n", os.Getenv("USER"))
		os.Exit(1)
	}

	// Print keys prefixed with `jangle_`
	fmt.Printf("Keys in the Keychain for user '%s':\n", os.Getenv("USER"))
	for _, key := range keys {
		fmt.Printf("  - %s\n", key)
	}
}

// delete removes a secret by name from the macOS Keychain using the specified
// arguments. Validates input and handles errors.
func delete(args []string) {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "help") {
		fmt.Println(deleteUsage)
		return
	}

	if len(args) != 1 {
		fmt.Println(deleteUsage)
		os.Exit(1)
	}

	keyName := addPrefix(args[0])

	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", keyName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Failed to remove key '%s' from the Keychain.\n", stripPrefix(keyName))
		os.Exit(1)
	}

	fmt.Printf("Successfully removed key '%s' from the Keychain.\n", stripPrefix(keyName))
}
