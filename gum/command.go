package gum

import (
	"fmt"
	"os"
)

type Command interface {
	Execute()

	Empty() bool
}

type EmptyCommand struct {
}

func (c EmptyCommand) Execute() {
	fmt.Println("Did not find a Gradle nor Maven project.")
	os.Exit(-1)
}

func (c EmptyCommand) Empty() bool {
	return true
}
