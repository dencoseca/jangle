package stores

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewExportFile(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectedErr error
	}{
		{
			name:        "valid path",
			path:        "test_file.txt",
			expectedErr: nil,
		},
		{
			name:        "empty path",
			path:        "",
			expectedErr: nil,
		},
		{
			name:        "special characters in path",
			path:        "test@file#$.txt",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randomPrefix := fmt.Sprintf("%d_", time.Now().UnixNano())
			tt.path = filepath.Join(os.TempDir(), randomPrefix+tt.path)

			defer func() {
				if tt.path != "" {
					_ = os.Remove(tt.path)
				}
			}()

			exportFile := NewExportFile(tt.path)

			if exportFile == nil {
				t.Fatalf("Expected non-nil ExportFile, got nil")
			}

			if exportFile.fileName != tt.path {
				t.Errorf("Expected fileName %q, got %q", tt.path, exportFile.fileName)
			}
		})
	}
}
