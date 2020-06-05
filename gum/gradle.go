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
	explicitProjectDir   string
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
	debug, nargs := GrabFlag("-gd", nargs)

	if debug {
		fmt.Println("nearest              = ", nearest)
		fmt.Println("args                 = ", nargs)
		fmt.Println("rootBuildFile        = ", c.rootBuildFile)
		fmt.Println("buildFile            = ", c.buildFile)
		fmt.Println("settingsFile         = ", c.settingsFile)
		fmt.Println("explicitBuildFile    = ", c.explicitBuildFile)
		fmt.Println("explicitSettingsFile = ", c.explicitSettingsFile)
		fmt.Println("explicitProjectDir   = ", c.explicitProjectDir)
		fmt.Println("")
	}

	if len(c.explicitProjectDir) > 0 {
		banner = append(banner, "to run project at '"+c.explicitProjectDir+"':")
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
		} else if len(c.rootBuildFile) > 0 {
			args = append(args, "-b")
			args = append(args, c.rootBuildFile)
			banner = append(banner, "to run buildFile '"+c.rootBuildFile+"':")
			buildFileSet = true
		}

		if len(c.explicitSettingsFile) > 0 {
			if !buildFileSet {
				banner = append(banner, "with settings at '"+c.explicitSettingsFile+"':")
			}
		} else if len(c.settingsFile) > 0 {
			args = append(args, "-c")
			args = append(args, c.settingsFile)
			if !buildFileSet {
				banner = append(banner, "with settings at '"+c.settingsFile+"':")
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
	explicitProjectDirSet, explicitProjectDir := findExplicitProjectDir(args)

	gradlew, noWrapper := resolveGradleWrapperExecutable(args)
	explicitBuildFileSet, explicitBuildFile := findExplicitGradleBuildFile(args)
	explicitSettingsFileSet, explicitSettingsFile := findExplicitGradleSettingsFile(args)
	settingsFile, noSettings := findGradleSettingsFile(pwd, args)
	buildFile, noBuildFile := findGradleBuildFile(pwd, args)

	var executable string
	if noWrapper == nil {
		executable = gradlew
	} else if noGradle == nil {
		warnNoGradleWrapper(quiet, explicit)
		executable = gradle
	} else {
		warnNoGradle(quiet, explicit)

		if explicit {
			os.Exit(-1)
		}
		return nil
	}

	if explicitProjectDirSet {
		return gradleCommand{
			quiet:              quiet,
			executable:         executable,
			args:               args,
			explicitProjectDir: explicitProjectDir}
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

	rootBuildFile, noRootBuildFile := findGradleRootFile(filepath.Join(pwd, ".."), args)

	if noRootBuildFile != nil {
		rootBuildFile = buildFile
	}

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
				fmt.Println("No Gradle project found")
				fmt.Println()
				os.Exit(-1)
			}
			return nil
		}
	}

	return gradleCommand{
		quiet:                quiet,
		executable:           executable,
		args:                 args,
		buildFile:            buildFile,
		rootBuildFile:        rootBuildFile,
		settingsFile:         settingsFile,
		explicitSettingsFile: explicitSettingsFile}
}

func resolveGradleWrapperExecutable(args []string) (string, error) {
	pwd := getWorkingDir()
	projectDirSet, projectDir := findExplicitProjectDir(args)

	if projectDirSet {
		return findGradleWrapperExec(projectDir)
	}
	return findGradleWrapperExec(pwd)
}

func warnNoGradleWrapper(quiet bool, explicit bool) {
	if !quiet && explicit {
		fmt.Printf("No %s set up for this project. ", resolveGradleWrapperExec())
		fmt.Println("Please consider setting one up.")
		fmt.Println("(https://gradle.org/docs/current/userguide/gradle_wrapper.html)")
		fmt.Println()
	}
}

func warnNoGradle(quiet bool, explicit bool) {
	if !quiet && explicit {
		fmt.Printf("No %s found in path. Please install Gradle.", resolveGradleExec())
		fmt.Println("(https://gradle.org/docs/current/userguide/installation.html)")
		fmt.Println()
	}
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
		file, _ = filepath.Abs(file)
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
		file, _ = filepath.Abs(file)
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
		file, _ = filepath.Abs(file)
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
		path := filepath.Join(dir, buildFiles[i])
		if fileExists(path) {
			return filepath.Abs(path)
		}
	}

	return findGradleRootFile(parentdir, args)
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
