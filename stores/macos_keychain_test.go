package stores

import (
	"testing"
)

func TestNewMacOSKeychainStore(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedNamespace string
	}{
		{
			name:              "when no input is passed the namespace should be the default",
			input:             "",
			expectedNamespace: "jangle_",
		},
		{
			name:              "if an input is passed the namespace should be the input",
			input:             "prefix_",
			expectedNamespace: "prefix_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var store *MacOSKeychainStore

			if tt.input == "" {
				store = NewMacOSKeychainStore()
			} else {
				store = NewMacOSKeychainStore(tt.input)
			}

			if store.namespace != tt.expectedNamespace {
				t.Errorf("expected namespace to be: %s, got %s", tt.expectedNamespace, store.namespace)
			}
		})
	}
}
