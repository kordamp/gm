// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2020-2021 Andres Almiray.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gum

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GradleCommand defines an executable Gradle command
type GradleCommand struct {
	context              Context
	config               *Config
	executable           string
	args                 *ParsedArgs
	explicitProjectDir   string
	rootDir              string
	buildFile            string
	explicitBuildFile    string
	rootBuildFile        string
	settingsFile         string
	explicitSettingsFile string
}

// Execute executes the given command
func (c GradleCommand) Execute() {
	c.doConfigureGradle()
	c.doExecuteGradle()
}

func (c *GradleCommand) doConfigureGradle() {
	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using gradle at '"+c.executable+"'")
	nearest := c.args.HasGumFlag("gn")
	debug := c.args.HasGumFlag("gd")
	skipReplace := c.args.HasGumFlag("gr")

	if debug {
		c.config.setDebug(debug)
	}
	if skipReplace {
		c.config.gradle.setReplace(!skipReplace)
	}
	c.debugConfig()
	oargs := c.args.Args
	rargs := replaceGradleTasks(c.config, c.args)

	if len(c.explicitProjectDir) > 0 {
		banner = append(banner, "to run project at '"+c.explicitProjectDir+"':")
	} else {
		var buildFileSet bool
		if len(c.explicitBuildFile) > 0 {
			args = append(args, "-b")
			args = append(args, c.explicitBuildFile)
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
			pwd, _ := filepath.Abs(c.context.GetWorkingDir())
			settingsDir, _ := filepath.Abs(filepath.Dir(c.settingsFile))

			if c.rootDir != settingsDir || c.rootDir != pwd {
				args = append(args, "-c")
				args = append(args, c.settingsFile)
			}

			if !buildFileSet {
				banner = append(banner, "with settings at '"+c.settingsFile+"':")
			}
		}
	}

	args = appendSafe(args, c.args.Tool)
	c.args.Args = appendSafe(args, rargs)

	c.debugGradle(oargs, rargs)

	if !c.config.general.quiet {
		fmt.Println(strings.Join(banner, " "))
	}
}

