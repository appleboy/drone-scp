//go:build windows
// +build windows

package main

func rmcmd(target string) string {
	return "DEL /F /S " + target
}

func mkdircmd(target string) string {
	return "if not exist " + target + " mkdir " + target
}
