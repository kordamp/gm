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

func TestGradleTaskSubstitutionAppendFlag(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "verify", "-S"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureGradle()
	for _, arg := range cmd.args.Args {
		if arg == filepath.Join(pwd, "build.gradle") {
			t.Errorf("args: explicit build file found in args")
		}
	}
	if len(cmd.args.Args) != 3 && cmd.args.Args[len(cmd.args.Args)-2] != "build" {
		t.Errorf("args: got verify, want build")
	}
}

func TestGradleTaskSubstitutionPrependFlag(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "-S", "verify"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureGradle()
	for _, arg := range cmd.args.Args {
		if arg == filepath.Join(pwd, "build.gradle") {
			t.Errorf("args: explicit build file found in args")
		}
	}
	if len(cmd.args.Args) != 3 && cmd.args.Args[len(cmd.args.Args)-1] != "build" {
		t.Errorf("args: got verify, want build")
	}
}
func TestGradleSingleWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "verify"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureGradle()
	for _, arg := range cmd.args.Args {
		if arg == filepath.Join(pwd, "build.gradle") {
			t.Errorf("args: explicit build file found in args")
		}
	}
	if len(cmd.args.Args) != 2 && cmd.args.Args[len(cmd.args.Args)-1] != "build" {
		t.Errorf("args: got verify, want build")
	}
}

func TestGradleSingleWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "gradle")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleParentWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleParentWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-without-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "gradle")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithExplicitBuildFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-explicit", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "-b", filepath.Join(pwd, "explicit.gradle")})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, ""},
		{"BuildFile", cmd.buildFile, ""},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, filepath.Join(pwd, "explicit.gradle")},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureGradle()
	if cmd.args.Args[0] != "-b" || cmd.args.Args[1] != filepath.Join(pwd, "explicit.gradle") {
		t.Errorf("args: invalid build file")
	}
}

func TestGradleWithExplicitSettingsFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "-c", filepath.Join(pwd, "..", "settings.gradle")})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithExplicitProjectDir(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "-p", filepath.Join(pwd, "..")})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, ""},
		{"BuildFile", cmd.buildFile, ""},
		{"SettingsFile", cmd.settingsFile, ""},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, filepath.Join(pwd, "..")},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithNearestBuildFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-conventional-child", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "-gn"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "child.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithoutExecutables(t *testing.T) {
	// given:
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd != nil {
		t.Error("Expected a nil command but got something")
	}
}

func TestGradleReplaceWithExactMatch(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "verify"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	cmd.doConfigureGradle()
	if len(cmd.args.Args) != 2 && cmd.args.Args[len(cmd.args.Args)-1] != "build" {
		t.Errorf("args: got verify, want build")
	}
}

func TestGradleReplaceWithSubMatch(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", ":subproject:verify"})
	cmd := FindGradle(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	cmd.doConfigureGradle()
	if len(cmd.args.Args) != 2 && cmd.args.Args[len(cmd.args.Args)-1] != ":subproject:build" {
		t.Errorf("args: got :subproject:verify, want b:subproject:build")
	}
}
