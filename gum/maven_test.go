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
	"path/filepath"
	"testing"
)

func TestMavenSingleWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "single-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindMaven(context, []string{"build"})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "mvnw")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "pom.xml")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "pom.xml")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureMaven()
	if len(cmd.args) != 2 && cmd.args[len(cmd.args)-1] != "verify" {
		t.Errorf("args: got build, want verify")
	}
}

func TestMavenSingleWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "single-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindMaven(context, []string{})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "mvn")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "pom.xml")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "pom.xml")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestMavenParentWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "parent-with-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindMaven(context, []string{})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "mvnw")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "pom.xml")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "pom.xml")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestMavenParentWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "parent-without-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindMaven(context, []string{})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "mvn")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "pom.xml")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "pom.xml")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestMavenWithExplicitBuildFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "parent-with-explicit", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindMaven(context, []string{"-f", filepath.Join(pwd, "explicit.xml")})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "mvn")},
		{"RootBuildFile", cmd.rootBuildFile, ""},
		{"BuildFile", cmd.buildFile, ""},
		{"ExplicitBuildFile", cmd.explicitBuildFile, filepath.Join(pwd, "explicit.xml")},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestMavenWithNearestBuildFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "parent-without-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindMaven(context, []string{"-gn"})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "mvn")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "pom.xml")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "pom.xml")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestMavenWithoutExecutables(t *testing.T) {
	// given:
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "maven", "single-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{}}

	// when:
	cmd := FindMaven(context, []string{})

	// then:
	if cmd != nil {
		t.Error("Expected a nil command but got something")
	}
}
