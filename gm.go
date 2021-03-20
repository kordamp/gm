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

	bachBuild := args.HasGumFlag("gb")
	gradleBuild := args.HasGumFlag("gg")
	mavenBuild := args.HasGumFlag("gm")
	jbangBuild := args.HasGumFlag("gj")
	antBuild := args.HasGumFlag("ga")
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
		fmt.Println("  -ga\tforce Ant build")
		fmt.Println("  -gb\tforce Bach build")
		fmt.Println("  -gc\tdisplays current configuration and quits")
		fmt.Println("  -gd\tdisplays debug information")
		fmt.Println("  -gg\tforce Gradle build")
		fmt.Println("  -gh\tdisplays help information")
		fmt.Println("  -gj\tforce JBang execution")
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
	if bachBuild {
		count = count + 1
	}
	if antBuild {
		count = count + 1
	}

	if count > 1 {
		fmt.Println("You cannot define -gb, -gg, -gm, gj, or -ga flags at the same time")
		os.Exit(-1)
	}

	if gradleBuild {
		gum.FindGradle(gum.NewDefaultContext(true), &args).Execute()
	} else if mavenBuild {
		gum.FindMaven(gum.NewDefaultContext(true), &args).Execute()
	} else if jbangBuild {
		gum.FindJbang(gum.NewDefaultContext(true), &args).Execute()
	} else if bachBuild {
		gum.FindBach(gum.NewDefaultContext(true), &args).Execute()
	} else if antBuild {
		gum.FindAnt(gum.NewDefaultContext(true), &args).Execute()
	} else {
		gum.FindTool(&args)
	}
}

func normalize(s string) string {
	if len(s) > 0 {
		return s
	}
	return "undefined"
}
