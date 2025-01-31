package stores

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

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
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error: failed to add '%s' to the Keychain", name)
	}

	return nil
}

func (mk MacOSKeychainStore) Update(name, value string) error {
	err := mk.Delete(name)
	if err != nil {
		return err
	}

	err = mk.Set(name, value)
	if err != nil {
		return err
	}

	return nil
}

func (mk MacOSKeychainStore) Delete(name string) error {
	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", mk.prefixName(name))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error: Failed to Delete '%s' from the Keychain", name)
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
