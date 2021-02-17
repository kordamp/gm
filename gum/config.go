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

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/grignaak/tribool"
	"github.com/pelletier/go-toml"
)

// Config defines configuration settings for Gum
type Config struct {
	theme   theme
	general general
	gradle  gradle
	maven   maven
	jbang   jbang
	bach    bach
}

type theme struct {
	t Theme

	name    string
	symbol  [2]uint8
	section [2]uint8
	key     [2]uint8
	boolean [2]uint8
	literal [2]uint8
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

type bach struct {
	version string
}

func (c *Config) print() {
	c.theme.t.PrintSection("theme")
	c.theme.t.PrintKeyValueLiteral("name", c.theme.name)
	if isInstanceOf(c.theme.t, (*ColoredTheme)(nil)) {
		c.theme.t.PrintKeyValueArrayI("symbol", c.theme.symbol)
		c.theme.t.PrintKeyValueArrayI("section", c.theme.section)
		c.theme.t.PrintKeyValueArrayI("key", c.theme.key)
		c.theme.t.PrintKeyValueArrayI("boolean", c.theme.boolean)
		c.theme.t.PrintKeyValueArrayI("literal", c.theme.literal)
	}
	c.theme.t.PrintSection("general")
	c.theme.t.PrintKeyValueBoolean("quiet", c.general.quiet)
	c.theme.t.PrintKeyValueBoolean("debug", c.general.debug)
	c.theme.t.PrintKeyValueArrayS("discovery", c.general.discovery)
	c.theme.t.PrintSection("gradle")
	c.theme.t.PrintKeyValueBoolean("replace", c.gradle.replace)
	c.theme.t.PrintKeyValueBoolean("defaults", c.gradle.defaults)
	if len(c.gradle.mappings) > 0 {
		c.theme.t.PrintSection("gradle.mappings")
		c.theme.t.PrintMap(c.gradle.mappings)
	}
	c.theme.t.PrintSection("maven")
	c.theme.t.PrintKeyValueBoolean("replace", c.maven.replace)
	c.theme.t.PrintKeyValueBoolean("defaults", c.maven.defaults)
	if len(c.maven.mappings) > 0 {
		c.theme.t.PrintSection("maven.mappings")
		c.theme.t.PrintMap(c.maven.mappings)
	}
	c.theme.t.PrintSection("jbang")
	c.theme.t.PrintKeyValueArrayS("discovery", c.jbang.discovery)
	c.theme.t.PrintSection("bach")
	c.theme.t.PrintKeyValueLiteral("version", c.bach.version)
}

func newConfig() *Config {
	return &Config{
		theme: theme{
			t:       DarkTheme,
			name:    "dark",
			symbol:  [2]uint8{255, 0},
			section: [2]uint8{28, 0},
			key:     [2]uint8{160, 0},
			boolean: [2]uint8{99, 0},
			literal: [2]uint8{33, 0}},
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
			discovery: make([]string, 0)},
		bach: bach{
			version: ""}}
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
		c.bach.merge(nil)
	} else {
		c.general.merge(&other.general)
		c.gradle.merge(&other.gradle)
		c.maven.merge(&other.maven)
		c.jbang.merge(&other.jbang)
		c.bach.merge(&other.bach)
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

func (b *bach) merge(other *bach) {
	if len(b.version) == 0 && other != nil && len(other.version) > 0 {
		b.version = other.version
	} else {
		b.version = "16.0.2"
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

	resolveSectionTheme(t, config)
	resolveSectionGeneral(t, config)
	resolveSectionGradle(t, config)
	resolveSectionMaven(t, config)
	resolveSectionJbang(t, config)
	resolveSectionBach(t, config)

	return config
}

func resolveSectionTheme(t *toml.Tree, config *Config) {
	tt := t.Get("theme")
	if tt != nil {
		table := tt.(*toml.Tree)
		v := table.Get("name")
		if v != nil {
			config.theme.name = v.(string)
		}
		if config.theme.name == "none" {
			config.theme.t = NoneTheme
			return
		}

		v = table.Get("symbol")
		if v != nil {
			data := v.([]interface{})
			config.theme.symbol = [2]uint8{uint8(data[0].(int64)), uint8(data[1].(int64))}
		} else {
			config.theme.symbol = [2]uint8{125, 0}
		}
		v = table.Get("section")
		if v != nil {
			data := v.([]interface{})
			config.theme.section = [2]uint8{uint8(data[0].(int64)), uint8(data[1].(int64))}
		} else {
			config.theme.symbol = [2]uint8{47, 0}
		}
		v = table.Get("key")
		if v != nil {
			data := v.([]interface{})
			config.theme.key = [2]uint8{uint8(data[0].(int64)), uint8(data[1].(int64))}
		} else {
			config.theme.symbol = [2]uint8{130, 0}
		}
		v = table.Get("boolean")
		if v != nil {
			data := v.([]interface{})
			config.theme.boolean = [2]uint8{uint8(data[0].(int64)), uint8(data[1].(int64))}
		} else {
			config.theme.symbol = [2]uint8{200, 0}
		}
		v = table.Get("literal")
		if v != nil {
			data := v.([]interface{})
			config.theme.literal = [2]uint8{uint8(data[0].(int64)), uint8(data[1].(int64))}
		} else {
			config.theme.symbol = [2]uint8{23, 0}
		}
	}

	switch strings.ToLower(config.theme.name) {
	case "light":
		config.theme.t = LightTheme
		break
	case "custom":
		config.theme.t = &ColoredTheme{
			name:    "custom",
			symbol:  color.S256(config.theme.symbol[0], config.theme.symbol[1]),
			section: color.S256(config.theme.section[0], config.theme.section[1]),
			key:     color.S256(config.theme.key[0], config.theme.key[1]),
			boolean: color.S256(config.theme.boolean[0], config.theme.boolean[1]),
			literal: color.S256(config.theme.literal[0], config.theme.literal[1])}
		break
	case "dark":
	default:
		config.theme.t = DarkTheme
		break
	}
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

func resolveSectionBach(t *toml.Tree, config *Config) {
	tt := t.Get("bach")
	if tt != nil {
		table := tt.(*toml.Tree)
		v := table.Get("version")
		if v != nil {
			config.bach.version = v.(string)
		}
	}
}
