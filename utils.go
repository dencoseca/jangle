package main

import "strings"

const NAMESPACE_PREFIX = "jangle_"

func addPrefix(key string) string {
	return NAMESPACE_PREFIX + key
}

func stripPrefix(key string) string {
	return strings.TrimPrefix(key, NAMESPACE_PREFIX)
}
