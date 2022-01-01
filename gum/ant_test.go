// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2020-2022 Andres Almiray.
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

func TestAntSingle(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "ant", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "ant", "single"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindAnt(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "ant")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.xml")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestAntParent(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "ant", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "ant", "parent", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindAnt(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "ant")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.xml")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestAntWithExplicitBuildFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "ant", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "ant", "parent-with-explicit", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "-f", filepath.Join(pwd, "explicit.xml")})
	cmd := FindAnt(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "ant")},
		{"BuildFile", cmd.buildFile, ""},
		{"ExplicitBuildFile", cmd.explicitBuildFile, filepath.Join(pwd, "explicit.xml")},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureAnt()
	if cmd.args.Args[0] != "-f" || cmd.args.Args[1] != filepath.Join(pwd, "explicit.xml") {
		t.Errorf("args: invalid build file")
	}
}

func TestAntWithoutExecutables(t *testing.T) {
	// given:
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "ant", "single"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindAnt(context, &args)

	// then:
	if cmd != nil {
		t.Error("Expected a nil command but got something")
	}
}
