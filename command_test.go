package main

import "testing"

// Unit tests for rmcmd and mkdircmd
func TestCommands(t *testing.T) {
	// Test rmcmd on Windows
	os1 := "windows"
	target1 := "C:\\path\\to\\file"
	expected1 := "DEL /F /S " + target1
	actual1 := rmcmd(os1, target1)
	if actual1 != expected1 {
		t.Errorf("rmcmd(%s, %s) = %s; expected %s", os1, target1, actual1, expected1)
	}

	// Test rmcmd on Unix-based system
	os2 := "unix"
	target2 := "/path/to/folder"
	expected2 := "rm -rf " + target2
	actual2 := rmcmd(os2, target2)
	if actual2 != expected2 {
		t.Errorf("rmcmd(%s, %s) = %s; expected %s", os2, target2, actual2, expected2)
	}

	// Test mkdircmd on Windows
	os3 := "windows"
	target3 := "C:\\path\\to\\folder"
	expected3 := "if not exist " + target3 + " mkdir " + target3
	actual3 := mkdircmd(os3, target3)
	if actual3 != expected3 {
		t.Errorf("mkdircmd(%s, %s) = %s; expected %s", os3, target3, actual3, expected3)
	}

	// Test mkdircmd on Unix-based system
	os4 := "unix"
	target4 := "/path/to/folder"
	expected4 := "mkdir -p " + target4
	actual4 := mkdircmd(os4, target4)
	if actual4 != expected4 {
		t.Errorf("mkdircmd(%s, %s) = %s; expected %s", os4, target4, actual4, expected4)
	}
}
