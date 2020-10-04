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

func TestJbangJavaWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "java-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "foo"})
	cmd := FindJbang(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "jbang")},
		{"SourceFile", cmd.sourceFile, filepath.Join(pwd, "hello.java")},
		{"ExplicitSourceFile", cmd.explicitSourceFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureJbang()
	if len(cmd.args.Args) != 2 {
		t.Errorf("invalid arg count")
	}
}

func TestJbangJavaWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "java-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "foo"})
	cmd := FindJbang(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "jbang")},
		{"SourceFile", cmd.sourceFile, filepath.Join(pwd, "hello.java")},
		{"ExplicitSourceFile", cmd.explicitSourceFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestJbangJavaWithExplicitFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "java-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "zzz.java", "foo"})
	cmd := FindJbang(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "jbang")},
		{"SourceFile", cmd.sourceFile, ""},
		{"ExplicitSourceFile", cmd.explicitSourceFile, filepath.Join(pwd, "zzz.java")},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestJbangJshWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "jsh-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "foo"})
	cmd := FindJbang(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "jbang")},
		{"SourceFile", cmd.sourceFile, filepath.Join(pwd, "hello.jsh")},
		{"ExplicitSourceFile", cmd.explicitSourceFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}

	cmd.doConfigureJbang()
	if len(cmd.args.Args) != 2 {
		t.Errorf("invalid arg count")
	}
}

func TestJbangJshWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "jsh-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "foo"})
	cmd := FindJbang(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "jbang")},
		{"SourceFile", cmd.sourceFile, filepath.Join(pwd, "hello.jsh")},
		{"ExplicitSourceFile", cmd.explicitSourceFile, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestJbangJshWithExplicitFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "jsh-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	args := ParseArgs([]string{"-gq", "zzz.jsh", "foo"})
	cmd := FindJbang(context, &args)

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
		return
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "jbang")},
		{"SourceFile", cmd.sourceFile, ""},
		{"ExplicitSourceFile", cmd.explicitSourceFile, filepath.Join(pwd, "zzz.jsh")},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestJbangWithoutExecutables(t *testing.T) {
	// given:
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "jbang", "java-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{}}

	// when:
	args := ParseArgs([]string{"-gq"})
	cmd := FindJbang(context, &args)

	// then:
	if cmd != nil {
		t.Error("Expected a nil command but got something")
	}
}
