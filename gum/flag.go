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
	"strings"
)

// GrabFlag finds a boolean flag in the given args, removing it from the slice if found
func GrabFlag(f string, args []string) (bool, []string) {
	if len(args) == 0 {
		// no args to be checked
		return false, args
	}

	for i := range args {
		s := args[i]
		if s == f {
			newArgs := make([]string, len(args)-1)
			for j := 0; j < i; j++ {
				newArgs[j] = args[j]
			}
			for j := i; j < len(newArgs); j++ {
				newArgs[j] = args[j+1]
			}
			return true, newArgs
		}
	}

	return false, args
}

func findFlag(f string, args []string) bool {
	if len(args) == 0 {
		// no args to be checked
		return false
	}

	for i := range args {
		s := args[i]
		if s == f {
			return true
		}
	}

	return false
}

func findFlagValue(flag string, args []string) (bool, string) {
	if len(args) == 0 {
		return false, ""
	}

	for i := range args {
		s := args[i]
		if flag == s {
			// next argument should contain the value we want
			if i+1 < len(args) {
				return true, args[i+1]
			}
			return false, ""
		}
		// check if format is flag=value
		parts := strings.Split(s, "=")
		if len(parts) == 2 && parts[0] == flag {
			return true, parts[1]
		}
	}

	return false, ""
}

func replaceArgs(args []string, replacements map[string]string) []string {
	nargs := make([]string, 0)

	for i := range args {
		key := args[i]
		val := replacements[key]

		if len(val) > 0 {
			nargs = append(nargs, val)
		} else {
			nargs = append(nargs, key)
		}
	}

	return nargs
}
