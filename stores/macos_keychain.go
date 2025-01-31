package stores

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// MacOSKeychainStore provides methods to interact with the macOS Keychain for
// storing, retrieving, updating, and deleting secrets. The namespace field
// specifies a prefix for stored secret keys to avoid conflicts with other
// entries.
type MacOSKeychainStore struct {
	namespace string
}

// NewMacOSKeychainStore initializes and returns a new MacOSKeychainStore
// instance with a predefined namespace.
func NewMacOSKeychainStore() *MacOSKeychainStore {
	return &MacOSKeychainStore{
		namespace: "jangle_",
	}
}

// Get retrieves a secret from the macOS Keychain by its name and returns the
// secret value or an error if not found.
func (m MacOSKeychainStore) Get(name string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", m.prefixName(name), "-w")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error: '%s' not found", name)
	}

	return string(output), nil
}

// Set stores a key-value pair in the macOS Keychain. Returns an error if adding
// the key fails.
func (m MacOSKeychainStore) Set(name, value string) error {
	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", m.prefixName(name), "-w", value)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: failed to add '%s' to the Keychain", name)
	}

	return nil
}

// Update attempts to update an existing secret in the macOS Keychain by deleting
// and re-adding the specified key-value pair. If the deletion or addition
// process fails, an error is returned.
func (m MacOSKeychainStore) Update(name, value string) error {
	if err := m.Delete(name); err != nil {
		return err
	}

	if err := m.Set(name, value); err != nil {
		return err
	}

	return nil
}

// Delete removes a secret identified by its name from the macOS Keychain.
// Returns an error if the operation fails.
func (m MacOSKeychainStore) Delete(name string) error {
	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", m.prefixName(name))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: Failed to Delete '%s' from the Keychain", name)
	}

	return nil
}

// List retrieves a list of secret names stored in the macOS Keychain under the
// specified namespace. Returns an error if access fails.
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

// prefixName appends the namespace to the provided name to ensure unique
// identification within the Keychain.
func (m MacOSKeychainStore) prefixName(name string) string {
	return m.namespace + name
}
