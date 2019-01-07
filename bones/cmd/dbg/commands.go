package dbg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/ppu"
)

var (
	// Cmds maps a Command to it's name and alias.
	Cmds map[string]*Command
	cmds []Command

	alias map[string]int

	n *bones.NES
)

func Init(nes *bones.NES) {
	n = nes

	cmds = commands()
	Cmds = mapCmds(cmds)

	vec := n.Vectors()
	alias = map[string]int{
		"NMI":   vec[0],
		"Reset": vec[1],
		"IRQ":   vec[2],
	}
}

// mapCmds creates a mapping of each command to it's name and alias.
func mapCmds(cmds []Command) map[string]*Command {
	cmdsMap := map[string]*Command{}

	for i := range cmds {
		cmdsMap[cmds[i].name] = &cmds[i]

		for _, a := range cmds[i].alias {
			cmdsMap[a] = &cmds[i]
		}
	}

	return cmdsMap
}

// commands inits all the interractive debugger's commands.
func commands() []Command {
	return []Command{
		Command{
			name:  "break",
			alias: []string{"b"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				addr := parseAddr(args[0])

				n.Break(int(addr))
				fmt.Printf("Breakpoint set at $%04x\n", addr)

				return false
			},
			validArgs: argsAddrValidator(cpu.RAMSize),

			descr: "Set a breakpoint",
			usage: "break <address>",
			hstr:  "Sets a breakpoint at the specified address in hex",
		},
		Command{
			name:  "breakpoints",
			alias: []string{"bps"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				for _, addr := range n.List() {
					fmt.Printf("%04x\n", addr)
				}
				return false
			},

			descr: "List breakpoints",
		},
		Command{
			name:  "delete",
			alias: []string{"del", "d"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				addr := parseAddr(args[0])

				ok := n.Delete(int(addr))
				if !ok {
					fmt.Printf("No breakpoint set at $%04x\n", addr)
					return false
				}

				fmt.Printf("Deleted breakpoint at $%04x\n", addr)
				return false
			},
			validArgs: argsAddrValidator(cpu.RAMSize),

			descr: "Delete a breakpoint",
			usage: "delete <address>",
			hstr:  "Delete a set breakpoint at the specified address in hex",
		},
		Command{
			name:  "deleteall",
			alias: []string{"da"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				n.DeleteAll()
				return false
			},

			descr: "Delete all breakpoints",
		},
		Command{
			name:  "alias",
			alias: []string{},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				addr := parseAddr(args[0])

				alias[args[1]] = addr
				fmt.Printf("Alias set for $%04x\n", addr)
				return false
			},
			validArgs: func(args []string) (ok bool) {
				ok = argsLenValidator([]int{2})(args)
				if !ok {
					return false
				}
				ok = argsAddrValidator(cpu.RAMSize)(args[:1])
				if !ok {
					return false
				}

				// Return false if alias taken
				addr, ok := alias[args[1]]
				if ok {
					fmt.Printf("Alias '%s' already taken for $%04x\n",
						args[1], addr)
					return false
				}

				return true
			},

			descr: "Alias an address",
			usage: "alias <address> <alias>",
			hstr:  "Alias an address with a name. Aliased addresses are handled like addresses, and can be used to set and delete breakpoints, as well as printing values from RAM.",
		},
		Command{
			name:  "aliases",
			alias: []string{},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				for a, addr := range alias {
					fmt.Printf("%04x: %s\n", addr, a)
				}
				return false
			},

			descr: "List aliases",
		},
		Command{
			name:  "continue",
			alias: []string{"c"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				n.Continue()
				return true
			},

			descr: "Run the CPU until next break or error",
		},
		Command{
			name:  "next",
			alias: []string{"n"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				n.Next()
				return true
			},

			descr: "Step over to next opcode",
		},
		Command{
			name:  "exit",
			alias: []string{"quit", "q"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				os.Exit(0)
				return true
			},

			descr: "Exit the debugger",
		},
		Command{
			name:  "help",
			alias: []string{"h", "?"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				printHelp(args)
				return false
			},
			validArgs: argsLenValidator([]int{0, 1}),

			descr: "Get a list of commands or help on each",
			usage: "help [command]",
			hstr:  "Type 'help' to get a list of commands, or help about a specific command by appending it's name",
		},
		Command{
			name:  "clear",
			alias: []string{},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				// TODO: support windows :(
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
				return false
			},

			descr: "Clear the screen",
		},
		Command{
			name:  "print",
			alias: []string{"p"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				// An error check is guranteed to be unnecessary as bounds are
				// checked earlier

				addr := parseAddr(args[0])
				d, _ := b.RAM.Observe(int(addr))

				fmt.Printf("$%04x: 0x%02x\n", int(addr), d)
				return false
			},
			validArgs: argsAddrValidator(cpu.RAMSize),

			descr: "Print a value from RAM",
			usage: "print <address>",
			hstr:  "Prints the hex value from RAM at a given address in hex",
		},
		Command{
			name:  "vprint",
			alias: []string{"vp"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				addr := parseAddr(args[0])

				fmt.Printf("$%04x: 0x%02x\n", int(addr), b.VRAM.Read(int(addr)))
				return false
			},
			validArgs: argsAddrValidator(ppu.RAMSize),

			descr: "Print a value from VRAM",
			usage: "vprint <address>",
			hstr:  "Prints the hex value from VRAM at a given address in hex",
		},
		Command{
			name:  "regs",
			alias: []string{},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				fmt.Println(strings.Trim(fmt.Sprintf("%+v", b.Reg), "&{}"))
				return false
			},

			descr: "Prints the cpu's registers' status",
		},
		Command{
			name:  "list",
			alias: []string{"ls"},

			cmd: func(b bones.BreakState, args []string) (fin bool) {
				List(b)
				return false
			},

			descr: "Display source code and current location",
		},
	}
}
