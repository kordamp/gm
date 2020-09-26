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

package main

import (
	"fmt"
	"os"

	"github.com/kordamp/gm/gum"
)

var gmVersion string
var gmBuildCommit string
var gmBuildTimestamp string

func main() {
	args := gum.ParseArgs(os.Args[1:])

	gradleBuild := args.HasGumFlag("gg")
	mavenBuild := args.HasGumFlag("gm")
	jbangBuild := args.HasGumFlag("gj")
	quiet := args.HasGumFlag("gq")
	version := args.HasGumFlag("gv")
	help := args.HasGumFlag("gh")

	if version {
		fmt.Println("------------------------------------------------------------")
		fmt.Println("gm " + normalize(gmVersion))
		fmt.Println("------------------------------------------------------------")
		fmt.Println("Build time: " + normalize(gmBuildTimestamp))
		fmt.Println("Revision:   " + normalize(gmBuildCommit))
		fmt.Println("------------------------------------------------------------")
		os.Exit(0)
	}

	if help {
		fmt.Println("Usage of gm:")
		fmt.Println("  -gd\tdisplays debug information")
		fmt.Println("  -gg\tforce Gradle build")
		fmt.Println("  -gh\tdisplays help information")
		fmt.Println("  -gj\tforce jbang execution")
		fmt.Println("  -gm\tforce Maven build")
		fmt.Println("  -gn\texecutes nearest build file")
		fmt.Println("  -gq\trun gm in quiet mode")
		fmt.Println("  -gr\tdo not replace goals/tasks")
		fmt.Println("  -gv\tdisplays version information")
		os.Exit(0)
	}

	count := 0
	if gradleBuild {
		count = count + 1
	}
	if mavenBuild {
		count = count + 1
	}
	if jbangBuild {
		count = count + 1
	}

	if count > 1 {
		fmt.Println("You cannot define -gg, -gm, or -gj flags at the same time")
		os.Exit(-1)
	}

	if gradleBuild {
		cmd := gum.FindGradle(gum.NewDefaultContext(quiet, true), &args)
		cmd.Execute()
	} else if mavenBuild {
		cmd := gum.FindMaven(gum.NewDefaultContext(quiet, true), &args)
		cmd.Execute()
	} else if jbangBuild {
		cmd := gum.FindJbang(gum.NewDefaultContext(quiet, true), &args)
		cmd.Execute()
	} else {
		findTool(quiet, &args)
	}
}

// Attempts to execute gradlew/gradle first then mvnw/mvn
func findTool(quiet bool, args *gum.ParsedArgs) {
	context := gum.NewDefaultContext(quiet, false)

	gradle := gum.FindGradle(context, args)
	if gradle != nil {
		gradle.Execute()
		os.Exit(0)
	}

	maven := gum.FindMaven(context, args)
	if maven != nil {
		maven.Execute()
		os.Exit(0)
	}

	jbang := gum.FindJbang(context, args)
	if jbang != nil {
		jbang.Execute()
		os.Exit(0)
	}

	fmt.Println("Did not find a Gradle, Maven, or jbang project")
	os.Exit(-1)
}

func normalize(s string) string {
	if len(s) > 0 {
		return s
	}
	return "undefined"
}
