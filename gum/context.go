// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2020-2022 Andres Almiray.
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
	"fmt"
	"os"
	"runtime"
	"strings"
)

// DefaultContext is the Context used by default
type DefaultContext struct {
	explicit bool
}

// NewDefaultContext creates a new DefaultContext with the given state
func NewDefaultContext(explicit bool) DefaultContext {
	return DefaultContext{explicit: explicit}
}

// IsExplicit whether a given tool was specified
func (c DefaultContext) IsExplicit() bool {
	return c.explicit
}

// IsWindows checks if the current OS is Windows
func (c DefaultContext) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// CheckIsExecutable checks if the given file has executable bits
func (c DefaultContext) CheckIsExecutable(file string) {
	if !c.IsWindows() {
		fileInfo, _ := os.Stat(file)
		if fileInfo.Mode().Perm()&0111 == 0 {
			fmt.Printf("%s is not executable", file)
			fmt.Println()
			c.Exit(-1)
		}
	}
}

// GetWorkingDir returns the current working dir
func (c DefaultContext) GetWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return pwd
}

// GetPaths gets the paths in $PATH
func (c DefaultContext) GetPaths() []string {
	return strings.Split(c.getPathFromEnv(), string(os.PathListSeparator))
}

// Gets the PATH environment variable
func (c DefaultContext) getPathFromEnv() string {
	if c.IsWindows() {
		return os.Getenv("Path")
	}

	return os.Getenv("PATH")
}

// GetHomeDir gets the home directory from environment
func (c DefaultContext) GetHomeDir() string {
	if c.IsWindows() {
		return os.Getenv("APPDATA")
	}

	return os.Getenv("HOME")
}

// FileExists checks if a file exists
func (c DefaultContext) FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Exit causes the current program to exit with the given status code.
func (c DefaultContext) Exit(code int) {
	os.Exit(code)
}

// -----------------------------------------------

type testContext struct {
	quiet      bool
	explicit   bool
	windows    bool
	workingDir string
	homeDir    string
	paths      []string
	exitCode   int
}

func (c testContext) IsQuiet() bool {
	return c.quiet
}

func (c testContext) IsExplicit() bool {
	return c.explicit
}

func (c testContext) CheckIsExecutable(file string) {

}

func (c testContext) IsWindows() bool {
	return c.windows
}

func (c testContext) GetWorkingDir() string {
	return c.workingDir
}

func (c testContext) GetHomeDir() string {
	return c.homeDir
}

func (c testContext) GetPaths() []string {
	return c.paths
}

func (c testContext) FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (c testContext) Exit(code int) {
	c.exitCode = code
}
