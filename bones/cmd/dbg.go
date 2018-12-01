package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/io"
	"github.com/m4ntis/bones/ppu"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

var (
	// dbgCmd represents the dbg cli command
	dbgCmd = &cobra.Command{
		Use:   "dbg",
		Short: "Debug an iNES program",
		Long: `The dbg command is used to debug NES roms, in iNES format.

The command prompts up an interactive gdb style debugger, waiting for user
input before executing each command. It has pretty basic functionality including
breakpoints, next and continue instructions, as well as printing registers, ram
and vram values.

The debugger also opens a separate window for displaying the ppu's output. It
may appear as frozen or 'Not responding' or anything of the sort, but that is
because frames are renedered on it only when the ppu finishes a frame and the
display is blocked on that.

For full documentation and options run the 'help' command in the interractive
debugger.
`,
		Run: func(cmd *cobra.Command, args []string) {
			rom := openRom(args)

			ctrl := new(io.Controller)
			disp := io.NewDisplay(ctrl, false, 4.0)

			n = bones.New(disp, ctrl, bones.ModeDebug)

			go n.Start(rom)
			go startInteractiveDbg()
			disp.Run()
		},
	}
)

func openRom(args []string) *ines.ROM {
	if len(args) != 1 {
		fmt.Println("Usage:\n  bones dbg <romname>.nes")
		os.Exit(1)
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s:\n%s\n", filename, err)
		os.Exit(1)
	}

	rom, err := ines.Parse(f)
	if err != nil {
		fmt.Printf("Error parsing iNES file %s:\n%s\n", filename, err)
		os.Exit(1)
	}

	return rom
}

func startInteractiveDbg() {
	fmt.Println("Type 'help' for list of commands.")
	for data := range n.Breaks {
		displayBreak(data)
		interact(data)
	}
}

func displayBreak(data bones.BreakState) {
	for i, inst := range data.Code {
		if i == data.PCIdx {
			fmt.Printf("=> %04x: %s\n", inst.Addr, inst.Text)
			continue
		}
		fmt.Printf("   %04x: %s\n", inst.Addr, inst.Text)
	}

	if data.Err != nil {
		fmt.Println(data.Err)
	}
}

func interact(data bones.BreakState) {
	finished := false
	for !finished {
		finished = handleUserInput(data)
	}
}

func handleUserInput(data bones.BreakState) (finished bool) {
	var input string
	defer func() { lastInput = input }()

	for input == "" {
		input, _ = line.Prompt("(dbg) ")
		line.AppendHistory(input)
		if input == "" {
			input = lastInput
		}
	}

	args := strings.Fields(input)

	cmd, ok := cmdsMap[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", input)
		return false
	}

	return cmd.run(data, args[1:])
}

