package gum

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type gradleCommand struct {
	quiet                bool
	executable           string
	args                 []string
	projectDir           string
	buildFile            string
	explicitBuildFile    string
	rootBuildFile        string
	settingsFile         string
	explicitSettingsFile string
}

func (c gradleCommand) Execute() {
	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using gradle at '"+c.executable+"'")
	nearest, nargs := GrabFlag("-gn", c.args)

	if len(c.projectDir) > 0 {
		banner = append(banner, "to run project at '"+c.projectDir+"':")
	} else {
		var buildFileSet bool
		if len(c.explicitBuildFile) > 0 {
			banner = append(banner, "to run buildFile '"+c.explicitBuildFile+"':")
			buildFileSet = true
		} else if nearest && len(c.buildFile) > 0 {
			args = append(args, "-b")
			args = append(args, c.buildFile)
			banner = append(banner, "to run buildFile '"+c.buildFile+"':")
			buildFileSet = true
		} else {
			args = append(args, "-b")
			args = append(args, c.rootBuildFile)
			banner = append(banner, "to run buildFile '"+c.rootBuildFile+"':")
			buildFileSet = true
		}

		if len(c.settingsFile) > 0 {
			args = append(args, "-c")
			args = append(args, c.settingsFile)
			if !buildFileSet {
				banner = append(banner, "with settings at '"+c.settingsFile+"':")
			}
		} else if len(c.explicitSettingsFile) > 0 {
			if !buildFileSet {
				banner = append(banner, "with settings at '"+c.explicitSettingsFile+"':")
			}
		}
	}

	for i := range nargs {
		args = append(args, nargs[i])
	}

	if !c.quiet {
		fmt.Println(strings.Join(banner, " "))
	}

	cmd := exec.Command(c.executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// FindGradle finds and executes gradlew/gradle
func FindGradle(quiet bool, explicit bool, args []string) Command {
	pwd := getWorkingDir()

	gradle, noGradle := findGradleExec()
	projectDirSet, projectDir := findExplicitProjectDir(args)

	var gradlew string
	var noWrapper error
	if projectDirSet {
		gradlew, noWrapper = findGradleWrapperExec(projectDir)
	} else {
		gradlew, noWrapper = findGradleWrapperExec(pwd)
	}

	explicitBuildFileSet, explicitBuildFile := findExplicitGradleBuildFile(args)
	explicitSettingsFileSet, explicitSettingsFile := findExplicitGradleSettingsFile(args)
	settingsFile, noSettings := findGradleSettingsFile(pwd, args)
	buildFile, noBuildFile := findGradleBuildFile(pwd, args)

	var executable string
	if noWrapper == nil {
		executable = gradlew
	} else if noGradle == nil {
		if !quiet && explicit {
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
		} else {
			return nil
		}
	}

	if projectDirSet {
		return gradleCommand{
			quiet:      quiet,
			executable: executable,
			args:       args,
			projectDir: projectDir}
	}

	if explicitBuildFileSet {
		if explicitSettingsFileSet {
			return gradleCommand{
				quiet:                quiet,
				executable:           executable,
				args:                 args,
				explicitBuildFile:    explicitBuildFile,
				explicitSettingsFile: explicitSettingsFile}
		}
		return gradleCommand{
			quiet:             quiet,
			executable:        executable,
			args:              args,
			explicitBuildFile: explicitBuildFile,
			settingsFile:      settingsFile}
	}

	rootBuildFile, _ := findGradleRootFile(pwd, args)

	if noBuildFile != nil {
		if explicitSettingsFileSet {
			if !quiet {
				fmt.Printf("Did not find a suitable Gradle build file but %s is specified", explicitSettingsFile)
				fmt.Println()
			}
			return gradleCommand{
				quiet:                quiet,
				executable:           executable,
				args:                 args,
				buildFile:            buildFile,
				rootBuildFile:        rootBuildFile,
				explicitSettingsFile: explicitSettingsFile}
		} else if noSettings == nil {
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
				return nil
			}
		}
	}

	return gradleCommand{
		quiet:         quiet,
		executable:    executable,
		args:          args,
		buildFile:     buildFile,
		rootBuildFile: rootBuildFile,
		settingsFile:  settingsFile}
}

// Finds the gradle executable
func findGradleExec() (string, error) {
	gradle := resolveGradleExec()
	paths := getPaths()

	for i := range paths {
		name := filepath.Join(paths[i], gradle)
		if fileExists(name) {
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
	if fileExists(path) {
		return filepath.Abs(path)
	}

	return findGradleWrapperExec(parentdir)
}

func findExplicitProjectDir(args []string) (bool, string) {
	found, file := findFlag("-p", args)
	if !found {
		found, file = findFlag("--project-dir", args)
	}

	if found {
		return true, file
	}

	return false, ""
}

func findExplicitGradleBuildFile(args []string) (bool, string) {
	found, file := findFlag("-b", args)
	if !found {
		found, file = findFlag("--build-file", args)
	}

	if found {
		return true, file
	}

	return false, ""
}

func findExplicitGradleSettingsFile(args []string) (bool, string) {
	found, file := findFlag("-c", args)
	if !found {
		found, file = findFlag("--settings-file", args)
	}

	if found {
		return true, file
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
		if fileExists(path) {
			return filepath.Abs(path)
		}
	}

	return findGradleBuildFile(parentdir, args)
}

// Finds settings.gradle(.kts)
// Unless explicit -c settingsFile is given in args
func findGradleSettingsFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find Gradle settings file")
	}

	var settingsFiles [2]string
	settingsFiles[0] = "settings.gradle"
	settingsFiles[1] = "settings.gradle.kts"

	for i := range settingsFiles {
		path := filepath.Join(dir, settingsFiles[i])
		if fileExists(path) {
			return filepath.Abs(path)
		}
	}

	return findGradleSettingsFile(parentdir, args)
}

// Finds the root build file
func findGradleRootFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root build file")
	}

	var buildFiles [2]string
	buildFiles[0] = "build.gradle"
	buildFiles[1] = "build.gradle.kts"

	for i := range buildFiles {
		currentBuild := filepath.Join(dir, buildFiles[i])
		parentBuild := filepath.Join(parentdir, buildFiles[i])
		if fileExists(currentBuild) && !fileExists(parentBuild) {
			return filepath.Abs(currentBuild)
		}
	}

	return findGradleBuildFile(parentdir, args)
}

// Resolves the gradlew executable (OS dependent)
func resolveGradleWrapperExec() string {
	if isWindows() {
		return "gradlew.bat"
	}
	return "gradlew"
}

// Resolves the gradle executable (OS dependent)
func resolveGradleExec() string {
	if isWindows() {
		return "gradle.bat"
	}
	return "gradle"
}
