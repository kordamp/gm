package gum

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func mvnCmd(executable string, buildFile string, args []string) Command {
	cmd := Command{Executable: executable}
	cmd.BuildFile = buildFile
	cmd.Args = args
	return cmd
}

// Finds and executes mvnw/mvn
func FindMaven(quiet bool, explicit bool, args []string) Command {
	pwd := GetWorkingDir()

	mvnw, noWrapper := findMavenWrapperExec(pwd)
	mvn, noMaven := findMavenExec()
	explicitBuildFileSet, explicitBuildFile := findExplicitMavenBuildFile(args)

	var executable string
	if noWrapper == nil {
		executable = mvnw
	} else if noMaven == nil {
		if !quiet {
			fmt.Printf("No %s set up for this project. ", resolveMavenWrapperExec())
			fmt.Println("Please consider setting one up.")
			fmt.Println("(https://maven.apache.org/)")
			fmt.Println()
		}
		executable = mvn
	} else {
		if !quiet {
			fmt.Printf("No %s found in path. Please install Maven.", resolveMavenExec())
			fmt.Println("(https://maven.apache.org/download.cgi)")
			fmt.Println()
		}

		if explicit {
			os.Exit(-1)
		} else {
			return EmptyCmd()
		}
	}

	if explicitBuildFileSet {
		return mvnCmd(executable, explicitBuildFile, args)
	}

	rootBuildFile, _ := findMavenRootFile(pwd, args)
	_, noBuildFile := findMavenBuildFile(pwd, args)

	if noBuildFile != nil {
		if explicit {
			fmt.Println("No Maven project found.")
			fmt.Println()
			os.Exit(-1)
		} else {
			return EmptyCmd()
		}
	}

	return mvnCmd(executable, rootBuildFile, args)
}

// Finds the maven executable
func findMavenExec() (string, error) {
	maven := resolveMavenExec()
	paths := GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], maven)
		if FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(maven + " not found")
}

// Finds the Maven wrapper (if it exists)
func findMavenWrapperExec(dir string) (string, error) {
	wrapper := resolveMavenWrapperExec()
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New(wrapper + " not found")
	}

	path := filepath.Join(dir, wrapper)
	if FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenWrapperExec(parentdir)
}

func findExplicitMavenBuildFile(args []string) (bool, string) {
	found, buildFile := FindFlag("-f", args)
	if !found {
		found, buildFile = FindFlag("--file", args)
	}

	if found {
		return true, buildFile
	}

	return false, ""
}

// Finds the nearest pom.xml
func findMavenBuildFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenBuildFile(parentdir, args)
}

// Finds the root pom.xml
func findMavenRootFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root pom.xml")
	}

	currentPom := filepath.Join(dir, "pom.xml")
	parentPom := filepath.Join(parentdir, "pom.xml")
	if FileExists(currentPom) && !FileExists(parentPom) {
		return filepath.Abs(currentPom)
	}

	return findGradleBuildFile(parentdir, args)
}

// Resolves the mvnw executable (OS dependent)
func resolveMavenWrapperExec() string {
	if IsWindows() {
		return "mvnw.bat"
	}
	return "mvnw"
}

// Resolves the mvn executable (OS dependent)
func resolveMavenExec() string {
	if IsWindows() {
		return "mvn.bat"
	}
	return "mvn"
}
