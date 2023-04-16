package main

func rmcmd(os, target string) string {
	switch os {
	case "windows":
		return "DEL /F /S " + target
	case "unix":
		return "rm -rf " + target
	}
	return ""
}

func mkdircmd(os, target string) string {
	switch os {
	case "windows":
		return "if not exist " + target + " mkdir " + target
	case "unix":
		return "mkdir -p " + target
	}

	return ""
}
