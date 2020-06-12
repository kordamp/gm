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

const (
	// VERSION is current Gum version
	VERSION string = "0.4.0"
)

func main() {
	var args []string
	gradleBuild, args := gum.GrabFlag("-gg", os.Args[1:])
	mavenBuild, args := gum.GrabFlag("-gm", args)
	quiet, args := gum.GrabFlag("-gq", args)
	version, args := gum.GrabFlag("-gv", args)
	help, args := gum.GrabFlag("-gh", args)

	if version {
		fmt.Println("gm " + VERSION)
		os.Exit(0)
	}

	if help {
		fmt.Println("Usage of gm:")
		fmt.Println("  -gd\tdisplays debug information")
		fmt.Println("  -gg\tforce Gradle build")
		fmt.Println("  -gh\tdisplays help information")
		fmt.Println("  -gm\tforce Maven build")
		fmt.Println("  -gn\texecutes nearest build file")
		fmt.Println("  -gq\trun gm in quiet mode")
		fmt.Println("  -gr\tdo not replace goals/tasks")
		fmt.Println("  -gv\tdisplays version information")
		os.Exit(-1)
	}

	if gradleBuild && mavenBuild {
		fmt.Println("You cannot define both -gg and -gm flags at the same time")
		os.Exit(-1)
	}

	if gradleBuild {
		cmd := gum.FindGradle(gum.NewDefaultContext(quiet, true), args)
		cmd.Execute()
	} else if mavenBuild {
		cmd := gum.FindMaven(gum.NewDefaultContext(quiet, true), args)
		cmd.Execute()
	} else {
		findGradleOrMaven(quiet, args)
	}
}

// Attempts to execute gradlew/gradle first then mvnw/mvn
func findGradleOrMaven(quiet bool, args []string) {
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

	fmt.Println("Did not find a Gradle nor Maven project")
	os.Exit(-1)
}
