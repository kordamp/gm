package gum

// Command defines an executable command (gradle/maven)
type Command interface {
	// Executes the given command
	Execute()
}