func (c *GradleCommand) doExecuteGradle() {
	cmd := exec.Command(c.executable, c.args.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func (c *GradleCommand) debugConfig() {
	if c.args.HasGumFlag("gc") {
		c.config.print()
		os.Exit(0)
	}
}

func (c *GradleCommand) debugGradle(oargs []string, rargs []string) {
	if c.config.general.debug {
		fmt.Println("nearest              = ", c.args.HasGumFlag("gn"))
		fmt.Println("replace              = ", c.config.gradle.replace)
		fmt.Println("pwd                  = ", c.context.GetWorkingDir())
		fmt.Println("rootDir              = ", c.rootDir)
		fmt.Println("rootBuildFile        = ", c.rootBuildFile)
		fmt.Println("buildFile            = ", c.buildFile)
		fmt.Println("settingsFile         = ", c.settingsFile)
		fmt.Println("explicitBuildFile    = ", c.explicitBuildFile)
		fmt.Println("explicitSettingsFile = ", c.explicitSettingsFile)
		fmt.Println("explicitProjectDir   = ", c.explicitProjectDir)
		fmt.Println("original args        = ", oargs)
		if c.config.gradle.replace {
			fmt.Println("replaced args        = ", rargs)
		}
		fmt.Println("actual args          = ", c.args.Args)
		fmt.Println("")
	}
}

func replaceGradleTasks(config *Config, args *ParsedArgs) []string {
	if config.gradle.replace {
		return replaceArgs(args.Args, config.gradle.mappings, true)
	}

	return args.Args
}

// FindGradle finds and executes gradlew/gradle
func FindGradle(context Context, args *ParsedArgs) *GradleCommand {
	pwd := context.GetWorkingDir()

	gradle, noGradle := findGradleExec(context)
	explicitProjectDirSet, explicitProjectDir := findExplicitProjectDir(args)

	gradlew, noWrapper := resolveGradleWrapperExecutable(context, args)
	explicitBuildFileSet, explicitBuildFile := findExplicitGradleBuildFile(args)
	explicitSettingsFileSet, explicitSettingsFile := findExplicitGradleSettingsFile(args)
	settingsFile, noSettings := findGradleSettingsFile(context, pwd)
	buildFile, noBuildFile := findGradleBuildFile(context, pwd)

	sf := settingsFile
	if explicitBuildFileSet {
		sf = explicitBuildFile
	}

	rootBuildFile, noRootBuildFile := findGradleRootFile(context, filepath.Join(pwd, ".."), args, sf)
	rootdir := resolveGradleRootDir(context, explicitProjectDir, explicitBuildFile, explicitSettingsFile, buildFile, rootBuildFile, settingsFile)
	config := ReadConfig(context, rootdir)
	quiet := args.HasGumFlag("gq")

	if quiet {
		config.setQuiet(quiet)
	}

	var executable string
	if noWrapper == nil {
		executable = gradlew
	} else if noGradle == nil {
		warnNoGradleWrapper(context, config)
		executable = gradle
	} else {
		warnNoGradle(context, config)

		if context.IsExplicit() {
			context.Exit(-1)
		}
		return nil
	}

	if explicitProjectDirSet {
		return &GradleCommand{
			context:            context,
			config:             config,
			executable:         executable,
			args:               args,
			rootDir:            rootdir,
			explicitProjectDir: explicitProjectDir}
	}

	if explicitBuildFileSet {
		if explicitSettingsFileSet {
			return &GradleCommand{
				context:              context,
				config:               config,
				executable:           executable,
				args:                 args,
				rootDir:              rootdir,
				explicitBuildFile:    explicitBuildFile,
				explicitSettingsFile: explicitSettingsFile}
		}
		return &GradleCommand{
			context:           context,
			config:            config,
			executable:        executable,
			args:              args,
			rootDir:           rootdir,
			explicitBuildFile: explicitBuildFile,
			settingsFile:      settingsFile}
	}

	if noRootBuildFile != nil {
		rootBuildFile = buildFile
	}

	if noBuildFile != nil {
		if explicitSettingsFileSet {
			if !config.general.quiet {
				fmt.Printf("Did not find a suitable Gradle build file but %s is specified", explicitSettingsFile)
				fmt.Println()
			}
			return &GradleCommand{
				context:              context,
				config:               config,
				executable:           executable,
				args:                 args,
				rootDir:              rootdir,
				buildFile:            buildFile,
				rootBuildFile:        rootBuildFile,
				explicitSettingsFile: explicitSettingsFile}
		} else if noSettings == nil {
			if !config.general.quiet {
				fmt.Printf("Did not find a suitable Gradle build file but found %s", settingsFile)
				fmt.Println()
			}
		} else {
			if context.IsExplicit() {
				fmt.Println("No Gradle project found")
				fmt.Println()
				context.Exit(-1)
			}
			return nil
		}
	}

	return &GradleCommand{
		context:              context,
		config:               config,
		executable:           executable,
		args:                 args,
		rootDir:              rootdir,
		buildFile:            buildFile,
		rootBuildFile:        rootBuildFile,
		settingsFile:         settingsFile,
		explicitSettingsFile: explicitSettingsFile}
}

func resolveGradleRootDir(context Context,
	explicitProjectDir string,
	explicitBuildFile string,
	explicitSettingsFile string,
	buildFile string,
	rootBuildFile string,
	settingsFile string) string {

	if context.FileExists(explicitProjectDir) {
		return explicitProjectDir
	} else if context.FileExists(explicitBuildFile) {
		return filepath.Dir(explicitBuildFile)
	} else if context.FileExists(rootBuildFile) {
		return filepath.Dir(rootBuildFile)
	} else if context.FileExists(explicitSettingsFile) {
		return filepath.Dir(explicitSettingsFile)
	} else if context.FileExists(settingsFile) {
		return filepath.Dir(settingsFile)
	}

	dir, _ := filepath.Abs(filepath.Dir(buildFile))
	return dir
}

func resolveGradleWrapperExecutable(context Context, args *ParsedArgs) (string, error) {
	pwd := context.GetWorkingDir()
	projectDirSet, projectDir := findExplicitProjectDir(args)

	if projectDirSet {
		return findGradleWrapperExec(context, projectDir)
	}
	return findGradleWrapperExec(context, pwd)
}

func warnNoGradleWrapper(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s set up for this project. ", resolveGradleWrapperExec(context))
		fmt.Println("Please consider setting one up.")
		fmt.Println("(https://gradle.org/docs/current/userguide/gradle_wrapper.html)")
		fmt.Println()
	}
}

func warnNoGradle(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s found in path. Please install Gradle.", resolveGradleExec(context))
		fmt.Println("(https://gradle.org/docs/current/userguide/installation.html)")
		fmt.Println()
	}
}

