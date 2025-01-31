package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// NAMESPACE_PREFIX defines a default prefix used for namespacing keys to ensure
// uniqueness and avoid collisions.
const NAMESPACE_PREFIX = "jangle_"

// janglercPath stores the path to the user's .janglerc configuration file in
// their home directory.
var janglercPath = os.Getenv("HOME") + "/.janglerc"

// main serves as the entry point for the program, handling command-line
// arguments to execute various keychain operations.
func main() {
	if len(os.Args) < 2 {
		fmt.Println(mainUsage)
		os.Exit(1)
	}

	ensureJangleFileExists()

	command := os.Args[1]

	secret := Secret{
		Name:  getArg(2, ""),
		Value: getArg(3, ""),
	}

	switch command {
	case "--help", "help":
		fmt.Println(mainUsage)
	case "set":
		secret.set()
	case "get":
		secret.get()
	case "update":
		secret.update()
	case "ls":
		list()
	case "delete":
		secret.remove()
	default:
		fmt.Println(mainUsage)
		os.Exit(1)
	}
}

// the list retrieves and displays all keys related to "jangle" stored in the
// macOS Keychain for the current user.
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
		fmt.Println(errorStyle(fmt.Sprintf("No secrets found for user '%s'.", os.Getenv("USER"))))
		os.Exit(1)
	}

	fmt.Println(headerStyle(fmt.Sprintf("Secrets for user '%s':\n", os.Getenv("USER"))))
	for _, key := range keys {
		fmt.Println("- " + key)
	}
}

// ensureJangleFileExists checks if the .janglerc file exists and creates it if
// not, exiting on any creation error.
func ensureJangleFileExists() {
	if _, err := os.Stat(janglercPath); os.IsNotExist(err) {
		file, err := os.Create(janglercPath)
		if err != nil {
			fmt.Println(errorStyle(fmt.Sprintf("Error creating the .janglerc file: %v", err)))
			os.Exit(1)
		}
		file.Close()
	}
}

// getArg returns the command-line argument at the given index if available,
// otherwise returns the defaultValue.
func getArg(index int, defaultValue string) string {
	if len(os.Args) > index {
		return os.Args[index]
	}

	return defaultValue
}
