package gum

// Defines an executable command
type Command interface {
	// Executes the given command
	Execute()
}
