package gum

import (
	"fmt"
	"os"
	"os/exec"
)

type Command struct {
	Executable string
	BuildFile  string
	Settings   string
	Args       []string
}

func (c *Command) Execute() {
	args := make([]string, 0)

	if len(c.BuildFile) > 0 {
		args = append(args, "-b")
		args = append(args, c.BuildFile)
	}
	if len(c.Settings) > 0 {
		args = append(args, "-c")
		args = append(args, c.Settings)
	}
	for i := range c.Args {
		args = append(args, c.Args[i])
	}

	fmt.Println(args)

	cmd := exec.Command(c.Executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func (c *Command) empty() bool {
	return len(c.Executable) < 1
}

func EmptyCmd() Command {
	return Command{}
}
