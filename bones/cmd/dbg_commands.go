package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/dbg"
	"github.com/m4ntis/bones/ppu"
)

// dbgCommand represents a command of the interactive debugger
//
// The return value of the function is set if the user interaction has ended,
// and the debugger should return to waiting for the next breakpoint to be hit.
type dbgCommand struct {
	name    string
	aliases []string

	cmd func(data dbg.BreakState, args []string) bool

	description string
	usage       string
	hString     string
}

func createCommands() map[string]*dbgCommand {
	cmds := []dbgCommand{
		dbgCommand{
			name:    "break",
			aliases: []string{"b"},

			cmd: func(data dbg.BreakState, args []string) bool {
				if len(args) != 1 {
					fmt.Println("break command takes exactly one argument")
					return false
				}

				addr, err := strconv.ParseInt(args[0], 16, 32)
				if err != nil {
					fmt.Println("break command only takes a numeric value in base 16")
					return false
				}

				ok := dw.Break(int(addr))
				if !ok {
					fmt.Printf("$%04x isn't a valid break address\n", addr)
					return false
				}

				fmt.Printf("Breakpoint set at $%04x\n", addr)
				return false
			},

			description: "Set a breakpoint",
			usage:       "break <address>",
			hString:     "Sets a breakpoint at the specified address in base 16",
		},
		dbgCommand{
			name:    "breakpoints",
			aliases: []string{"bps"},

			cmd: func(data dbg.BreakState, args []string) bool {
				for _, addr := range dw.List() {
					fmt.Printf("%04x\n", addr)
				}
				return false
			},

			description: "List breakpoints",
			usage:       "",
			hString:     "",
		},
		dbgCommand{
			name:    "delete",
			aliases: []string{"del", "d"},

			cmd: func(data dbg.BreakState, args []string) bool {
				if len(args) != 1 {
					fmt.Println("delete command takes exactly one argument")
					return false
				}

				addr, err := strconv.ParseInt(args[0], 16, 32)
				if err != nil {
					fmt.Println("delete command only takes a numeric value in base 16")
					return false
				}

				ok := dw.Delete(int(addr))
				if !ok {
					fmt.Printf("There is no breakpoint set at $%04x\n", addr)
					return false
				}

				fmt.Printf("Deleted breakpoint at $%04x\n", addr)
				return false
			},

			description: "Delete a breakpoint",
			usage:       "delete <address>",
			hString:     "Delete a set breakpoint at the specified address in base 16",
		},
		dbgCommand{
			name:    "deleteall",
			aliases: []string{"da"},

			cmd: func(data dbg.BreakState, args []string) bool {
				dw.DeleteAll()
				return false
			},

			description: "Delete all breakpoints",
			usage:       "",
			hString:     "",
		},
		dbgCommand{
			name:    "continue",
			aliases: []string{"c"},

			cmd: func(data dbg.BreakState, args []string) bool {
				dw.Continue()
				return true
			},

			description: "Let the CPU continue till the next break",
			usage:       "",
			hString:     "",
		},
		dbgCommand{
			name:    "next",
			aliases: []string{"n"},

			cmd: func(data dbg.BreakState, args []string) bool {
				dw.Next()
				return true
			},

			description: "Step over to next opcode",
			usage:       "",
			hString:     "",
		},
		dbgCommand{
			name:    "exit",
			aliases: []string{"quit", "q"},

			cmd: func(data dbg.BreakState, args []string) bool {
				os.Exit(0)
				return true
			},

			description: "Exit the debugger",
			usage:       "",
			hString:     "",
		},
		dbgCommand{
			name:    "help",
			aliases: []string{"h", "?"},

			cmd: func(data dbg.BreakState, args []string) bool {
				printHelp(args)
				return false
			},

			description: "Get list of commands or help on each",
			usage:       "help [command]",
			hString:     "Type 'help' to get a list of commands, or help about a specific command by appending it's name",
		},
		dbgCommand{
			name:    "clear",
			aliases: []string{},

			cmd: func(data dbg.BreakState, args []string) bool {
				// TODO: support windows :(
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
				return false
			},

			description: "Clear the screen",
			usage:       "",
			hString:     "",
		},
		dbgCommand{
			name:    "print",
			aliases: []string{"p"},

			cmd: func(data dbg.BreakState, args []string) bool {
				if len(args) != 1 {
					fmt.Println("print command takes exactly one argument")
					return false
				}

				addr, err := strconv.ParseInt(args[0], 16, 32)
				if err != nil || addr < 0 || addr >= cpu.RAMSize {
					fmt.Printf("print command only takes a numeric value between 0 and 0x%x\n", cpu.RAMSize)
					return false
				}

				// An error check is guranteed to be unnecessary as we checked
				// bounds earlier
				d, _ := data.RAM.Observe(int(addr))
				fmt.Printf("$%04x: 0x%02x\n", int(addr), d)
				return false
			},

			description: "Print a value from RAM",
			usage:       "print <address>",
			hString:     "Prints the hex value from RAM at a given address in hex",
		},
		dbgCommand{
			name:    "vprint",
			aliases: []string{"vp"},

			cmd: func(data dbg.BreakState, args []string) bool {
				if len(args) != 1 {
					fmt.Println("vprint command takes exactly one argument")
					return false
				}

				addr, err := strconv.ParseInt(args[0], 16, 32)
				if err != nil || addr > ppu.RAMSize {
					fmt.Printf("vprint command only takes a numeric value between 0 and 0x%x\n", cpu.RAMSize)
					return false
				}

				fmt.Printf("$%04x: 0x%02x\n", int(addr), data.VRAM.Read(int(addr)))
				return false
			},

			description: "Print a value from VRAM",
			usage:       "vprint <address>",
			hString:     "Prints the hex value from VRAM at a given address in hex",
		},
		dbgCommand{
			name:    "regs",
			aliases: []string{},

			cmd: func(data dbg.BreakState, args []string) bool {
				fmt.Println(strings.Trim(fmt.Sprintf("%+v", data.Reg), "&{}"))
				return false
			},

			description: "Prints the cpu's registers' status",
			usage:       "",
			hString:     "",
		},
		dbgCommand{
			name:    "list",
			aliases: []string{"ls"},

			cmd: func(data dbg.BreakState, args []string) bool {
				displayBreak(data)
				return false
			},

			description: "Display source code and current location",
			usage:       "",
			hString:     "",
		},
	}

	return initCmdsMap(cmds)
}

