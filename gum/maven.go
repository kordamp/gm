// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2020 Andres Almiray.
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
	args              []string
	buildFile         string
	explicitBuildFile string
	rootBuildFile     string
}

// Execute executes the given command
func (c MavenCommand) Execute() {
	c.doConfigureMaven()
	c.doExecuteMaven()
}

func (c *MavenCommand) doConfigureMaven() {
	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using maven at '"+c.executable+"'")
	nearest, oargs := GrabFlag("-gn", c.args)
	debug, oargs := GrabFlag("-gd", oargs)
	replaceSet := findFlag("-gr", args)
	skipReplace, oargs := GrabFlag("-gr", oargs)

	c.config.setDebug(debug)
	if replaceSet {
		c.config.maven.setReplace(!skipReplace)
	}
	rargs := replaceMavenGoals(c.config, oargs)

	if len(c.explicitBuildFile) > 0 {
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

	for i := range rargs {
		args = append(args, rargs[i])
	}
	c.args = args

	c.debugMaven(nearest, oargs, rargs, args)

	if !c.config.general.quiet {
		fmt.Println(strings.Join(banner, " "))
	}
}

func (c *MavenCommand) doExecuteMaven() {
	cmd := exec.Command(c.executable, c.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func (c *MavenCommand) debugMaven(nearest bool, oargs []string, rargs []string, args []string) {
	if c.config.general.debug {
		fmt.Println("nearest            = ", nearest)
		fmt.Println("rootBuildFile      = ", c.rootBuildFile)
		fmt.Println("buildFile          = ", c.buildFile)
		fmt.Println("explicitBuildFile  = ", c.explicitBuildFile)
		fmt.Println("original args      = ", oargs)
		if c.config.maven.replace {
			fmt.Println("replaced args      = ", rargs)
		}
		fmt.Println("actual args        = ", args)
		fmt.Println("")
	}
}

func replaceMavenGoals(config *Config, args []string) []string {
	var nargs []string = args

	if config.maven.replace {
		nargs = replaceArgs(args, config.maven.mappings)
	}

	return nargs
}

// FindMaven finds and executes mvnw/mvn
func FindMaven(context Context, args []string) *MavenCommand {
	pwd := context.GetWorkingDir()

	mvnw, noWrapper := findMavenWrapperExec(context, pwd)
	mvn, noMaven := findMavenExec(context)
	explicitBuildFileSet, explicitBuildFile := findExplicitMavenBuildFile(args)

	rootBuildFile, noRootBuildFile := findMavenRootFile(context, filepath.Join(pwd, ".."), args)
	buildFile, noBuildFile := findMavenBuildFile(context, pwd, args)
	rootdir := resolveMavenRootDir(context, explicitBuildFile, buildFile, rootBuildFile)
	config := ReadConfig(context, rootdir)
	config.setQuiet(context.IsQuiet())

	var executable string
	if noWrapper == nil {
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
		fmt.Println("Please consider setting one up.")
		fmt.Println("(https://maven.apache.org/)")
		fmt.Println()
	}
}

func warnNoMaven(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s found in path. Please install Maven.", resolveMavenExec(context))
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

func findExplicitMavenBuildFile(args []string) (bool, string) {
	found, file := findFlagValue("-f", args)
	if !found {
		found, file = findFlagValue("--file", args)
	}

	if found {
		file, _ = filepath.Abs(file)
		return true, file
	}

	return false, ""
}

// Finds the nearest pom.xml
func findMavenBuildFile(context Context, dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenBuildFile(context, parentdir, args)
}

// Finds the root pom.xml
func findMavenRootFile(context Context, dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenRootFile(context, parentdir, args)
}

// Resolves the mvnw executable (OS dependent)
func resolveMavenWrapperExec(context Context) string {
	if context.IsWindows() {
		return "mvnw.bat"
	}
	return "mvnw"
}

// Resolves the mvn executable (OS dependent)
func resolveMavenExec(context Context) string {
	if context.IsWindows() {
		return "mvn.bat"
	}
	return "mvn"
}
