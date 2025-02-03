package stores

import (
	"os"
	"strings"
	"testing"
)

func TestNewExportFile(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectedErr bool
	}{
		{
			name:        "valid path",
			path:        "test_file",
			expectedErr: false,
		},
		{
			name:        "no file name",
			path:        "",
			expectedErr: true,
		},
		{
			name:        "special characters in path",
			path:        "test@file#$.txt",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exportFile, err := NewExportFile(tt.path)
			if tt.expectedErr {
				if err == nil {
					t.Fatalf("expected an error but didn't get one")
				} else {
					return
				}
			}

			if exportFile == nil {
				t.Fatalf("expected non-nil ExportFile, got nil")
			}

			if exportFile.filePath != tt.path {
				t.Fatalf("expected fileName %q, got %q", tt.path, exportFile.filePath)
			}
		})
	}
}

func TestNewExportFile_Set(t *testing.T) {
	tests := []struct {
		name          string
		secretName    string
		expectedEntry string
		expectedErr   bool
	}{
		{
			name:          "valid secret name",
			secretName:    "SECRET",
			expectedEntry: "export SECRET=$(jangle get SECRET)\n",
			expectedErr:   false,
		},
		{
			name:          "no secret name",
			secretName:    "",
			expectedEntry: "",
			expectedErr:   true,
		},
		{
			name:          "special characters in secret name",
			secretName:    "GH_TOKEN()&*12-34",
			expectedEntry: "export GH_TOKEN()&*12-34=$(jangle get GH_TOKEN()&*12-34)\n",
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := "testdata/jangle_exports"

			exportFile, err := NewExportFile(fileName)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(fileName)

			err = exportFile.Set(tt.secretName)
			if tt.expectedErr {
				if err == nil {
					t.Fatalf("expected an error but didn't get one")
				} else {
					return
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			content, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if string(content) != tt.expectedEntry {
				t.Fatalf("expected file content %q, got %q", tt.expectedEntry, string(content))
			}
		})
	}
}

func TestNewExportFile_Delete(t *testing.T) {
	tests := []struct {
		name            string
		initialEntries  []string
		deleteName      string
		expectedContent string
		expectedErr     bool
	}{
		{
			name:            "existing entry",
			initialEntries:  []string{"export SECRET=$(jangle get SECRET)\n"},
			deleteName:      "SECRET",
			expectedContent: "",
			expectedErr:     false,
		},
		{
			name:            "non-existing entry",
			initialEntries:  []string{"export SECRET=$(jangle get SECRET)\n"},
			deleteName:      "NON_EXISTENT",
			expectedContent: "export SECRET=$(jangle get SECRET)\n",
			expectedErr:     false,
		},
		{
			name:            "empty name",
			initialEntries:  []string{"export SECRET=$(jangle get SECRET)\n"},
			deleteName:      "",
			expectedContent: "export SECRET=$(jangle get SECRET)\n",
			expectedErr:     false,
		},
		{
			name:            "special characters in entry",
			initialEntries:  []string{"export GH_TOKEN@123=$(jangle get GH_TOKEN@123)\n"},
			deleteName:      "GH_TOKEN@123",
			expectedContent: "",
			expectedErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := "testdata/jangle_exports"
			defer os.Remove(fileName)

			// Setup: write initial entries to the file
			if len(tt.initialEntries) > 0 {
				content := strings.Join(tt.initialEntries, "")
				err := os.WriteFile(fileName, []byte(content), 0644)
				if err != nil {
					t.Fatalf("failed to set up test file: %v", err)
				}
			}

			// Create ExportFile instance
			exportFile, err := NewExportFile(fileName)
			if err != nil {
				t.Fatal(err)
			}

			// Call Delete method
			err = exportFile.Delete(tt.deleteName)
			if tt.expectedErr {
				if err == nil {
					t.Fatalf("expected an error but didn't get one")
				} else {
					return
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify content
			content, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if string(content) != tt.expectedContent {
				t.Fatalf("expected file content %q, got %q", tt.expectedContent, string(content))
			}
		})
	}
}
