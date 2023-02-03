//go:build !windows
// +build !windows

package main

func rmcmd(target string) string {
	return "rm -rf " + target
}

func mkdircmd(target string) string {
	return "mkdir -p " + target
}
