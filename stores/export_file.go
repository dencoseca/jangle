package stores

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ExportFile struct {
	fileName string
}

func NewExportFile(path string) *ExportFile {
	return &ExportFile{
		fileName: path,
	}
}

func (cf ExportFile) Set(name string) error {
	exportLine := fmt.Sprintf("export %s=$(jangle get %s)\n", name, name)

	file, err := os.OpenFile(cf.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error: failed to write to '%s'", cf.fileName)
	}
	defer file.Close()

	if _, err := file.WriteString(exportLine); err != nil {
		return fmt.Errorf("error: failed to write to '%s'", cf.fileName)
	}

	return nil
}

func (cf ExportFile) Delete(name string) error {
	file, err := os.OpenFile(cf.fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error: Failed to open '%s'", cf.fileName)
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
		return fmt.Errorf("error reading from '%s': %v", cf.fileName, scannerErr)
	}

	if err := os.WriteFile(cf.fileName, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("error writing to '%s': %v", cf.fileName, err)
	}

	return nil
}
