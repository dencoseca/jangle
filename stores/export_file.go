package stores

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ExportFile represents a file used to manage export statements, with operations
// for appending and deleting entries.
type ExportFile struct {
	fileName string
}

// NewExportFile creates and returns a pointer to an ExportFile object
// initialized with the specified file path.
func NewExportFile(path string) *ExportFile {
	return &ExportFile{
		fileName: path,
	}
}

// Set appends an export statement for the given name to the file managed by the
// ExportFile instance. Returns an error if the file cannot be opened, written
// to, or closed properly.
func (cf ExportFile) Set(name string) error {
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)\n", name, name)

	file, err := os.OpenFile(cf.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error: failed to write to '%s: %w'", cf.fileName, err)
	}
	defer file.Close()

	if _, err := file.WriteString(exportLine); err != nil {
		return fmt.Errorf("error: failed to write to '%s': %w", cf.fileName, err)
	}

	return nil
}

// Delete removes an export statement for the provided name from the file managed
// by the ExportFile instance. Returns an error if the file cannot be read,
// written to, or closed properly.
func (cf ExportFile) Delete(name string) error {
	file, err := os.OpenFile(cf.fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error: Failed to open '%s': %w", cf.fileName, err)
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

	if scannerErr := scanner.Err(); scannerErr != nil {
		return fmt.Errorf("error reading from '%s': %w", cf.fileName, scannerErr)
	}

	if err := os.WriteFile(cf.fileName, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("error writing to '%s': %w", cf.fileName, err)
	}

	return nil
}
