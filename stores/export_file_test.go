package stores

import (
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
			path:        "test_file.txt",
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
