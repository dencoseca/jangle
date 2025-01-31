package stores

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var envVariableFile = os.Getenv("HOME") + "/.janglerc"

type MacOSKeychainStore struct {
	prefix string
}

func NewMacOSKeychainStore() *MacOSKeychainStore {
	return &MacOSKeychainStore{
		prefix: "jangle_",
	}
}

func (mk MacOSKeychainStore) Get(name string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", mk.prefixName(name), "-w")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error: '%s' not found", name)
	}

	return string(output), nil

}

func (mk MacOSKeychainStore) Set(name, value string) error {
	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", mk.prefixName(name), "-w", value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: failed to add '%s' to the Keychain", name)
	}

	// export the secret in .janglerc
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)\n", name, name)

	file, err := os.OpenFile(envVariableFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error: failed to write to '%s'", envVariableFile)
	}
	defer file.Close()

	if _, err := file.WriteString(exportLine); err != nil {
		return fmt.Errorf("error: failed to write to '%s'", envVariableFile)
	}

	return nil
}

func (mk MacOSKeychainStore) Update(name, value string) error {
	// Delete the old key (if it exists)
	cmdDelete := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", mk.prefixName(name))
	cmdDelete.Run() // Ignore errors since we are updating it anyway

	// Add a new key
	cmdAdd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", mk.prefixName(name), "-w", value)
	if err := cmdAdd.Run(); err != nil {
		return fmt.Errorf("error: Failed to Update '%s' in the Keychain", name)
	}

	return nil
}

func (mk MacOSKeychainStore) Delete(name string) error {
	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", mk.prefixName(name))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: Failed to Delete '%s' from the Keychain", name)
	}

	// Delete the corresponding export line from .janglerc
	file, err := os.OpenFile(envVariableFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error: Failed to open '%s'", envVariableFile)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	exportLine := fmt.Sprintf("export %s=$(jangle Get %s)", name, name)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != exportLine {
			lines = append(lines, line)
		}
	}

	if scannerErr := scanner.Err(); scannerErr != nil {
		return fmt.Errorf("error reading from '%s': %v", envVariableFile, scannerErr)
	}

	// Rewrite the file with the updated content
	if err := os.WriteFile(envVariableFile, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("error writing to '%s': %v", envVariableFile, err)
	}

	return nil
}

func (mk MacOSKeychainStore) List() ([]string, error) {
	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error accessing the Keychain: %v", err)
	}

	regex := regexp.MustCompile(`"svce"<blob>="jangle_(.*?)"`)
	matches := regex.FindAllStringSubmatch(string(output), -1)

	var secretNames []string
	for _, match := range matches {
		if len(match) > 1 {
			secretNames = append(secretNames, match[1])
		}
	}

	return secretNames, nil
}

func (mk MacOSKeychainStore) prefixName(name string) string {
	return mk.prefix + name
}
