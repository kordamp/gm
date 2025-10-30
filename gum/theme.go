// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2020-2025 Andres Almiray.
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
	"strings"

	"github.com/gookit/color"
)

// DarkTheme dark colored theme
var DarkTheme *ColoredTheme = &ColoredTheme{
	name:    "dark",
	symbol:  color.S256(255, 0),
	section: color.S256(28, 0),
	key:     color.S256(160, 0),
	boolean: color.S256(99, 0),
	literal: color.S256(33, 0)}

// LightTheme light colored theme
var LightTheme *ColoredTheme = &ColoredTheme{
	name:    "light",
	symbol:  color.S256(0, 255),
	section: color.S256(28, 255),
	key:     color.S256(160, 255),
	boolean: color.S256(99, 255),
	literal: color.S256(33, 255)}

// NoneTheme non colored theme
var NoneTheme *noneTheme = &noneTheme{}

// ColoredTheme provides a color theme
type ColoredTheme struct {
	name    string
	symbol  *color.Style256
	section *color.Style256
	key     *color.Style256
	boolean *color.Style256
	literal *color.Style256
}

// PrintSection prints a section header such as [section]
func (t *ColoredTheme) PrintSection(section string) {
	t.symbol.Print("[")
	t.section.Print(section)
	t.symbol.Println("]")
}

// PrintKeyValueBoolean prints a key/value pair as key = value
func (t *ColoredTheme) PrintKeyValueBoolean(key string, value bool) {
	t.key.Print(key)
	t.symbol.Print(" = ")
	t.boolean.Println(value)
}

// PrintKeyValueLiteral prints a key/value pair as key = "value"
func (t *ColoredTheme) PrintKeyValueLiteral(key string, value string) {
	t.key.Print(key)
	t.symbol.Print(" = \"")
	t.literal.Print(value)
	t.literal.Println("\"")
}

// PrintKeyValueArrayS prints a key/value pair as key = ["v1", "v2"]
func (t *ColoredTheme) PrintKeyValueArrayS(key string, value []string) {
	t.key.Print(key)
	t.symbol.Print(" = [")

	for i, w := range value {
		if i != 0 {
			t.symbol.Print(", ")
		}
		t.literal.Print("\"")
		t.literal.Print(w)
		t.literal.Print("\"")
	}

	t.symbol.Println("]")
}

// PrintKeyValueArrayI prints a key/value pair as key = [i1, i2]
func (t *ColoredTheme) PrintKeyValueArrayI(key string, value [2]uint8) {
	t.key.Print(key)
	t.symbol.Print(" = [")
	t.symbol.Print(value[0])
	t.symbol.Print(", ")
	t.symbol.Print(value[1])
	t.symbol.Println("]")
}

// PrintMap prints a map with each entry as key = "value"
func (t *ColoredTheme) PrintMap(value map[string]string) {
	for k, v := range value {
		if strings.Contains(k, ":") {
			t.key.Print("\"" + k + "\"")
		} else {
			t.key.Print(k)
		}
		t.symbol.Print(" = ")
		t.literal.Print("\"")
		t.literal.Print(v)
		t.literal.Println("\"")
	}
}

type noneTheme struct {
}

// PrintSection prints a section header such as [section]
func (t *noneTheme) PrintSection(section string) {
	fmt.Println("[" + section + "]")
}

// PrintKeyValueBoolean prints a key/value pair as key = value
func (t *noneTheme) PrintKeyValueBoolean(key string, value bool) {
	fmt.Println(key+" =", value)
}

// PrintKeyValueLiteral prints a key/value pair as key = "value"
func (t *noneTheme) PrintKeyValueLiteral(key string, value string) {
	fmt.Print(key)
	fmt.Print(" = \"")
	fmt.Print(value)
	fmt.Println("\"")
}

// PrintKeyValueArrayS prints a key/value pair as key = ["v1", "v2"]
func (t *noneTheme) PrintKeyValueArrayS(key string, value []string) {
	fmt.Print(key)
	fmt.Print(" = ")
	fmt.Print("[")

	for i, w := range value {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Print("\"")
		fmt.Print(w)
		fmt.Print("\"")
	}

	fmt.Println("]")
}

// PrintKeyValueArrayI prints a key/value pair as key = [i1, i2]
func (t *noneTheme) PrintKeyValueArrayI(key string, value [2]uint8) {
	fmt.Print(key)
	fmt.Print(" = [")
	fmt.Print(value[0])
	fmt.Print(", ")
	fmt.Print(value[1])
	fmt.Println("]")
}

// PrintMap prints a map with each entry as key = "value"
func (t *noneTheme) PrintMap(value map[string]string) {
	for k, v := range value {
		if strings.Contains(k, ":") {
			fmt.Print("\"" + k + "\"")
		} else {
			fmt.Print(k)
		}
		fmt.Print(" = ")
		fmt.Print("\"")
		fmt.Print(v)
		fmt.Println("\"")
	}
}
