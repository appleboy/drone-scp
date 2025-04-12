//go:build windows
// +build windows

package main

import (
	"strings"
)

func getRealPath(path string) string {
	return "/" + strings.ReplaceAll(strings.ReplaceAll(path, ":", ""), "\\", "/")
}
