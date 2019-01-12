package dbg

import (
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/ppu"
)

// commands inits all the interractive debugger's commands.
func commands(dbg *Debugger) []Command {
	return []Command{
		Command{
			name:  "break",
			alias: []string{"b"},

			cmd:       dbg.breakCmd,
			validArgs: dbg.argsAddrValidator(cpu.RAMSize),

			descr: "Set a breakpoint",
			usage: "break <address>",
			hstr:  "Sets a breakpoint at the specified address in hex",
		},
		Command{
			name:  "breakpoints",
			alias: []string{"bps"},

			cmd: dbg.breakpointsCmd,

			descr: "List breakpoints",
		},
		Command{
			name:  "delete",
			alias: []string{"del", "d"},

			cmd:       dbg.deleteCmd,
			validArgs: dbg.argsAddrValidator(cpu.RAMSize),

			descr: "Delete a breakpoint",
			usage: "delete <address>",
			hstr:  "Delete a set breakpoint at the specified address in hex",
		},
		Command{
			name:  "deleteall",
			alias: []string{"da"},

			cmd: dbg.deleteallCmd,

			descr: "Delete all breakpoints",
		},
		Command{
			name:  "alias",
			alias: []string{},

			cmd:       dbg.aliasCmd,
			validArgs: dbg.validAliasCmd,

			descr: "Alias an address",
			usage: "alias <address> <alias>",
			hstr:  "Alias an address with a name. Aliased addresses are handled like addresses, and can be used to set and delete breakpoints, as well as printing values from RAM.",
		},
		Command{
			name:  "aliases",
			alias: []string{},

			cmd: dbg.aliasesCmd,

			descr: "List aliases",
		},
		Command{
			name:  "continue",
			alias: []string{"c"},

			cmd: dbg.continueCmd,

			descr: "Run the CPU until next break or error",
		},
		Command{
			name:  "next",
			alias: []string{"n"},

			cmd: dbg.nextCmd,

			descr: "Step over to next opcode",
		},
		Command{
			name:  "exit",
			alias: []string{"quit", "q"},

			cmd: dbg.exitCmd,

			descr: "Exit the debugger",
		},
		Command{
			name:  "help",
			alias: []string{"h", "?"},

			cmd:       dbg.helpCmd,
			validArgs: argsLenValidator([]int{0, 1}),

			descr: "Get a list of commands or help on each",
			usage: "help [command]",
			hstr:  "Type 'help' to get a list of commands, or help about a specific command by appending it's name",
		},
		Command{
			name:  "clear",
			alias: []string{},

			cmd: dbg.clearCmd,

			descr: "Clear the screen",
		},
		Command{
			name:  "print",
			alias: []string{"p"},

			cmd:       dbg.printCmd,
			validArgs: dbg.argsAddrValidator(cpu.RAMSize),

			descr: "Print a value from RAM",
			usage: "print <address>",
			hstr:  "Prints the hex value from RAM at a given address in hex",
		},
		Command{
			name:  "vprint",
			alias: []string{"vp"},

			cmd:       dbg.vprintCmd,
			validArgs: dbg.argsAddrValidator(ppu.RAMSize),

			descr: "Print a value from VRAM",
			usage: "vprint <address>",
			hstr:  "Prints the hex value from VRAM at a given address in hex",
		},
		Command{
			name:  "regs",
			alias: []string{},

			cmd: dbg.regsCmd,

			descr: "Prints the cpu's registers' status",
		},
		Command{
			name:  "list",
			alias: []string{"ls"},

			cmd: dbg.listCmd,

			descr: "Display source code and current location",
		},
	}
}
