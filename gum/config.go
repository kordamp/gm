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
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/grignaak/tribool"
	"github.com/pelletier/go-toml"
)

// Config defines configuration settings for Gum
type Config struct {
	general general
	gradle  gradle
	maven   maven
	jbang   jbang
}

type general struct {
	quiet     bool
	debug     bool
	discovery []string

	q tribool.Tribool
	d tribool.Tribool
}

type gradle struct {
	replace  bool
	defaults bool
	mappings map[string]string

	r tribool.Tribool
	d tribool.Tribool
}

type maven struct {
	replace  bool
	defaults bool
	mappings map[string]string

	r tribool.Tribool
	d tribool.Tribool
}

type jbang struct {
	discovery []string
}

func (c *Config) print() {
	fmt.Println("[general]")
	fmt.Println("quiet =", c.general.quiet)
	fmt.Println("debug =", c.general.debug)
	fmt.Println("discovery =", formatSlice(c.general.discovery))
	fmt.Println("")
	fmt.Println("[gradle]")
	fmt.Println("replace =", c.gradle.replace)
	fmt.Println("defaults =", c.gradle.defaults)
	if len(c.gradle.mappings) > 0 {
		fmt.Println("[gradle.mappings]")
		printMappings(c.gradle.mappings)
	}
	fmt.Println("")
	fmt.Println("[maven]")
	fmt.Println("replace =", c.maven.replace)
	fmt.Println("defaults =", c.maven.defaults)
	if len(c.maven.mappings) > 0 {
		fmt.Println("[maven.mappings]")
		printMappings(c.maven.mappings)
	}
	fmt.Println("")
	fmt.Println("[jbang]")
	fmt.Println("discovery =", formatSlice(c.jbang.discovery))
}

func formatSlice(s []string) string {
	var buffer bytes.Buffer
	buffer.WriteString("[")

	for i, w := range s {
		if i != 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString("\"")
		buffer.WriteString(w)
		buffer.WriteString("\"")
	}

	buffer.WriteString("]")
	return buffer.String()
}

func printMappings(mappings map[string]string) {
	for k, v := range mappings {
		if strings.Contains(k, ":") {
			fmt.Print("\"" + k + "\"")
		} else {
			fmt.Print(k)
		}
		fmt.Println(" = \"" + v + "\"")
	}
}

func newConfig() *Config {
	return &Config{
		general: general{
			q:         tribool.Maybe,
			d:         tribool.Maybe,
			discovery: make([]string, 0)},
		gradle: gradle{
			r:        tribool.Maybe,
			d:        tribool.Maybe,
			mappings: make(map[string]string)},
		maven: maven{
			r:        tribool.Maybe,
			d:        tribool.Maybe,
			mappings: make(map[string]string)},
		jbang: jbang{
			discovery: make([]string, 0)}}
}

func (c *Config) setQuiet(b bool) {
	c.general.quiet = b
}

func (c *Config) setDebug(b bool) {
	c.general.debug = b
}

func (g *gradle) setReplace(b bool) {
	g.replace = b
}

func (m *maven) setReplace(b bool) {
	m.replace = b
}

func (c *Config) merge(other *Config) {
	if other == nil {
		c.general.merge(nil)
		c.gradle.merge(nil)
		c.maven.merge(nil)
		c.jbang.merge(nil)
	} else {
		c.general.merge(&other.general)
		c.gradle.merge(&other.gradle)
		c.maven.merge(&other.maven)
		c.jbang.merge(&other.jbang)
	}
}

func (g *general) merge(other *general) {
	if g.q != tribool.Maybe || other == nil {
		g.quiet = g.q.WithMaybeAsFalse()
	} else if other != nil {
		g.quiet = other.q.WithMaybeAsFalse()
	}

	if g.d != tribool.Maybe || other == nil {
		g.debug = g.d.WithMaybeAsFalse()
	} else {
		g.debug = other.d.WithMaybeAsFalse()
	}

	if len(g.discovery) != 3 && other != nil {
		g.discovery = other.discovery
	}
}

func (g *gradle) merge(other *gradle) {
	if g.r != tribool.Maybe || other == nil {
		g.replace = g.r.WithMaybeAsTrue()
	} else {
		g.replace = other.r.WithMaybeAsTrue()
	}

	if g.d != tribool.Maybe || other == nil {
		g.defaults = g.d.WithMaybeAsTrue()
	} else {
		g.defaults = other.d.WithMaybeAsTrue()
	}

	mp := make(map[string]string)
	if g.defaults {
		mp = map[string]string{
			"compile":         "classes",
			"package":         "assemble",
			"verify":          "build",
			"install":         "publishToMavenLocal",
			"exec:java":       "run",
			"dependency:tree": "dependencies"}
	}
	if other != nil {
		for k, v := range other.mappings {
			mp[k] = v
		}
	}
	for k, v := range g.mappings {
		mp[k] = v
	}
	g.mappings = mp
}

