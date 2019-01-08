package cmd

import (
	"fmt"
	"strings"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/bones/cmd/dbg"
	"github.com/m4ntis/bones/io"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

var (
	n *bones.NES

	line      = liner.NewLiner()
	lastInput = ""
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
			rom := openRom(cmd.Use, args)

			ctrl := new(io.Controller)
			disp := io.NewDisplay(ctrl, displayFPS, scale)

			n = bones.New(disp, ctrl, bones.ModeDebug)
			n.Load(rom)
			dbg.Init(n)

			go n.Start()
			go startInteractiveDbg()

			disp.Run()
		},
	}
)

func startInteractiveDbg() {
	fmt.Println("Type 'help' for list of commands.")
	for b := range n.Breaks {
		dbg.List(b)
		interact(b)
	}
}

func interact(b bones.BreakState) {
	fin := false
	for !fin {
		fin = handleUserInput(b)
	}
}

func handleUserInput(b bones.BreakState) (fin bool) {
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

	cmd, ok := dbg.Cmds[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", input)
		return false
	}

	return cmd.Run(b, args[1:])
}

func init() {
	rootCmd.AddCommand(dbgCmd)

	flags := dbgCmd.Flags()

	flags.BoolVar(&displayFPS, "display-fps",
		false, "Display small FPS counter")
	flags.Float64VarP(&scale, "scale", "s", 4.0,
		"Set display scaling (240x256 * scale)")

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
