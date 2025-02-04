package stores

import (
	"os"
	"testing"
)

func TestNewExportFile(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedFilePath string
	}{
		{
			name:             "uses default filePath when no path is provided",
			input:            "",
			expectedFilePath: "/mockHome/.jangle_exports",
		},
		{
			name:             "when a filePath is provided it uses that path",
			input:            "/Users/test.user/.test_file",
			expectedFilePath: "/Users/test.user/.test_file",
		},
	}

	err := os.Setenv("HOME", "/mockHome")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var exportFile *ExportFile
			if tt.input == "" {
				exportFile = NewExportFile()
			} else {
				exportFile = NewExportFile(tt.input)
			}

			if exportFile.filePath != tt.expectedFilePath {
				t.Errorf("wrong export filePath: got %s, want %s", exportFile.filePath, tt.expectedFilePath)
			}
		})
	}
}

func TestExportFile_Set(t *testing.T) {
	tests := []struct {
		name                string
		initialFileContent  string
		expectedFileContent string
	}{
		{
			name:                "when the file is empty it appends the export line",
			initialFileContent:  "",
			expectedFileContent: "export TEST_TOKEN=$(jangle get TEST_TOKEN)\n",
		},
		{
			name:                "when the file has a single entry it appends the export line and retains the existing content",
			initialFileContent:  "export TEST_TOKEN=$(jangle get TEST_TOKEN)\n",
			expectedFileContent: "export TEST_TOKEN=$(jangle get TEST_TOKEN)\nexport TEST_TOKEN=$(jangle get TEST_TOKEN)\n",
		},
		{
			name:                "when the file has multiple entries it appends the export line and retains the existing content",
			initialFileContent:  "export TEST_TOKEN=$(jangle get TEST_TOKEN)\nexport TEST_TOKEN=$(jangle get TEST_TOKEN)\n",
			expectedFileContent: "export TEST_TOKEN=$(jangle get TEST_TOKEN)\nexport TEST_TOKEN=$(jangle get TEST_TOKEN)\nexport TEST_TOKEN=$(jangle get TEST_TOKEN)\n",
		},
	}

	err := os.Setenv("HOME", "./testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := createTestFile(t, tt.initialFileContent)
			defer cleanup()

			exportFile := NewExportFile()

			err = exportFile.Set("TEST_TOKEN")
			if err != nil {
				t.Fatal(err)
			}

			assertFileContent(t, tt.expectedFileContent)
		})
	}
}

const testFilePath = "./testdata/.jangle_exports"

func createTestFile(t *testing.T, initialFileContent string) func() {
	err := os.MkdirAll("./testdata", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if initialFileContent != "" {
		_, err := file.WriteString(initialFileContent)
		if err != nil {
			t.Fatal(err)
		}
		_ = file.Sync()
	}

	cleanup := func() {
		// Clean up the testdata file after the test
		err := os.Remove(testFilePath)
		if err != nil {
			t.Fatal(err)
		}
	}

	return cleanup
}

func assertFileContent(t *testing.T, expectedContent string) {
	// Read the file and verify its content
	fileContent, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(fileContent) != expectedContent {
		t.Errorf("wrong file content: got %q, want %q", string(fileContent), expectedContent)
	}
}
