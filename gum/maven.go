// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2020-2023 Andres Almiray.
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

// MavenCommand defines an executable Maven command
type MavenCommand struct {
	context           Context
	config            *Config
	executable        string
	args              *ParsedArgs
	buildFile         string
	explicitBuildFile string
	rootBuildFile     string
}

// Execute executes the given command
func (c MavenCommand) Execute() int {
	c.doConfigureMaven()
	return c.doExecuteMaven()
}

func (c *MavenCommand) doConfigureMaven() {
	c.context.CheckIsExecutable(c.executable)

	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using maven at '"+c.executable+"'")
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
	otargs := c.args.Tool
	oargs := c.args.Args
	rtargs, rargs := replaceMavenGoals(c.config, c.args)

	if len(c.explicitBuildFile) > 0 {
		args = append(args, "-f")
		args = append(args, c.explicitBuildFile)
		banner = append(banner, "to run buildFile '"+c.explicitBuildFile+"':")
	} else if nearest && len(c.buildFile) > 0 {
		args = append(args, "-f")
		args = append(args, c.buildFile)
		banner = append(banner, "to run buildFile '"+c.buildFile+"':")
	} else if len(c.rootBuildFile) > 0 {
		args = append(args, "-f")
		args = append(args, c.rootBuildFile)
		banner = append(banner, "to run buildFile '"+c.rootBuildFile+"':")
	}

	args = appendSafe(args, rtargs)
	c.args.Args = appendSafe(args, rargs)

	c.debugMaven(otargs, oargs, rtargs, rargs)

	if !c.config.general.quiet {
		fmt.Println(strings.Join(banner, " "))
	}
}

func (c *MavenCommand) doExecuteMaven() int {
	cmd := exec.Command(c.executable, c.args.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	var exerr *exec.ExitError
	if errors.As(err, &exerr) {
		return exerr.ExitCode()
	}
	return 0
}

func (c *MavenCommand) debugConfig() {
	if c.args.HasGumFlag("gc") {
		c.config.print()
		os.Exit(0)
	}
}

func (c *MavenCommand) debugMaven(otargs []string, oargs []string, rtargs []string, rargs []string) {
	if c.config.general.debug {
		fmt.Println("nearest            = ", c.args.HasGumFlag("gn"))
		fmt.Println("replace            = ", c.config.maven.replace)
		fmt.Println("pwd                = ", c.context.GetWorkingDir())
		fmt.Println("rootBuildFile      = ", c.rootBuildFile)
		fmt.Println("buildFile          = ", c.buildFile)
		fmt.Println("explicitBuildFile  = ", c.explicitBuildFile)
		fmt.Println("original tool args = ", otargs)
		if c.config.maven.replace {
			fmt.Println("replaced tool args = ", rtargs)
		}
		fmt.Println("original args      = ", oargs)
		if c.config.maven.replace {
			fmt.Println("replaced args      = ", rargs)
		}
		fmt.Println("actual args        = ", c.args.Args)
		fmt.Println("")
	}
}

func replaceMavenGoals(config *Config, args *ParsedArgs) ([]string, []string) {
	if config.maven.replace {
		return replaceArgs(args.Tool, config.maven.mappings, false), replaceArgs(args.Args, config.maven.mappings, false)
	}

	return args.Tool, args.Args
}

// FindMaven finds and executes mvnw/mvn
func FindMaven(context Context, args *ParsedArgs) *MavenCommand {
	pwd := context.GetWorkingDir()

	mvnw, noWrapper := findMavenWrapperExec(context, pwd)
	mvn, noMaven := findMavenExec(context)
	mvnd, noMvnd := findMvndExec(context)
	explicitBuildFileSet, explicitBuildFile := findExplicitMavenBuildFile(args)

	rootBuildFile, noRootBuildFile := findMavenRootFile(context, filepath.Join(pwd, ".."))
	buildFile, noBuildFile := findMavenBuildFile(context, pwd)
	rootdir := resolveMavenRootDir(context, explicitBuildFile, buildFile, rootBuildFile)
	config := ReadConfig(context, rootdir)
	quiet := args.HasGumFlag("gq")

	if quiet {
		config.setQuiet(quiet)
	}

	var executable string
	if config.maven.mvnd && noMvnd == nil {
		executable = mvnd
	} else if noWrapper == nil {
		executable = mvnw
	} else if noMaven == nil {
		warnNoMavenWrapper(context, config)
		executable = mvn
	} else {
		warnNoMaven(context, config)

		if context.IsExplicit() {
			context.Exit(-1)
		}
		return nil
	}

	if explicitBuildFileSet {
		return &MavenCommand{
			context:           context,
			config:            config,
			executable:        executable,
			args:              args,
			explicitBuildFile: explicitBuildFile}
	}

	if noRootBuildFile != nil {
		rootBuildFile = buildFile
	}

	if noBuildFile != nil {
		if context.IsExplicit() {
			fmt.Println("No Maven project found")
			fmt.Println()
			context.Exit(-1)
		}
		return nil
	}

	return &MavenCommand{
		context:       context,
		config:        config,
		executable:    executable,
		args:          args,
		rootBuildFile: rootBuildFile,
		buildFile:     buildFile}
}

func resolveMavenRootDir(context Context,
	explicitBuildFile string,
	buildFile string,
	rootBuildFile string) string {

	if context.FileExists(explicitBuildFile) {
		return filepath.Dir(explicitBuildFile)
	} else if context.FileExists(rootBuildFile) {
		return filepath.Dir(rootBuildFile)
	}
	return filepath.Dir(buildFile)
}

func warnNoMavenWrapper(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s set up for this project. ", resolveMavenWrapperExec(context))
		fmt.Println()
		fmt.Println("Please consider setting one up.")
		fmt.Println("(https://maven.apache.org/)")
		fmt.Println()
	}
}

