package main

// This function returns the appropriate command for removing a file/directory based on the operating system.
func rmcmd(os, target string) string {
	switch os {
	case "windows":
		// On Windows, use DEL command to delete files and folders recursively
		return "DEL /F /S " + target
	case "unix":
		// On Unix-based systems, use rm command to delete files and folders recursively
		return "rm -rf " + target
	}
	// Return an empty string if the operating system is not recognized
	return ""
}

// This function returns the appropriate command for creating a directory based on the operating system.
func mkdircmd(os, target string) string {
	switch os {
	case "windows":
		// On Windows, use mkdir command to create directory and check if it exists
		return "if not exist " + target + " mkdir " + target
	case "unix":
		// On Unix-based systems, use mkdir command with -p option to create directories recursively
		return "mkdir -p " + target
	}
	// Return an empty string if the operating system is not recognized
	return ""
}
