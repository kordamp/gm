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

// Command defines an executable command (gradle/maven)
type Command interface {
	// Execute executes the given command
	Execute()
}

// Context provides an abstraction over the OS and Environment as required by Gum
type Context interface {
	// IsExplicit whether a given tool was specified
	IsExplicit() bool

	// IsWindows checks if the current OS is Windows
	IsWindows() bool

	// GetWorkingDir returns the current working dir
	GetWorkingDir() string

	// GetHomeDir gets the home directory from environment
	GetHomeDir() string

	// GetPaths gets the paths in $PATH
	GetPaths() []string

	// FileExists checks if a file exists
	FileExists(name string) bool

	// Exit causes the current program to exit with the given status code.
	Exit(code int)
}
