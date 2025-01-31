package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Secret struct {
	Name  string
	Value string
}

// prefixedName returns the secret's name prefixedName with the NAMESPACE_PREFIX constant
// for ensuring uniqueness.
func (s Secret) prefixedName() string {
	return NAMESPACE_PREFIX + s.Name
}

// set adds a named secret with a specified value to the macOS Keychain and
// updates the .janglerc file with the export directive.
func (s Secret) set() {
	if s.Name == "" || s.Value == "" {
		fmt.Println(setUsage)
		os.Exit(1)
	}

	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", s.prefixedName(), "-w", s.Value)
	if err := cmd.Run(); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to add '%s' to the Keychain.", s.Name)))
		os.Exit(1)
	}

	// export the secret in .janglerc
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)\n", s.Name, s.Name)

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

	fmt.Println(successStyle(fmt.Sprintf("Successfully added '%s'.\n", s.Name)))
	fmt.Println("Source your terminal configuration or restart your shell to use the environment variable.")
}

// get retrieves a secret value by name from the macOS Keychain or exits with an
// error if the key is not found.
func (s Secret) get() {
	if s.Name == "" {
		fmt.Println(getUsage)
		os.Exit(1)
	}

	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", s.prefixedName(), "-w")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: '%s' not found.", s.Name)))
		os.Exit(1)
	}

	fmt.Println(strings.TrimSpace(string(output)))
}

// update updates an existing key-value pair in the macOS Keychain or adds a new
// one if the key does not exist.
func (s Secret) update() {
	if s.Name == "" || s.Value == "" {
		fmt.Println(updateUsage)
		os.Exit(1)
	}

	// Delete the old key (if it exists)
	cmdDelete := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", s.prefixedName())
	cmdDelete.Run() // Ignore errors since we are updating it anyway

	// Add a new key
	cmdAdd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", s.prefixedName(), "-w", s.Value)
	if err := cmdAdd.Run(); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to update '%s' in the Keychain.", s.Name)))
		os.Exit(1)
	}

	fmt.Println(successStyle(fmt.Sprintf("Successfully updated '%s'.\n", s.Name)))
	fmt.Println("Source your terminal configuration or restart your shell to use the updated environment variable.")
}

// remove deletes a secret by name from the macOS Keychain and removes the
// corresponding entry from .janglerc.
func (s Secret) remove() {
	if s.Name == "" {
		fmt.Println(deleteUsage)
		os.Exit(1)
	}

	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", s.prefixedName())
	if err := cmd.Run(); err != nil {
		fmt.Println(errorStyle(fmt.Sprintf("Error: Failed to remove '%s' from the Keychain.", s.Name)))
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
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)", s.Name, s.Name)
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

	fmt.Println(successStyle(fmt.Sprintf("Successfully removed '%s'.\n", s.Name)))
	fmt.Println(fmt.Sprintf("To remove the environment variable restart your terminal or run: unset %s", s.Name))
}
