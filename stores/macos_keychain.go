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

// NewMacOSKeychainStore initializes a new MacOSKeychainStore with an optional
// prefix for key namespaces.
func NewMacOSKeychainStore(prefix ...string) *MacOSKeychainStore {
	var namespace string

	if len(prefix) > 0 {
		namespace = prefix[0]
	} else {
		namespace = "jangle_"
	}

	return &MacOSKeychainStore{
		namespace: namespace,
	}
}

// Get retrieves a secret from the macOS Keychain by its name and returns the
// secret value or an error if not found.
func (m MacOSKeychainStore) Get(name string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", m.namespaced(name), "-w")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error: '%s' not found: %w", name, err)
	}

	return string(output), nil
}

// Set stores a key-value pair in the macOS Keychain. Returns an error if adding
// the key fails.
func (m MacOSKeychainStore) Set(name, value string) error {
	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("USER"), "-s", m.namespaced(name), "-w", value)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: failed to add '%s' to the Keychain: %w", name, err)
	}

	return nil
}

// Update attempts to update an existing secret in the macOS Keychain by deleting
// and re-adding the specified key-value pair. If the deletion or addition
// process fails, an error is returned.
func (m MacOSKeychainStore) Update(name, value string) error {
	if err := m.Delete(name); err != nil {
		return fmt.Errorf("error: failed to update '%s': %w", name, err)
	}

	if err := m.Set(name, value); err != nil {
		return fmt.Errorf("error: failed to update '%s': %w", name, err)
	}

	return nil
}

// Delete removes a secret identified by its name from the macOS Keychain.
// Returns an error if the operation fails.
func (m MacOSKeychainStore) Delete(name string) error {
	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("USER"), "-s", m.namespaced(name))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: failed to delete '%s' from the Keychain: %w", name, err)
	}

	return nil
}

// List retrieves a list of secret names stored in the macOS Keychain under the
// specified namespace. Returns an error if access fails.
func (m MacOSKeychainStore) List() ([]string, error) {
	cmd := exec.Command("security", "dump-keychain")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error accessing the Keychain: %w", err)
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

// namespaced appends the namespace to the provided name to ensure unique
// identification within the Keychain.
func (m MacOSKeychainStore) namespaced(name string) string {
	return m.namespace + name
}