func (m *maven) merge(other *maven) {
	if m.r != tribool.Maybe || other == nil {
		m.replace = m.r.WithMaybeAsTrue()
	} else {
		m.replace = other.r.WithMaybeAsTrue()
	}

	if m.d != tribool.Maybe || other == nil {
		m.defaults = m.d.WithMaybeAsTrue()
	} else {
		m.defaults = other.d.WithMaybeAsTrue()
	}

	mp := make(map[string]string)
	if m.defaults {
		mp = map[string]string{
			"classes":             "compile",
			"jar":                 "package",
			"assemble":            "package",
			"build":               "verify",
			"publishToMavenLocal": "install",
			"puTML":               "install",
			"check":               "verify",
			"run":                 "exec:java",
			"dependencies":        "dependency:tree"}
	}
	if other != nil {
		for k, v := range other.mappings {
			mp[k] = v
		}
	}
	for k, v := range m.mappings {
		mp[k] = v
	}
	m.mappings = mp
}

func (j *jbang) merge(other *jbang) {
	if len(j.discovery) == 0 && other != nil && len(other.discovery) == 3 {
		j.discovery = make([]string, 3)
		copy(j.discovery, other.discovery)
	}
}

// ReadUserConfig reads user config
func ReadUserConfig(context Context) *Config {
	homedir := context.GetHomeDir()
	tomlfile := filepath.Join(homedir, ".gm.toml")
	if context.IsWindows() {
		tomlfile = filepath.Join(homedir, "Gum", "gm.toml")
	}

	return ReadConfigFile(context, tomlfile)
}

// ReadConfig reads and merges project & user config
func ReadConfig(context Context, rootdir string) *Config {
	uconfig := ReadUserConfig(context)
	pconfig := ReadConfigFile(context, filepath.Join(rootdir, ".gm.toml"))

	pconfig.merge(uconfig)

	return pconfig
}

// ReadConfigFile reads the given TOML config file
func ReadConfigFile(context Context, path string) *Config {
	config := newConfig()

	if !context.FileExists(path) {
		return config
	}

	doc, err := ioutil.ReadFile(path)
	if err == nil {
		toml.Unmarshal(doc, &config)
	} else {
		fmt.Println(err)
	}

	t, err := toml.LoadBytes(doc)
	if err != nil {
		return config
	}

	resolveSectionGeneral(t, config)
	resolveSectionGradle(t, config)
	resolveSectionMaven(t, config)
	resolveSectionJbang(t, config)

	return config
}

func resolveSectionGeneral(t *toml.Tree, config *Config) {
	tt := t.Get("general")
	if tt != nil {
		table := tt.(*toml.Tree)
		v := table.Get("quiet")
		if v != nil {
			config.general.q = tribool.FromBool(v.(bool))
		}
		v = table.Get("debug")
		if v != nil {
			config.general.d = tribool.FromBool(v.(bool))
		}
		v = table.Get("discovery")
		if v != nil {
			data := v.([]interface{})
			config.general.discovery = make([]string, len(data))
			for i, e := range data {
				config.general.discovery[i] = e.(string)
			}
		}
	}
}

func resolveSectionGradle(t *toml.Tree, config *Config) {
	tt := t.Get("gradle")
	if tt != nil {
		table := tt.(*toml.Tree)
		v := table.Get("replace")
		if v != nil {
			config.gradle.r = tribool.FromBool(v.(bool))
		}
		v = table.Get("defaults")
		if v != nil {
			config.gradle.d = tribool.FromBool(v.(bool))
		}
		v = table.Get("mappings")
		if v != nil {
			m := v.(*toml.Tree)
			for i := range m.Keys() {
				key := m.Keys()[i]
				config.gradle.mappings[key] = m.Get(key).(string)
			}
		}
	}
}
func resolveSectionMaven(t *toml.Tree, config *Config) {
	tt := t.Get("maven")
	if tt != nil {
		table := tt.(*toml.Tree)
		v := table.Get("replace")
		if v != nil {
			config.maven.r = tribool.FromBool(v.(bool))
		}
		v = table.Get("defaults")
		if v != nil {
			config.maven.d = tribool.FromBool(v.(bool))
		}
		v = table.Get("mappings")
		if v != nil {
			m := v.(*toml.Tree)
			for i := range m.Keys() {
				key := m.Keys()[i]
				config.maven.mappings[key] = m.Get(key).(string)
			}
		}
	}
}

func resolveSectionJbang(t *toml.Tree, config *Config) {
	tt := t.Get("jbang")
	if tt != nil {
		table := tt.(*toml.Tree)
		v := table.Get("discovery")
		if v != nil {
			data := v.([]interface{})
			config.jbang.discovery = make([]string, len(data))
			for i, e := range data {
				config.jbang.discovery[i] = e.(string)
			}
		}
	}
}
