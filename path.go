//go:build !windows
// +build !windows

package main

func getRealPath(path string) string {
	return path
}