func warnNoMaven(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s found in path. Please install Maven.", resolveMavenExec(context))
		fmt.Println()
		fmt.Println("(https://maven.apache.org/download.cgi)")
		fmt.Println()
	}
}

// Finds the maven executable
func findMavenExec(context Context) (string, error) {
	maven := resolveMavenExec(context)
	paths := context.GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], maven)
		if context.FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(maven + " not found")
}

// Finds the mvnd executable
func findMvndExec(context Context) (string, error) {
	mvnd := resolveMvndExec(context)
	paths := context.GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], mvnd)
		if context.FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(mvnd + " not found")
}

// Finds the Maven wrapper (if it exists)
func findMavenWrapperExec(context Context, dir string) (string, error) {
	wrapper := resolveMavenWrapperExec(context)
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New(wrapper + " not found")
	}

	path := filepath.Join(dir, wrapper)
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenWrapperExec(context, parentdir)
}

func findExplicitMavenBuildFile(args *ParsedArgs) (bool, string) {
	found, file, shrunkArgs := findFlagValue("-f", args.Tool)
	args.Tool = shrunkArgs
	if !found {
		found, file, shrunkArgs = findFlagValue("--file", args.Tool)
		args.Tool = shrunkArgs
	}

	if found {
		file, _ = filepath.Abs(file)
		return true, file
	}

	return false, ""
}

// Finds the nearest pom.xml
func findMavenBuildFile(context Context, dir string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenBuildFile(context, parentdir)
}

// Finds the root pom.xml
func findMavenRootFile(context Context, dir string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenRootFile(context, parentdir)
}

// Resolves the mvnw executable (OS dependent)
func resolveMavenWrapperExec(context Context) string {
	if context.IsWindows() {
		return "mvnw.cmd"
	}
	return "mvnw"
}

// Resolves the mvn executable (OS dependent)
func resolveMavenExec(context Context) string {
	if context.IsWindows() {
		return "mvn.cmd"
	}
	return "mvn"
}

// Resolves the mvnd executable (OS dependent)
func resolveMvndExec(context Context) string {
	if context.IsWindows() {
		return "mvnd.cmd"
	}
	return "mvnd"
}
