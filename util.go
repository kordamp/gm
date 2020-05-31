package gum

import (
	"os"
	"runtime"
	"strings"
)

func IsFlagSet(f string, args []string) (bool, []string) {
	if len(args) == 0 {
		// no args to be checked
		return false, args
	}

	for i := range args {
		s := args[i]
		if s == f {
			newArgs := make([]string, len(args)-1)
			for j := 0; j < i; j++ {
				newArgs[j] = args[j]
			}
			for j := i; j < len(newArgs); j++ {
				newArgs[j] = args[j+1]
			}
			return true, newArgs
		}
	}

	return false, args
}

func FindFlag(flag string, args []string) (bool, string) {
	if len(args) == 0 {
		return false, ""
	}

	for i := range args {
		s := args[i]
		if flag == s {
			// next argument should contain the value we want
			if i+1 < len(args) {
				return true, args[i+1]
			} else {
				return false, ""
			}
		} else {
			// check if format is flag=value
			parts := strings.Split(s, "=")
			if len(parts) == 2 && parts[0] == flag {
				return true, parts[1]
			}
		}
	}

	return false, ""
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// Returns the current working dir
func GetWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return pwd
}

// Gets the paths in $PATH
func GetPaths() []string {
	return strings.Split(getPathFromEnv(), string(os.PathListSeparator))
}

// Gets the PATH environment variable
func getPathFromEnv() string {
	if IsWindows() {
		return os.Getenv("Path")
	}

	return os.Getenv("PATH")
}

// Checks if a file exists
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
