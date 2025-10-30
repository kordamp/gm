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
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// given:
	home, _ := filepath.Abs(filepath.Join("..", "tests", "home"))
	root, _ := filepath.Abs(filepath.Join("..", "tests", "toml"))

	context := testContext{
		explicit:   true,
		windows:    false,
		workingDir: root,
		homeDir:    home,
		paths:      []string{home, root}}

	// when:
	config := ReadConfig(context, root)

	// then:
	if config == nil {
		t.Error("Expected config but got nil")
	}

	var checks = []struct {
		title            string
		actual, expected bool
	}{
		{"quiet", config.general.quiet, false},
		{"debug", config.general.debug, true},
		{"gradle.replace", config.gradle.replace, true},
		{"gradle.defaults", config.gradle.defaults, true},
		{"maven.replace", config.maven.replace, true},
		{"maven.defaults", config.maven.defaults, true},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %t, want %t", check.title, check.actual, check.expected)
		}
	}

	if config.theme.name != "dark" {
		t.Errorf("theme.name: got %s, want dark", config.theme.name)
	}

	if config.gradle.mappings["compile"] != "compileJava" {
		t.Errorf("gradle.mappings.compileJava: got %s, want %s", config.gradle.mappings["compile"], "compileJava")
	}

	if config.maven.mappings["compileJava"] != "compile" {
		t.Errorf("maven.mappings.compile: got %s, want %s", config.gradle.mappings["compileJava"], "compile")
	}
}