func initCmdsMap(cmds []dbgCommand) map[string]*dbgCommand {
	cmdsMap := map[string]*dbgCommand{}

	for i := range cmds {
		cmdsMap[cmds[i].name] = &cmds[i]

		for _, a := range cmds[i].aliases {
			cmdsMap[a] = &cmds[i]
		}
	}

	return cmdsMap
}

func printHelp(args []string) {
	if len(args) == 0 {
		fmt.Println(help)
		return
	}

	if len(args) > 1 {
		fmt.Println("'help' only takes up to 1 arguments")
		return
	}

	cmd, ok := dbgCommands[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", args[0])
		return
	}

	fmt.Println(cmd.description)
	if cmd.usage != "" {
		fmt.Println()
		fmt.Println("    " + cmd.usage)
	}
	if cmd.hString != "" {
		fmt.Println()
		fmt.Println(cmd.hString)
	}
}

func generateHelp() string {
	cmdDescriptions := map[string]string{}
	var sortedDescriptions []string
	var longestTitle int

	for _, cmd := range dbgCommands {
		title := generateCmdHelpTitle(cmd)
		if len(title) > longestTitle {
			longestTitle = len(title)
		}

		cmdDescriptions[title] = cmd.description
	}

	sortedDescriptions = make([]string, len(cmdDescriptions))
	i := 0
	for title := range cmdDescriptions {
		sortedDescriptions[i] = title
		i++
	}
	sort.Strings(sortedDescriptions)

	help := "The following commands are available:\n"
	for _, title := range sortedDescriptions {
		help += fmt.Sprintf("    %s %s %s\n", title,
			strings.Repeat("-", longestTitle-len(title)+1), cmdDescriptions[title])
	}

	help += "Type 'help' followed by a command for full documentation"
	return help
}

func generateCmdHelpTitle(cmd *dbgCommand) string {
	if len(cmd.aliases) == 0 {
		return cmd.name
	}
	return fmt.Sprintf("%s (alias: %s)", cmd.name, strings.Join(cmd.aliases, " | "))
}
