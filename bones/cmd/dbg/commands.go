package dbg

import (
	"fmt"
	"strings"

	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/ppu"
	"github.com/m4ntis/swerve"
)

// commands creates all the interractive debugger's commands.
func commands(dbg *Debugger) []swerve.Command {
	return []swerve.Command{
		swerve.Command{
			Name:    "break",
			Aliases: []string{"b"},

			Run: func(p swerve.Prompt, args []string) {
				addr := dbg.parseAddr(args[0])

				dbg.n.Break(int(addr))
				p.Printf("Breakpoint set at $%04x\n", addr)
			},
			ValidateArgs: dbg.argsAddrValidator(cpu.RAMSize),

			Desc:  "Set a breakpoint",
			Usage: "break <address>",
			Help:  "Sets a breakpoint at the specified address in hex",
		},
		swerve.Command{
			Name:    "breakpoints",
			Aliases: []string{"bps"},

			Run: func(p swerve.Prompt, args []string) {
				for _, addr := range dbg.n.List() {
					p.Printf("%04x\n", addr)
				}
			},

			Desc: "List breakpoints",
		},
		swerve.Command{
			Name:    "delete",
			Aliases: []string{"del", "d"},

			Run: func(p swerve.Prompt, args []string) {
				addr := dbg.parseAddr(args[0])

				ok := dbg.n.Delete(int(addr))
				if !ok {
					p.Printf("No breakpoint set at $%04x\n", addr)
					return
				}

				p.Printf("Deleted breakpoint at $%04x\n", addr)
			},
			ValidateArgs: dbg.argsAddrValidator(cpu.RAMSize),

			Desc:  "Delete a breakpoint",
			Usage: "delete <address>",
			Help:  "Delete a set breakpoint at the specified address in hex",
		},
		swerve.Command{
			Name:    "deleteall",
			Aliases: []string{"da"},

			Run: func(p swerve.Prompt, args []string) {
				dbg.n.DeleteAll()
			},

			Desc: "Delete all breakpoints",
		},
		swerve.Command{
			Name:    "alias",
			Aliases: []string{},

			Run: func(p swerve.Prompt, args []string) {
				addr := dbg.parseAddr(args[0])

				dbg.aliases[args[1]] = addr
				p.Printf("Alias set for $%04x\n", addr)
			},
			ValidateArgs: dbg.validAliasCmd,

			Desc:  "Alias an address",
			Usage: "alias <address> <alias>",
			Help:  "Alias an address with a name. Aliased addresses are handled like addresses, and can be used to set and delete breakpoints, as well as printing values from RAM.",
		},
		swerve.Command{
			Name:    "aliases",
			Aliases: []string{},

			Run: func(p swerve.Prompt, args []string) {
				for a, addr := range dbg.aliases {
					p.Printf("%04x: %s\n", addr, a)
				}
			},

			Desc: "List aliases",
		},
		swerve.Command{
			Name:    "continue",
			Aliases: []string{"c"},

			Run: func(p swerve.Prompt, args []string) {
				dbg.n.Continue()
				dbg.waitBreak()
			},

			Desc: "Run the CPU until next break or error",
		},
		swerve.Command{
			Name:    "next",
			Aliases: []string{"n"},

			Run: func(p swerve.Prompt, args []string) {
				dbg.n.Next()
				dbg.waitBreak()
			},

			Desc: "Step over to next opcode",
		},
		swerve.Command{
			Name:    "print",
			Aliases: []string{"p"},

			Run: func(p swerve.Prompt, args []string) {
				// An error check is guranteed to be unnecessary as bounds are
				// checked earlier

				addr := dbg.parseAddr(args[0])
				d, _ := dbg.n.RAM().Observe(int(addr))

				p.Printf("$%04x: 0x%02x\n", int(addr), d)
			},
			ValidateArgs: dbg.argsAddrValidator(cpu.RAMSize),

			Desc:  "Print a value from RAM",
			Usage: "print <address>",
			Help:  "Prints the hex value from RAM at a given address in hex",
		},
		swerve.Command{
			Name:    "vprint",
			Aliases: []string{"vp"},

			Run: func(p swerve.Prompt, args []string) {
				addr := dbg.parseAddr(args[0])

				p.Printf("$%04x: 0x%02x\n",
					int(addr), dbg.n.VRAM().Read(int(addr)))
			},
			ValidateArgs: dbg.argsAddrValidator(ppu.RAMSize),

			Desc:  "Print a value from VRAM",
			Usage: "vprint <address>",
			Help:  "Prints the hex value from VRAM at a given address in hex",
		},
		swerve.Command{
			Name:    "regs",
			Aliases: []string{},

			Run: func(p swerve.Prompt, args []string) {
				p.Println(strings.Trim(fmt.Sprintf("%+v", dbg.n.Reg()), "&{}"))
			},

			Desc: "Prints the cpu's registers' status",
		},
		swerve.Command{
			Name:    "list",
			Aliases: []string{"ls"},

			Run: func(p swerve.Prompt, args []string) {
				list(dbg.b)
			},

			Desc: "Display source code and current location",
		},
	}
}
