package cmd

import (
	"bufio"
	"fmt"
	"os"
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
	usage       string
	hString     string
}

var (
	dw *dbg.Worker

	breakVals chan dbg.BreakData

	dbgCommands map[string]*dbgCommand
	help        string

	reader *bufio.Reader

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
		input, _ = reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
	}

	args := strings.Fields(input)

	cmd, ok := dbgCommands[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", input)
		return false
	}

	return cmd.cmd(args[1:])
}

func init() {
	rootCmd.AddCommand(dbgCmd)

	reader = bufio.NewReader(os.Stdin)
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

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)
}
