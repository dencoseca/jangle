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

// NewExportFile creates a new ExportFile instance with the specified file path.
// Returns an error if the path is empty.
func NewExportFile(path string) (*ExportFile, error) {
	if path == "" {
		return nil, errors.New("error: no path provided")
	}

	return &ExportFile{
		filePath: path,
	}, nil
}

// Set appends an export statement for the given name to the file managed by the
// ExportFile instance. Returns an error if the file cannot be opened, written
// to, or closed properly.
func (cf ExportFile) Set(name string) error {
	if name == "" {
		return errors.New("error: no name provided")
	}

	exportLine := fmt.Sprintf("export %s=$(jangle get %s)\n", name, name)

	file, err := os.OpenFile(cf.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error: failed to write to '%s: %w'", cf.filePath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(exportLine); err != nil {
		return fmt.Errorf("error: failed to write to '%s': %w", cf.filePath, err)
	}

	return nil
}

// Delete removes an export statement for the provided name from the file managed
// by the ExportFile instance. Returns an error if the file cannot be read,
// written to, or closed properly.
func (cf ExportFile) Delete(name string) error {
	file, err := os.OpenFile(cf.filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error: Failed to open '%s': %w", cf.filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)", name, name)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != exportLine {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from '%s': %w", cf.filePath, err)
	}

	if err := os.WriteFile(cf.filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("error writing to '%s': %w", cf.filePath, err)
	}

	return nil
}
