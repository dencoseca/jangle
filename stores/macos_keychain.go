package stores

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type MacOSKeychainStore struct {
	namespace string
}

func NewMacOSKeychainStore() *MacOSKeychainStore {
	return &MacOSKeychainStore{
		namespace: "jangle_",
	}
}

func (m MacOSKeychainStore) Get(name string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", m.prefixName(name), "-w")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error: '%s' not found", name)
	}

	return string(output), nil
}

func (m MacOSKeychainStore) Set(name, value string) error {
	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", m.prefixName(name), "-w", value)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: failed to add '%s' to the Keychain", name)
	}

	return nil
}

func (m MacOSKeychainStore) Update(name, value string) error {
	if err := m.Delete(name); err != nil {
		return err
	}

	if err := m.Set(name, value); err != nil {
		return err
	}

	return nil
}

func (m MacOSKeychainStore) Delete(name string) error {
	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", m.prefixName(name))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: Failed to Delete '%s' from the Keychain", name)
	}

	return nil
}

func (m MacOSKeychainStore) List() ([]string, error) {
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

func (m MacOSKeychainStore) prefixName(name string) string {
	return m.namespace + name
}
