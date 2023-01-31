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
	"strings"
)

// ParsedArgs captures input args separated by responsibility
type ParsedArgs struct {
	Gum  map[string]struct{}
	Tool []string
	Args []string
}

// HasGumFlag finds if a given Gum flag is specified in the parsed args
func (a *ParsedArgs) HasGumFlag(flag string) bool {
	_, ok := a.Gum[flag]
	return ok
}

var gumFlags = []string{"ga", "gb", "gc", "gd", "gg", "gh", "gj", "gm", "gn", "gq", "gr", "gv"}

// ParseArgs parses input args and separates them between Gum, Tool, and Args
func ParseArgs(args []string) ParsedArgs {
	flags := ParsedArgs{
		Gum:  make(map[string]struct{}, 0),
		Tool: make([]string, 0),
		Args: make([]string, 0)}

	if len(args) == 0 {
		return flags
	}

	// 0 = gum
	// 1 = tool
	// 2 = args
	mode := 0

	for i := 0; i < len(args); i++ {
		s := strings.TrimSpace(args[i])

		switch mode {
		case 0:
			if s[0] == '-' && isGumFlag(s) {
				flags.Gum[s[1:]] = struct{}{}
			} else {
				mode = 1
				i = i - 1
			}
		case 1:
			if s[0] == '-' {
				flags.Tool = append(flags.Tool, s)
				if strings.Index(s, "=") == -1 {
					// grab the next value if available
					j := i + 1
					if j < len(args) {
						arg := args[j]
						// add it only if it's not another flag
						if arg[0] != '-' {
							flags.Tool = append(flags.Tool, arg)
							i = j
						}
					}
				}
			} else {
				mode = 2
				i = i - 1
			}
		case 2:
			flags.Args = append(flags.Args, s)
		}
	}

	return flags
}

func isGumFlag(flag string) bool {
	for _, f := range gumFlags {
		if flag[1:] == f {
			return true
		}
	}
	return false
}

func findFlagValue(flag string, args []string) (bool, string, []string) {
	if len(args) == 0 {
		return false, "", args
	}

	for i, s := range args {
		if flag == s {
			// next argument should contain the value we want
			if i+1 < len(args) {
				return true, args[i+1], shrinkSlice(args, i, 2)
			}
			return false, "", args
		}
		// check if format is flag=value
		parts := strings.Split(s, "=")
		if len(parts) == 2 && parts[0] == flag {
			return true, parts[1], shrinkSlice(args, i, 1)
		}
	}

	return false, "", args
}

func shrinkSlice(s []string, index int, length int) []string {
	shrunk := make([]string, len(s)-length)

	j := 0
	for i, e := range s {
		if i < index || i >= index+length {
			shrunk[j] = e
			j = j + 1
		}
	}

	return shrunk
}

func replaceArgs(args []string, replacements map[string]string, allowsSubMatch bool) []string {
	nargs := make([]string, 0)

	for _, key := range args {
		exactMatch := replacements[key]

		subMatch := ""

		if allowsSubMatch {
			semicolon := strings.LastIndex(key, ":")
			if semicolon > -1 {
				prefix := key[0:(semicolon + 1)]
				suffix := key[(semicolon + 1):]
				match := replacements[suffix]

				if len(match) > 0 {
					subMatch = prefix + match
				}
			}
		}

		if len(exactMatch) > 0 {
			nargs = append(nargs, exactMatch)
		} else if allowsSubMatch && len(subMatch) > 0 {
			nargs = append(nargs, subMatch)
		} else {
			nargs = append(nargs, key)
		}
	}

	return nargs
}
