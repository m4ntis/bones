package dbg

import (
	"fmt"
	"strings"
)

// Command represents a command of the interactive debugger.
type Command struct {
	name  string
	alias []string

	cmd       func(args []string)
	validArgs func(args []string) (ok bool)

	descr string
	usage string
	hstr  string
}

// Run runs the command.
//
// fin indicates whether user interraction has finished.
func (c *Command) Run(args []string) {
	if c.validArgs != nil {
		ok := c.validArgs(args)
		if !ok {
			return
		}
	}

	c.cmd(args)
}

func (c *Command) title() string {
	if len(c.alias) == 0 {
		return c.name
	}

	return fmt.Sprintf("%s (alias: %s)", c.name, strings.Join(c.alias, " | "))
}
