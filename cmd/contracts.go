package cmd

// Command represents a wg-access-server
// subcommand module
type Command interface {
	Name() string
	Run()
}