// Finds the gradle executable
func findGradleExec(context Context) (string, error) {
	gradle := resolveGradleExec(context)
	paths := context.GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], gradle)
		if context.FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(gradle + " not found")
}

// Finds the gradle wrapper (if it exists)
func findGradleWrapperExec(context Context, dir string) (string, error) {
	wrapper := resolveGradleWrapperExec(context)
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New(wrapper + " not found")
	}

	path := filepath.Join(dir, wrapper)
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findGradleWrapperExec(context, parentdir)
}

func findExplicitProjectDir(args *ParsedArgs) (bool, string) {
	found, file, shrunkArgs := findFlagValue("-p", args.Tool)
	args.Tool = shrunkArgs
	if !found {
		found, file, shrunkArgs = findFlagValue("--project-dir", args.Tool)
		args.Tool = shrunkArgs
	}

	if found {
		file, _ = filepath.Abs(file)
		return true, file
	}

	return false, ""
}

func findExplicitGradleBuildFile(args *ParsedArgs) (bool, string) {
	found, file, shrunkArgs := findFlagValue("-b", args.Tool)
	args.Tool = shrunkArgs
	if !found {
		found, file, shrunkArgs = findFlagValue("--build-file", args.Tool)
		args.Tool = shrunkArgs
	}
	if !found {
		found, file, shrunkArgs = findFlagValue("-b", args.Args)
		args.Args = shrunkArgs
	}
	if !found {
		found, file, shrunkArgs = findFlagValue("--build-file", args.Args)
		args.Args = shrunkArgs
	}

	if found {
		file, _ = filepath.Abs(file)
		return true, file
	}

	return false, ""
}

func findExplicitGradleSettingsFile(args *ParsedArgs) (bool, string) {
	found, file, shrunkArgs := findFlagValue("-c", args.Tool)
	args.Tool = shrunkArgs
	if !found {
		found, file, shrunkArgs = findFlagValue("--settings-file", args.Tool)
		args.Tool = shrunkArgs
	}
	if !found {
		found, file, shrunkArgs = findFlagValue("-c", args.Args)
		args.Args = shrunkArgs
	}
	if !found {
		found, file, shrunkArgs = findFlagValue("--settings-file", args.Args)
		args.Args = shrunkArgs
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
func findGradleBuildFile(context Context, dir string) (string, error) {
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
		if context.FileExists(path) {
			return filepath.Abs(path)
		}
	}

	return findGradleBuildFile(context, parentdir)
}

// Finds settings.gradle(.kts)
// Unless explicit -c settingsFile is given in args
func findGradleSettingsFile(context Context, dir string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find Gradle settings file")
	}

	var settingsFiles [2]string
	settingsFiles[0] = "settings.gradle"
	settingsFiles[1] = "settings.gradle.kts"

	for i := range settingsFiles {
		path := filepath.Join(dir, settingsFiles[i])
		if context.FileExists(path) {
			return filepath.Abs(path)
		}
	}

	return findGradleSettingsFile(context, parentdir)
}

// Finds the root build file
func findGradleRootFile(context Context, dir string, args *ParsedArgs, settingsFile string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root build file")
	}

	var buildFiles [2]string
	buildFiles[0] = "build.gradle"
	buildFiles[1] = "build.gradle.kts"

	for i := range buildFiles {
		path := filepath.Join(dir, buildFiles[i])
		if context.FileExists(path) {
			return filepath.Abs(path)
		}
	}

	if len(settingsFile) > 0 {
		settingsdir := filepath.Dir(settingsFile)
		if len(parentdir) <= len(settingsdir) {
			return "", errors.New("Did not find root build file")
		}
	}

	return findGradleRootFile(context, parentdir, args, settingsFile)
}

// Resolves the gradlew executable (OS dependent)
func resolveGradleWrapperExec(context Context) string {
	if context.IsWindows() {
		return "gradlew.bat"
	}
	return "gradlew"
}

// Resolves the gradle executable (OS dependent)
func resolveGradleExec(context Context) string {
	if context.IsWindows() {
		return "gradle.bat"
	}
	return "gradle"
}
