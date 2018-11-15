package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/io"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

var (
	n *bones.NES

	dbgCommands map[string]*dbgCommand
	help        string

	line      = liner.NewLiner()
	lastInput = ""

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

	cmd, ok := dbgCommands[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", input)
		return false
	}

	return cmd.cmd(data, args[1:])
}

func init() {
	rootCmd.AddCommand(dbgCmd)

	dbgCommands = createCommands()
	help = generateHelp()

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
