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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// JbangCommand defines an executable Jbang command
type JbangCommand struct {
	context            Context
	config             *Config
	executable         string
	args               []string
	sourceFile         string
	explicitSourceFile string
}

// Execute executes the given command
func (c JbangCommand) Execute() {
	c.doConfigureJbang()
	c.doExecuteJbang()
}

func (c *JbangCommand) doConfigureJbang() {
	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using jbang at '"+c.executable+"'")
	debug, oargs := GrabFlag("-gd", c.args)

	c.config.setDebug(debug)

	if len(c.explicitSourceFile) > 0 {
		banner = append(banner, "to run sourceFile '"+c.explicitSourceFile+"':")
	} else if len(c.sourceFile) > 0 {
		args = append(args, c.sourceFile)
		banner = append(banner, "to run sourceFile '"+c.sourceFile+"':")
	}

	for i := range oargs {
		args = append(args, oargs[i])
	}
	c.args = args

	c.debugJbang(oargs, args)

	if !c.config.general.quiet {
		fmt.Println(strings.Join(banner, " "))
	}
}

func (c *JbangCommand) doExecuteJbang() {
	cmd := exec.Command(c.executable, c.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func (c *JbangCommand) debugJbang(oargs []string, args []string) {
	if c.config.general.debug {
		fmt.Println("sourceFile         = ", c.sourceFile)
		fmt.Println("explicitSourceFile = ", c.explicitSourceFile)
		fmt.Println("original args      = ", oargs)
		fmt.Println("actual args        = ", args)
		fmt.Println("")
	}
}

// FindJbang finds and executes jbang
func FindJbang(context Context, args []string) *JbangCommand {
	pwd := context.GetWorkingDir()

	jbangw, noWrapper := findJbangWrapperExec(context, pwd)
	jbang, noJbang := findJbangExec(context)
	explicitSourceFileSet, explicitSourceFile := findExplicitJbangSourceFile(pwd, args)

	sourceFile, noSourceFile := findJbangSourceFile(context, pwd, args)
	rootdir := resolveJbangRootDir(context, explicitSourceFile, sourceFile)
	config := ReadConfig(context, rootdir)
	config.setQuiet(context.IsQuiet())

	var executable string
	if noWrapper == nil {
		executable = jbangw
	} else if noJbang == nil {
		warnNoJbangWrapper(context, config)
		executable = jbang
	} else {
		warnNoJbang(context, config)

		if context.IsExplicit() {
			context.Exit(-1)
		}
		return nil
	}

	if explicitSourceFileSet {
		return &JbangCommand{
			context:            context,
			config:             config,
			executable:         executable,
			args:               args,
			explicitSourceFile: explicitSourceFile}
	}

	if noSourceFile != nil {
		if context.IsExplicit() {
			fmt.Println("No jbang project found")
			fmt.Println()
			context.Exit(-1)
		}
		return nil
	}

	return &JbangCommand{
		context:    context,
		config:     config,
		executable: executable,
		args:       args,
		sourceFile: sourceFile}
}

func resolveJbangRootDir(context Context,
	explicitSourceFile string,
	sourceFile string) string {

	if context.FileExists(explicitSourceFile) {
		return filepath.Dir(explicitSourceFile)
	}
	return filepath.Dir(sourceFile)
}

func warnNoJbangWrapper(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s set up for this project. ", resolveJbangWrapperExec(context))
		fmt.Println("Please consider setting one up.")
		fmt.Println("(https://github.com/jbangdev)")
		fmt.Println()
	}
}

func warnNoJbang(context Context, config *Config) {
	if !config.general.quiet && context.IsExplicit() {
		fmt.Printf("No %s found in path. Please install jbang.", resolveJbangExec(context))
		fmt.Println("(https://github.com/jbangdev)")
		fmt.Println()
	}
}

// Finds the jbang executable
func findJbangExec(context Context) (string, error) {
	jbang := resolveJbangExec(context)
	paths := context.GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], jbang)
		if context.FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(jbang + " not found")
}

// Finds the Jbang wrapper (if it exists)
func findJbangWrapperExec(context Context, dir string) (string, error) {
	wrapper := resolveJbangWrapperExec(context)
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New(wrapper + " not found")
	}

	path := filepath.Join(dir, wrapper)
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return "", errors.New(wrapper + " not found")
}

func isLaunchableSource(source string) bool {
	if isLaunchableURL(source) {
		return true
	}
	return isLaunchableSourceFile(source)
}

func isLaunchableURL(source string) bool {
	if strings.HasPrefix(source, "http") {
		return true
	}
	match, _ := regexp.MatchString(".+@.+", source)
	return match
}

func isLaunchableSourceFile(source string) bool {
	if strings.HasSuffix(source, ".java") {
		return true

	} else if strings.HasSuffix(source, ".jsh") {
		return true
	}
	return false
}

func findExplicitJbangSourceFile(pwd string, args []string) (bool, string) {
	// grab the first non flag arg
	file := ""
	for i := range args {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			file = arg
			break
		}
	}

	if len(file) > 0 {
		if isLaunchableURL(file) {
			return true, file
		} else if isLaunchableSourceFile(file) {
			if filepath.IsAbs(file) {
				return true, file
			}
			return true, filepath.Join(pwd, file)
		}
	}

	return false, ""
}

// Finds the nearest source file
func findJbangSourceFile(context Context, dir string, args []string) (string, error) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return "", err
	}

	for i := range files {
		file := files[i]
		if isLaunchableSourceFile(file.Name()) {
			f, err := filepath.Abs(file.Name())
			return filepath.Join(dir, filepath.Base(f)), err
		}
	}

	return "", errors.New("Did not find a launchable source")
}

// Resolves the jbangw executable (OS dependent)
func resolveJbangWrapperExec(context Context) string {
	if context.IsWindows() {
		return "jbang.cmd"
	}
	return "jbang"
}

// Resolves the jbang executable (OS dependent)
func resolveJbangExec(context Context) string {
	if context.IsWindows() {
		return "jbang.cmd"
	}
	return "jbang"
}
