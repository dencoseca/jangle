package stores

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// ExportFile represents a file used to manage export statements, with operations
// for appending and deleting entries.
type ExportFile struct {
	filePath string
}

// NewExportFile creates a new ExportFile instance with the optional given path
// or defaults to "$HOME/.jangle_exports".
func NewExportFile(path ...string) *ExportFile {
	var filePath string

	if len(path) == 0 {
		filePath = os.Getenv("HOME") + "/.jangle_exports"
	} else {
		filePath = path[0]
	}

	return &ExportFile{
		filePath: filePath,
	}
}

// Set appends an export statement for the given name to the file managed by the
// ExportFile instance. Returns an error if the file cannot be opened, written
// to, or closed properly.
func (ef ExportFile) Set(name string) error {
	if name == "" {
		return errors.New("error: no name provided")
	}

	file, err := os.OpenFile(ef.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error: failed to write to '%s: %w'", ef.filePath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(ef.shellExportStatement(name)); err != nil {
		return fmt.Errorf("error: failed to write to '%s': %w", ef.filePath, err)
	}

	return nil
}

// Delete removes an export statement for the provided name from the file managed
// by the ExportFile instance. Returns an error if the file cannot be read,
// written to, or closed properly.
func (ef ExportFile) Delete(name string) error {
	file, err := os.OpenFile(ef.filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error: failed to open '%s': %w", ef.filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != ef.shellExportStatement(name) {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from '%s': %w", ef.filePath, err)
	}

	if err := os.WriteFile(ef.filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("error writing to '%s': %w", ef.filePath, err)
	}

	return nil
}

// shellExportStatement generates a shell export statement for the given variable
// name and its value using "jangle get".
func (ef ExportFile) shellExportStatement(name string) string {
	return fmt.Sprintf("export %s=$(jangle get %s)", name, name)
}
