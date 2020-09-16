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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/grignaak/tribool"
	"github.com/pelletier/go-toml"
)

// Config defines configuration settings for Gum
type Config struct {
	general general
	gradle  gradle
	maven   maven
}

type general struct {
	quiet bool
	debug bool

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

func newConfig() *Config {
	return &Config{
		general: general{
			q: tribool.Maybe,
			d: tribool.Maybe},
		gradle: gradle{
			r:        tribool.Maybe,
			d:        tribool.Maybe,
			mappings: make(map[string]string)},
		maven: maven{
			r:        tribool.Maybe,
			d:        tribool.Maybe,
			mappings: make(map[string]string)}}
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
	} else {
		c.general.merge(&other.general)
		c.gradle.merge(&other.gradle)
		c.maven.merge(&other.maven)
	}
}

func (g *general) merge(other *general) {
	if g.q != tribool.Maybe {
		g.quiet = g.q.WithMaybeAsFalse()
	} else {
		g.quiet = other.q.WithMaybeAsFalse()
	}

	if g.d != tribool.Maybe {
		g.debug = g.d.WithMaybeAsFalse()
	} else {
		g.debug = other.d.WithMaybeAsFalse()
	}
}

func (g *gradle) merge(other *gradle) {
	if g.r != tribool.Maybe {
		g.replace = g.r.WithMaybeAsFalse()
	} else {
		g.replace = other.r.WithMaybeAsTrue()
	}

	if g.d != tribool.Maybe {
		g.defaults = g.d.WithMaybeAsFalse()
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
	for k, v := range other.mappings {
		mp[k] = v
	}
	for k, v := range g.mappings {
		mp[k] = v
	}
	g.mappings = mp
}

func (m *maven) merge(other *maven) {
	if m.r != tribool.Maybe {
		m.replace = m.r.WithMaybeAsFalse()
	} else {
		m.replace = other.r.WithMaybeAsTrue()
	}

	if m.d != tribool.Maybe {
		m.defaults = m.d.WithMaybeAsFalse()
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
	for k, v := range other.mappings {
		mp[k] = v
	}
	for k, v := range m.mappings {
		mp[k] = v
	}
	m.mappings = mp
}

// ReadConfig reads and merges project & user config
func ReadConfig(context Context, rootdir string) *Config {
	homedir := context.GetHomeDir()
	tomlfile := filepath.Join(homedir, ".gm.toml")
	if context.IsWindows() {
		tomlfile = filepath.Join(homedir, "Gum", "gm.toml")
	}

	uconfig := ReadConfigFile(context, tomlfile)

	tomlfile = filepath.Join(rootdir, ".gm.toml")
	pconfig := ReadConfigFile(context, tomlfile)

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
