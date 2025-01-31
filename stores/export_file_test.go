package stores

import (
	"os"
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

			if exportFile.fileName != tt.path {
				t.Fatalf("expected fileName %q, got %q", tt.path, exportFile.fileName)
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
