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
	"os"
	"strings"
)

// FindTool Executes gradle/maven/jbang based on config discovery
func FindTool(args *ParsedArgs) {
	context := NewDefaultContext(false)
	config := ReadUserConfig(context)
	config.merge(nil)

	if len(config.general.discovery) == 4 {
		discoverTool(config, context, args)
	}

	doFindGradle(context, args)
	doFindMaven(context, args)
	doFindJbang(context, args)
	doFindBach(context, args)

	if args.HasGumFlag("gc") {
		config.print()
		os.Exit(0)
	} else {
		fmt.Println("Did not find a Gradle, Maven, Bach, or JBang project")
		os.Exit(-1)
	}
}

func discoverTool(config *Config, context Context, args *ParsedArgs) {
	for i := range config.general.discovery {
		tool := strings.TrimSpace(strings.ToLower(config.general.discovery[i]))

		switch tool {
		case "gradle":
			doFindGradle(context, args)
			break
		case "maven":
			doFindMaven(context, args)
			break
		case "jbang":
			doFindJbang(context, args)
			break
		case "bach":
			doFindBach(context, args)
			break
		default:
			fmt.Println("Unsupported tool: " + tool)
			os.Exit(-1)
		}
	}
}

func doFindGradle(context Context, args *ParsedArgs) {
	gradle := FindGradle(context, args)
	if gradle != nil {
		gradle.Execute()
		os.Exit(0)
	}
}

func doFindMaven(context Context, args *ParsedArgs) {
	maven := FindMaven(context, args)
	if maven != nil {
		maven.Execute()
		os.Exit(0)
	}
}

func doFindJbang(context Context, args *ParsedArgs) {
	jbang := FindJbang(context, args)
	if jbang != nil {
		jbang.Execute()
		os.Exit(0)
	}
}

func doFindBach(context Context, args *ParsedArgs) {
	bach := FindBach(context, args)
	if bach != nil {
		bach.Execute()
		os.Exit(0)
	}
}
