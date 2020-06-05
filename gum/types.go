package gum

// Command defines an executable command (gradle/maven)
type Command interface {
	// Execute executes the given command
	Execute()
}

// Context provides an abstraction over the OS and Environment as required by Gum
type Context interface {
	// IsQuiet whether Gum should stay silent or not
	IsQuiet() bool

	// IsExplicit whether a given tool was specified
	IsExplicit() bool

	// IsWindows checks if the current OS is Windows
	IsWindows() bool

	// GetWorkingDir returns the current working dir
	GetWorkingDir() string

	// GetPaths gets the paths in $PATH
	GetPaths() []string

	// FileExists checks if a file exists
	FileExists(name string) bool

	// Exit causes the current program to exit with the given status code.
	Exit(code int)
}
