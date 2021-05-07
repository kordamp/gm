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

	// CheckIsExecutable checks if the given file has executable bits
	CheckIsExecutable(file string)

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

// Theme defines a console theme for printing messages
type Theme interface {
	// PrintSection prints a section header such as [section]
	PrintSection(section string)

	// PrintKeyValueBoolean prints a key/value pair as key = value
	PrintKeyValueBoolean(key string, value bool)

	// PrintKeyValueLiteral prints a key/value pair as key = "value"
	PrintKeyValueLiteral(key string, value string)

	// PrintKeyValueArrayS prints a key/value pair as key = ["v1", "v2"]
	PrintKeyValueArrayS(key string, value []string)

	// PrintKeyValueArrayI prints a key/value pair as key = [i1, i2]
	PrintKeyValueArrayI(key string, value [2]uint8)

	// PrintMap prints a map with each entry as key = "value"
	PrintMap(value map[string]string)
}
