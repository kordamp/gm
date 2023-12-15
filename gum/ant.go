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

// AntCommand defines an executable Ant command
type AntCommand struct {
	context           Context
	config            *Config
	rootdir           string
	executable        string
	args              *ParsedArgs
	buildFile         string
	explicitBuildFile string
}

// Execute executes the given command
func (c AntCommand) Execute() int {
	c.doConfigureAnt()
	return c.doExecuteAnt()
}

func (c *AntCommand) doConfigureAnt() {
	c.context.CheckIsExecutable(c.executable)

	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using Ant at '"+c.executable+"'")

	debug := c.args.HasGumFlag("gd")

	if debug {
		c.config.setDebug(debug)
	}
	c.debugConfig()
	oargs := c.args.Args

	if len(c.explicitBuildFile) > 0 {
		args = append(args, "-f")
		args = append(args, c.explicitBuildFile)
		banner = append(banner, "to run buildFile '"+c.explicitBuildFile+"':")
	} else if len(c.buildFile) > 0 {
		args = append(args, "-f")
		args = append(args, c.buildFile)
		banner = append(banner, "to run buildFile '"+c.buildFile+"':")
	}

	args = appendSafe(args, c.args.Tool)
	args = append(args, "-Dbasedir="+c.rootdir)
	c.args.Args = appendSafe(args, oargs)

	c.debugAnt(c.config, oargs)

	if !c.config.general.quiet {
		fmt.Println(strings.Join(banner, " "))
	}
}

func (c *AntCommand) doExecuteAnt() int {
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

func (c *AntCommand) debugConfig() {
	if c.args.HasGumFlag("gc") {
		c.config.print()
		os.Exit(0)
	}
}

func (c *AntCommand) debugAnt(config *Config, oargs []string) {
	if c.config.general.debug {
		fmt.Println("rootdir            = ", c.rootdir)
		fmt.Println("executable         = ", c.executable)
		fmt.Println("buildFile          = ", c.buildFile)
		fmt.Println("explicitBuildFile  = ", c.explicitBuildFile)
		fmt.Println("original args      = ", oargs)
		fmt.Println("actual args        = ", c.args.Args)
		fmt.Println("")
	}
}

// FindAnt finds and executes Ant
func FindAnt(context Context, args *ParsedArgs) *AntCommand {
	pwd := context.GetWorkingDir()

	ant, noAnt := findAntExec(context)
	explicitBuildFileSet, explicitBuildFile := findExplicitAntBuildFile(args)
	buildFile, noBuildFile := findAntBuildFile(context, pwd)

	rootdir := resolveAntRootDir(context, explicitBuildFile, buildFile)
	config := ReadConfig(context, rootdir)
	quiet := args.HasGumFlag("gq")

	if quiet {
		config.setQuiet(quiet)
	}

	var executable string
	if noAnt == nil {
		executable = ant
	} else {
		warnNoAnt(context, config)

		if context.IsExplicit() {
			context.Exit(-1)
		}
		return nil
	}

	if explicitBuildFileSet {
		return &AntCommand{
			context:           context,
			config:            config,
			executable:        executable,
			args:              args,
			explicitBuildFile: explicitBuildFile}
	}

	if noBuildFile != nil {
		if context.IsExplicit() {
			fmt.Println("No Ant project found")
			fmt.Println()
			context.Exit(-1)
		}
		return nil
	}

	return &AntCommand{
		context:    context,
		config:     config,
		rootdir:    rootdir,
		executable: executable,
		args:       args,
		buildFile:  buildFile}
}

func resolveAntRootDir(context Context,
	explicitBuildFile string,
	buildFile string) string {

	if context.FileExists(explicitBuildFile) {
		return filepath.Dir(explicitBuildFile)
	}
	return filepath.Dir(buildFile)
}

func warnNoAnt(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s found in path. Please install Ant.", resolveAntExec(context))
		fmt.Println()
		fmt.Println("(https://ant.apache.org/bindownload.cgi)")
		fmt.Println()
	}
}

// Finds the ant executable
func findAntExec(context Context) (string, error) {
	ant := resolveAntExec(context)
	paths := context.GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], ant)
		if context.FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(ant + " not found")
}

func findExplicitAntBuildFile(args *ParsedArgs) (bool, string) {
	found, file, shrunkArgs := findFlagValue("-f", args.Tool)
	args.Tool = shrunkArgs
	if !found {
		found, file, shrunkArgs = findFlagValue("-file", args.Tool)
		args.Tool = shrunkArgs
	}
	if !found {
		found, file, shrunkArgs = findFlagValue("-buildfile", args.Tool)
		args.Tool = shrunkArgs
	}

	if found {
		file, _ = filepath.Abs(file)
		return true, file
	}

	return false, ""
}

// Finds the nearest build.xml
func findAntBuildFile(context Context, dir string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find build.xml")
	}

	path := filepath.Join(dir, "build.xml")
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findAntBuildFile(context, parentdir)
}

// Resolves the ant executable (OS dependent)
func resolveAntExec(context Context) string {
	if context.IsWindows() {
		return "ant.bat"
	}
	return "ant"
}
