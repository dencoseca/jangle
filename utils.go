package main

import "strings"

// NAMESPACE_PREFIX defines a default prefix used for namespacing keys to ensure
// uniqueness and avoid collisions.
const NAMESPACE_PREFIX = "jangle_"

// addPrefix appends a default namespace prefix to the provided key to ensure
// uniqueness and avoid naming collisions.
func addPrefix(key string) string {
	return NAMESPACE_PREFIX + key
}

// stripPrefix removes the predefined namespace prefix from the provided key and
// returns the resulting string.
func stripPrefix(key string) string {
	return strings.TrimPrefix(key, NAMESPACE_PREFIX)
}

// isHelpArg checks if the provided arguments contain exactly one element and if
// that element is "--help" or "help".
func isHelpArg(args []string) bool {
	return len(args) == 1 && (args[0] == "--help" || args[0] == "help")
}
