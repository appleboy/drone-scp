// +build windows

package main

import (
	"strings"
)

func getRealPath(path string) string {
	return "/" + strings.Replace(strings.Replace(path, ":", "", -1), "\\", "/", -1)
}
