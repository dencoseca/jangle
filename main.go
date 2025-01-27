package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const keyPrefix = "jangle_"

func addPrefix(key string) string {
	return keyPrefix + key
}

func stripPrefix(key string) string {
	return strings.TrimPrefix(key, keyPrefix)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jangle <command> [args...]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "set":
		jangleSet(os.Args[2:])
	case "get":
		jangleGet(os.Args[2:])
	case "update":
		jangleUpdate(os.Args[2:])
	case "ls":
		jangleList()
	case "rm":
		jangleRemove(os.Args[2:])
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands: set, get, update, ls, rm")
		os.Exit(1)
	}
}

func jangleSet(args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: jangle set <key_name> <key_value>")
		os.Exit(1)
	}

	keyName := addPrefix(args[0])
	keyValue := args[1]

	if keyName == "" {
		fmt.Println("Error: <key_name> cannot be empty.")
		os.Exit(1)
	}

	if keyValue == "" {
		fmt.Println("Error: <key_value> cannot be empty.")
		os.Exit(1)
	}

	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w", keyValue)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Failed to add key '%s' to the Keychain.\n", stripPrefix(keyName))
		os.Exit(1)
	}

	fmt.Printf("Successfully added key '%s' to the Keychain.\n", stripPrefix(keyName))
}

func jangleGet(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: jangle get <key_name>")
		os.Exit(1)
	}

	keyName := addPrefix(args[0])

	if keyName == "" {
		fmt.Println("Error: <key_name> cannot be empty.")
		os.Exit(1)
	}

	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", keyName, "-w")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: Key '%s' not found in the Keychain.\n", stripPrefix(keyName))
		os.Exit(1)
	}

	fmt.Println(strings.TrimSpace(string(output)))
}

func jangleUpdate(args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: jangle update <key_name> <new_value>")
		os.Exit(1)
	}

	keyName := addPrefix(args[0])
	newValue := args[1]

	if keyName == "" {
		fmt.Println("Error: <key_name> cannot be empty.")
		os.Exit(1)
	}

	if newValue == "" {
		fmt.Println("Error: <new_value> cannot be empty.")
		os.Exit(1)
	}

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

func jangleList() {
	// Execute the security command to dump all keychain items
	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.CombinedOutput() // Capture both stdout and stderr
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
		fmt.Printf("No jangle_ keys found in the Keychain for user '%s'.\n", os.Getenv("USER"))
		os.Exit(1)
	}

	// Print keys prefixed with `jangle_`
	fmt.Printf("jangle_ keys in the Keychain for user '%s':\n", os.Getenv("USER"))
	for _, key := range keys {
		fmt.Printf("  - %s\n", key)
	}
}
func jangleRemove(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: jangle rm <key_name>")
		os.Exit(1)
	}

	keyName := addPrefix(args[0])

	if keyName == "" {
		fmt.Println("Error: <key_name> cannot be empty.")
		os.Exit(1)
	}

	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", keyName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Failed to remove key '%s' from the Keychain.\n", stripPrefix(keyName))
		os.Exit(1)
	}

	fmt.Printf("Successfully removed key '%s' from the Keychain.\n", stripPrefix(keyName))
}
