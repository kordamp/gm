package main

import (
	"fmt"
	"os"

	"github.com/kordamp/gm/gum"
)

const (
	VERSION string = "0.1.0"
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
		fmt.Println("  -gg\tforce Gradle build")
		fmt.Println("  -gh\tdisplays help information")
		fmt.Println("  -gm\tforce Maven build")
		fmt.Println("  -gn\texecutes nearest build file")
		fmt.Println("  -gq\trun gm in quiet mode")
		fmt.Println("  -gv\tdisplays version information")
		os.Exit(-1)
	}

	if gradleBuild && mavenBuild {
		fmt.Println("You cannot define both -gg and -gm flags at the same time")
		os.Exit(-1)
	}

	var cmd gum.Command
	if gradleBuild {
		cmd = gum.FindGradle(quiet, true, args)
	} else if mavenBuild {
		cmd = gum.FindMaven(quiet, true, args)
	} else {
		cmd = findGradleOrMaven(quiet, args)
	}

	cmd.Execute()
}

// Attempts to execute gradlew/gradle first then mvnw/mvn
func findGradleOrMaven(quiet bool, args []string) gum.Command {
	cmd := gum.FindGradle(quiet, false, args)

	if !cmd.Empty() {
		return cmd
	}

	cmd = gum.FindMaven(quiet, false, args)

	if !cmd.Empty() {
		return cmd
	}

	return gum.EmptyCommand{}
}
