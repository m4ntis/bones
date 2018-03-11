package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/m4ntis/bones/dbg"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/models"
	"github.com/spf13/cobra"
)

// dbgCommand represents a command of the interactive debugger
//
// The return value of the function is set if the user interaction has ended,
// and the debugger should return to waiting for the next breakpoint to be hit.
type dbgCommand struct {
	name    string
	aliases []string

	cmd func(args []string) bool

	description string
	hString     string
}

var (
	dw *dbg.Worker

	breakVals chan dbg.BreakData

	dbgCommands map[string]*dbgCommand
	cmdsHelp    string

	// dbgCmd represents the dbg cli command
	dbgCmd = &cobra.Command{
		Use:   "dbg",
		Short: "Debug an iNES program",
		Long:  "The bones dbg command is used to debug NES roms, in iNES format.\n",
		Run: func(cmd *cobra.Command, args []string) {
			rom := openRom(args)
			breakVals = make(chan dbg.BreakData)
			dw = dbg.NewWorker(rom, breakVals)

			go dw.Start()
			startInteractiveDbg()
		},
	}
)

func openRom(args []string) *models.ROM {
	if len(args) != 1 {
		fmt.Println("Usage:\n  bones dbg <romname>.nes")
		os.Exit(1)
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s:\n%s\n", filename, err.Error())
		os.Exit(1)
	}

	rom, err := ines.Parse(f)
	if err != nil {
		fmt.Printf("Error parsing iNES file %s:\n%s\n", filename, err.Error())
		os.Exit(1)
	}

	return rom
}

func startInteractiveDbg() {
	fmt.Println("Type 'help' for list of commands.")
	for data := range breakVals {
		displayBreak(data)
		interact()
	}
}

func displayBreak(data dbg.BreakData) {
	instIdx := data.Disass.IndexOf(data.Reg.PC - 0x8000)
	inst := data.Disass.Code[instIdx]
	fmt.Printf("%+v\n", data.Reg)
	fmt.Printf("0x%04x: %s\n", inst.Addr, inst.Text)
}

func interact() {
	finished := false
	for !finished {
		finished = handleUserInput()
	}
}

func handleUserInput() (finished bool) {
	var input string
	for input == "" {
		fmt.Print("(dbg) ")
		fmt.Scanln(&input)
	}

	args := strings.Fields(input)

	cmd, ok := dbgCommands[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", input)
		return false
	}

	return cmd.cmd(args[1:])
}

func initCommands() {
	cmds := []dbgCommand{
		dbgCommand{
			name:    "next",
			aliases: []string{"n"},

			cmd: func(args []string) bool {
				dw.Next()
				return true
			},

			description: "Step over to next opcode",
			hString:     "",
		},
		dbgCommand{
			name:    "exit",
			aliases: []string{"quit", "q"},

			cmd: func(args []string) bool {
				os.Exit(0)
				return true
			},

			description: "Exit the debugger",
			hString:     "",
		},
		dbgCommand{
			name:    "help",
			aliases: []string{"h", "?"},

			cmd: func(args []string) bool {
				printHelp(args)
				return false
			},

			description: "Get list of commands or help on each",
			hString:     "",
		},
		dbgCommand{
			name:    "clear",
			aliases: []string{},

			cmd: func(args []string) bool {
				// TODO: support windows :(
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
				return false
			},

			description: "Get list of commands or help on each",
			hString:     "",
		},
	}

	dbgCommands = initCmdsMap(cmds)
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
		fmt.Println(cmdsHelp)
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

	fmt.Println(cmd.hString)
}

func init() {
	rootCmd.AddCommand(dbgCmd)

	initCommands()

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

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)
}
