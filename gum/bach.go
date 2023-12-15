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

// BachCommand defines an executable Bach command
type BachCommand struct {
	context    Context
	config     *Config
	rootdir    string
	executable string
	args       *ParsedArgs
}

// Execute executes the given command
func (c BachCommand) Execute() int {
	c.doConfigureBach()
	return c.doExecuteBach()
}

func (c *BachCommand) doConfigureBach() {
	execParts := strings.Split(c.executable, " ")
	c.context.CheckIsExecutable(execParts[0])

	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using Bach at '"+c.rootdir+"'")

	debug := c.args.HasGumFlag("gd")

	if debug {
		c.config.setDebug(debug)
	}
	c.debugConfig()
	oargs := c.args.Args

	c.executable = execParts[0]

	args = appendSafe(args, execParts[1:])
	args = appendSafe(args, c.args.Tool)
	c.args.Args = appendSafe(args, oargs)

	c.debugBach(c.config, oargs)

	if !c.config.general.quiet {
		fmt.Println(strings.Join(banner, " "))
	}
}

func (c *BachCommand) doExecuteBach() int {
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

func (c *BachCommand) debugConfig() {
	if c.args.HasGumFlag("gc") {
		c.config.print()
		os.Exit(0)
	}
}

func (c *BachCommand) debugBach(config *Config, oargs []string) {
	if c.config.general.debug {
		fmt.Println("rootdir            = ", c.rootdir)
		fmt.Println("executable         = ", c.executable)
		fmt.Println("original args      = ", oargs)
		fmt.Println("actual args        = ", c.args.Args)
		fmt.Println("")
	}
}

// FindBach finds and executes Bach
func FindBach(context Context, args *ParsedArgs) *BachCommand {
	pwd := context.GetWorkingDir()

	rootdir, noRootdir := resolveBachRootDir(context, pwd)
	config := ReadConfig(context, rootdir)
	quiet := args.HasGumFlag("gq")

	if quiet {
		config.setQuiet(quiet)
	}

	executable, noExecutable := findBachExecutable(context, config, rootdir)

	if noExecutable != nil {
		warnNoBach(context, config)

		if context.IsExplicit() {
			context.Exit(-1)
		}
		return nil
	}

	if noRootdir != nil {
		if context.IsExplicit() {
			fmt.Println("No Bach project found")
			fmt.Println()
			context.Exit(-1)
		}
		return nil
	}

	p, _ := filepath.Abs(pwd)
	r, _ := filepath.Abs(rootdir)
	if p != r {
		if context.IsExplicit() {
			fmt.Println("Bach must be invoked from " + rootdir)
			fmt.Println()
			context.Exit(-1)
		}
		return nil
	}

	return &BachCommand{
		context:    context,
		config:     config,
		rootdir:    rootdir,
		executable: executable,
		args:       args}
}

func resolveBachRootDir(context Context, dir string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root")
	}

	path := filepath.Join(dir, ".bach")
	if context.FileExists(path) {
		return filepath.Abs(dir)
	}

	return resolveBachRootDir(context, parentdir)
}

func warnNoBach(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Println("No java/jshell found in path. Please install Java 16+")
	}
}

func findBachExecutable(context Context, config *Config, dir string) (string, error) {
	java, noJava := findExecutable(context, dir, "java")
	jshell, noJshell := findExecutable(context, dir, "jshell")

	if noJava != nil {
		if noJshell != nil {
			return "", errors.New("No java nor jshell found")
		}
		return jshell + " https://github.com/sormuras/bach/releases/download/" + config.bach.version + "/build.jsh", nil
	}

	bin := filepath.Join(dir, ".bach", "bin")
	cache := filepath.Join(dir, ".bach", "cache")
	if context.FileExists(bin) {
		return java + " -p " + bin + " -m com.github.sormuras.bach build", nil
	} else if context.FileExists(cache) {
		return java + " -p " + cache + " -m com.github.sormuras.bach build", nil
	} else if noJshell == nil {
		return jshell + " https://github.com/sormuras/bach/releases/download/" + config.bach.version + "/build.jsh", nil
	} else {
		return "", errors.New("jshell not found")
	}
}

func findExecutable(context Context, dir string, cmd string) (string, error) {
	exec := resolveExec(context, cmd)
	paths := context.GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], exec)
		if context.FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(cmd + " not found")
}

func resolveExec(context Context, cmd string) string {
	if context.IsWindows() {
		return cmd + ".bat"
	}
	return cmd
}
