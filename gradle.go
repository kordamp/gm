package gum

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func gradleCmd(executable string, buildFile string, settings string, args []string) Command {
	cmd := Command{Executable: executable}
	cmd.BuildFile = buildFile
	cmd.Settings = settings
	cmd.Args = args
	return cmd
}

// Finds and executes gradlew/gradle
func FindGradle(quiet bool, explicit bool, args []string) Command {
	pwd := GetWorkingDir()

	gradlew, noWrapper := findGradleWrapperExec(pwd)
	gradle, noGradle := findGradleExec()
	explicitBuildFileSet, explicitBuildFile := findExplicitGradleBuildFile(args)

	var executable string
	if noWrapper == nil {
		executable = gradlew
	} else if noGradle == nil {
		if !quiet {
			fmt.Printf("No %s set up for this project. ", resolveGradleWrapperExec())
			fmt.Println("Please consider setting one up.")
			fmt.Println("(https://gradle.org/docs/current/userguide/gradle_wrapper.html)")
			fmt.Println()
		}
		executable = gradle
	} else {
		if !quiet {
			fmt.Printf("No %s found in path. Please install Gradle.", resolveGradleExec())
			fmt.Println("(https://gradle.org/docs/current/userguide/installation.html)")
			fmt.Println()
		}

		if explicit {
			os.Exit(-1)
		}
	}

	if explicitBuildFileSet {
		return gradleCmd(executable, explicitBuildFile, "", args)
	}

	settingsFile, noSettings := findSettingsFile(pwd, args)
	buildFile, noBuildFile := findGradleBuildFile(pwd, args)

	if noBuildFile != nil {
		if noSettings == nil {
			if !quiet {
				fmt.Printf("Did not find a suitable Gradle build file but found %s", settingsFile)
				fmt.Println()
			}
		} else {
			if explicit {
				fmt.Println("No Gradle project found.")
				fmt.Println()
				os.Exit(-1)
			} else {
				return EmptyCmd()
			}
		}
	}

	return gradleCmd(executable, buildFile, settingsFile, args)
}

// Finds the gradle executable
func findGradleExec() (string, error) {
	gradle := resolveGradleExec()
	paths := GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], gradle)
		if FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(gradle + " not found")
}

// Finds the gradle wrapper (if it exists)
func findGradleWrapperExec(dir string) (string, error) {
	wrapper := resolveGradleWrapperExec()
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New(wrapper + " not found")
	}

	path := filepath.Join(dir, wrapper)
	if FileExists(path) {
		return filepath.Abs(path)
	}

	return findGradleWrapperExec(parentdir)
}

func findExplicitGradleBuildFile(args []string) (bool, string) {
	found, buildFile := FindFlag("-b", args)
	if !found {
		found, buildFile = FindFlag("--build-file", args)
	}

	if found {
		return true, buildFile
	}

	return false, ""
}

// Finds the nearest Gradle build file
// Unless explicit -b buildFile is given in args
// Checks the following paths in order:
// - build.gradle
// - build.gradle.kts
// - ${basedir}.gradle
// - ${basedir}.gradle.kts
func findGradleBuildFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find Gradle build file")
	}

	var buildFiles [4]string
	buildFiles[0] = "build.gradle"
	buildFiles[1] = "build.gradle.kts"
	buildFiles[2] = filepath.Base(dir) + ".gradle"
	buildFiles[3] = filepath.Base(dir) + ".gradle.kts"

	for i := range buildFiles {
		path := filepath.Join(dir, buildFiles[i])
		if FileExists(path) {
			return filepath.Abs(path)
		}
	}

	return findGradleBuildFile(parentdir, args)
}

// Finds settings.gradle(.kts)
// Unless explicit -c settingsFile is given in args
func findSettingsFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find Gradle settings file")
	}

	var settingsFiles [2]string
	settingsFiles[0] = "settings.gradle"
	settingsFiles[1] = "settings.gradle.kts"

	for i := range settingsFiles {
		path := filepath.Join(dir, settingsFiles[i])
		if FileExists(path) {
			return filepath.Abs(path)
		}
	}

	return findGradleBuildFile(parentdir, args)
}

// Resolves the gradlew executable (OS dependent)
func resolveGradleWrapperExec() string {
	if IsWindows() {
		return "gradlew.bat"
	}
	return "gradlew"
}

// Resolves the gradle executable (OS dependent)
func resolveGradleExec() string {
	if IsWindows() {
		return "gradle.bat"
	}
	return "gradle"
}
