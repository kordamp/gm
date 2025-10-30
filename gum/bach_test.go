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
	"strings"
	"testing"
)

func TestBachProjectWithoutCache(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "bach", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "bach", "project-without-cache"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindBach(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "jshell")},
	}

	for _, check := range checks {
		if !strings.HasPrefix(check.actual, check.expected) {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestBachProjectWithCache(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "bach", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "bach", "project-with-cache"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindBach(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "java")},
	}

	for _, check := range checks {
		if !strings.HasPrefix(check.actual, check.expected) {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestBachWithoutExecutables(t *testing.T) {
	// given:
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "bach", "project-without-cache"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindBach(context, &args)

	// then:
	if cmd != nil {
		t.Error("Expected a nil command but got something")
	}
}
