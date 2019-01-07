package dbg

import (
	"fmt"
	"strings"

	"github.com/m4ntis/bones"
)

// Command represents a command of the interactive debugger.
type Command struct {
	name  string
	alias []string

	cmd       func(b bones.BreakState, args []string) (fin bool)
	validArgs func(args []string) (ok bool)

	descr string
	usage string
	hstr  string
}

// Run runs the command.
//
// fin indicates whether user interraction has finished.
func (c *Command) Run(b bones.BreakState, args []string) (fin bool) {
	if c.validArgs != nil {
		ok := c.validArgs(args)
		if !ok {
			return false
		}
	}

	return c.cmd(b, args)
}

func (c *Command) title() string {
	if len(c.alias) == 0 {
		return c.name
	}

	return fmt.Sprintf("%s (alias: %s)", c.name, strings.Join(c.alias, " | "))
}
