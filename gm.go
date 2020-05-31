package gum

import (
	"fmt"
	"os"
)

const (
	VERSION string = "0.1.0"
)

func main() {
	var args []string
	gradleBuild, args := IsFlagSet("-gg", os.Args[1:])
	mavenBuild, args := IsFlagSet("-gm", args)
	quiet, args := IsFlagSet("-gq", args)
	version, args := IsFlagSet("-gv", args)
	help, args := IsFlagSet("-gh", args)

	if version {
		fmt.Println("gm " + VERSION)
		os.Exit(0)
	}

	if help {
		fmt.Println("Usage of gm:")
		fmt.Println("\t-gg\tforce Gradle build")
		fmt.Println("\t-gh\tdisplays help information")
		fmt.Println("\t-gm\tforce Maven build")
		fmt.Println("\t-gq\trun gm in quiet mode")
		fmt.Println("\t-gv\tdisplays version information")
		os.Exit(-1)
	}

	if gradleBuild && mavenBuild {
		fmt.Println("You cannot define both -gg and -gm flags at the same time")
		os.Exit(-1)
	}

	var cmd Command
	if gradleBuild {
		cmd = FindGradle(quiet, true, args)
	} else if mavenBuild {
		cmd = FindMaven(quiet, true, args)
	} else {
		cmd = findGradleOrMaven(quiet, args)
	}

	cmd.Execute()
}

// Attempts to execute gradlew/gradle first then mvnw/mvn
func findGradleOrMaven(quiet bool, args []string) Command {
	cmd := FindGradle(quiet, false, args)

	if !cmd.empty() {
		return cmd
	}

	cmd = FindMaven(quiet, false, args)

	if !cmd.empty() {
		return cmd
	}

	fmt.Println("Did not find a Gradle nor Maven project.")
	os.Exit(-1)
	return EmptyCmd()
}