func init() {
	rootCmd.AddCommand(dbgCmd)

	cmdsMap = mapCmds(cmds)
	help = generateHelp(cmds)

	// Make bones dbg's usage be 'bones dbg <romname>.nes'
	dbgCmd.SetUsageTemplate(`Usage:
  bones dbg <romname>.nes{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

// dbgCommand represents a command of the interactive debugger
//
// The return value of the function is set if the user interaction has ended,
// and the debugger should return to waiting for the next breakpoint to be hit.
type dbgCommand struct {
	name    string
	aliases []string

	cmd  func(data bones.BreakState, args []string) bool
	argc int

	description string
	usage       string
	hString     string
}

func (d *dbgCommand) run(data bones.BreakState, args []string) (ok bool) {
	ok = d.validateArgsLength(args)
	if !ok {
		return false
	}

	return d.cmd(data, args)
}

func (d *dbgCommand) generateCmdHelpTitle() string {
	if len(d.aliases) == 0 {
		return d.name
	}

	return fmt.Sprintf("%s (alias: %s)", d.name, strings.Join(d.aliases, " | "))
}

func (d *dbgCommand) validateArgsLength(args []string) (ok bool) {
	if d.argc == 0 {
		return true
	}

	if len(args) != d.argc {
		if d.argc == 1 {
			fmt.Printf("%s command takes exactly 1 argument\n", d.name)
			return false
		}

		fmt.Printf("%s command takes exactly %d arguments\n", d.name, d.argc)
		return false
	}

	return true
}

// mapCmds maps commands to their title.
func mapCmds(cmds []dbgCommand) map[string]*dbgCommand {
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
		fmt.Println("help command only takes up to 1 arguments")
		return
	}

	cmd, ok := cmdsMap[args[0]]
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

// generateHelp generates an alphabetically sorted, multi-line help string for a
// list of commands, based on their name, aliases and description.
func generateHelp(cmds []dbgCommand) string {
	// descriptions maps a cmd's title to it's description
	descriptions := map[string]string{}

	// titles holds a sorted list of titles
	titles := make([]string, 0)

	// maxTitleLen holds the length of longest title
	maxTitleLen := 0

	// Map cmd title to description
	for _, cmd := range cmds {
		title := cmd.generateCmdHelpTitle()
		if len(title) > maxTitleLen {
			maxTitleLen = len(title)
		}

		descriptions[title] = cmd.description
	}

	// Sort titles
	for title := range descriptions {
		titles = append(titles, title)
	}
	sort.Strings(titles)

	help := "The following commands are available:\n"
	for _, title := range titles {
		help += fmt.Sprintf("    %s %s %s\n",
			title,
			strings.Repeat("-", maxTitleLen-len(title)+1),
			descriptions[title])
	}

	help += "Type 'help' followed by a command for full documentation"
	return help
}

var (
	n *bones.NES

	cmdsMap map[string]*dbgCommand
	help    string

	line      = liner.NewLiner()
	lastInput = ""

	cmds = []dbgCommand{
		dbgCommand{
			name:    "break",
			aliases: []string{"b"},

			cmd: func(data bones.BreakState, args []string) bool {
				addr, err := strconv.ParseInt(args[0], 16, 32)
				if err != nil {
					fmt.Println("break command only takes a numeric value in base 16")
					return false
				}

				n.Break(int(addr))
				fmt.Printf("Breakpoint set at $%04x\n", addr)

				return false
			},
			argc: 1,

			description: "Set a breakpoint",
			usage:       "break <address>",
			hString:     "Sets a breakpoint at the specified address in base 16",
		},
		dbgCommand{
			name:    "breakpoints",
			aliases: []string{"bps"},

			cmd: func(data bones.BreakState, args []string) bool {
				for _, addr := range n.List() {
					fmt.Printf("%04x\n", addr)
				}
				return false
			},

			description: "List breakpoints",
		},
		dbgCommand{
			name:    "delete",
			aliases: []string{"del", "d"},

			cmd: func(data bones.BreakState, args []string) bool {
				addr, err := strconv.ParseInt(args[0], 16, 32)
				if err != nil {
					fmt.Println("delete command only takes a numeric value in base 16")
					return false
				}

				ok := n.Delete(int(addr))
				if !ok {
					fmt.Printf("There is no breakpoint set at $%04x\n", addr)
					return false
				}

				fmt.Printf("Deleted breakpoint at $%04x\n", addr)
				return false
			},
			argc: 1,

			description: "Delete a breakpoint",
			usage:       "delete <address>",
			hString:     "Delete a set breakpoint at the specified address in base 16",
		},
		dbgCommand{
			name:    "deleteall",
			aliases: []string{"da"},

			cmd: func(data bones.BreakState, args []string) bool {
				n.DeleteAll()
				return false
			},

			description: "Delete all breakpoints",
		},
		dbgCommand{
			name:    "continue",
			aliases: []string{"c"},

			cmd: func(data bones.BreakState, args []string) bool {
				n.Continue()
				return true
			},

			description: "Continue running the CPU until next break or error",
		},
		dbgCommand{
			name:    "next",
			aliases: []string{"n"},

			cmd: func(data bones.BreakState, args []string) bool {
				n.Next()
				return true
			},

			description: "Step over to next opcode",
		},
		dbgCommand{
			name:    "exit",
			aliases: []string{"quit", "q"},

			cmd: func(data bones.BreakState, args []string) bool {
				os.Exit(0)
				return true
			},

			description: "Exit the debugger",
		},
		dbgCommand{
			name:    "help",
			aliases: []string{"h", "?"},

			cmd: func(data bones.BreakState, args []string) bool {
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

			cmd: func(data bones.BreakState, args []string) bool {
				// TODO: support windows :(
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
				return false
			},

			description: "Clear the screen",
		},
		dbgCommand{
			name:    "print",
			aliases: []string{"p"},

			cmd: func(data bones.BreakState, args []string) bool {
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
			argc: 1,

			description: "Print a value from RAM",
			usage:       "print <address>",
			hString:     "Prints the hex value from RAM at a given address in hex",
		},
		dbgCommand{
			name:    "vprint",
			aliases: []string{"vp"},

			cmd: func(data bones.BreakState, args []string) bool {
				addr, err := strconv.ParseInt(args[0], 16, 32)
				if err != nil || addr > ppu.RAMSize {
					fmt.Printf("vprint command only takes a numeric value between 0 and 0x%x\n", cpu.RAMSize)
					return false
				}

				fmt.Printf("$%04x: 0x%02x\n", int(addr), data.VRAM.Read(int(addr)))
				return false
			},
			argc: 1,

			description: "Print a value from VRAM",
			usage:       "vprint <address>",
			hString:     "Prints the hex value from VRAM at a given address in hex",
		},
		dbgCommand{
			name:    "regs",
			aliases: []string{},

			cmd: func(data bones.BreakState, args []string) bool {
				fmt.Println(strings.Trim(fmt.Sprintf("%+v", data.Reg), "&{}"))
				return false
			},

			description: "Prints the cpu's registers' status",
		},
		dbgCommand{
			name:    "list",
			aliases: []string{"ls"},

			cmd: func(data bones.BreakState, args []string) bool {
				displayBreak(data)
				return false
			},

			description: "Display source code and current location",
		},
	}
)
