package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/m4ntis/bones/dbg"
	"github.com/m4ntis/bones/drawer"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/models"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

var (
	d = drawer.New()

	dw *dbg.Worker

	breakVals chan dbg.BreakData

	dbgCommands map[string]*dbgCommand
	help        string

	line = liner.NewLiner()

	// dbgCmd represents the dbg cli command
	dbgCmd = &cobra.Command{
		Use:   "dbg",
		Short: "Debug an iNES program",
		Long:  "The bones dbg command is used to debug NES roms, in iNES format.\n",
		Run: func(cmd *cobra.Command, args []string) {
			rom := openRom(args)
			breakVals = make(chan dbg.BreakData)
			dw = dbg.NewWorker(rom, breakVals, d)

			go dw.Start()
			go startInteractiveDbg()
			d.Run()
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
		interact(data)
	}
}

func displayBreak(data dbg.BreakData) {
	var startIdx int
	var endIdx int

	instIdx := data.Disass.IndexOf(data.Reg.PC)

	if instIdx > 5 {
		startIdx = instIdx - 5
	}

	if instIdx+5 < len(data.Disass.Code) {
		endIdx = instIdx + 5
	}

	for i := startIdx; i <= endIdx; i++ {
		inst := data.Disass.Code[i]

		if i == instIdx {
			fmt.Printf("=> 0x%04x: %s\n", inst.Addr, inst.Text)
			continue
		}
		fmt.Printf("   0x%04x: %s\n", inst.Addr, inst.Text)
	}
}

func interact(data dbg.BreakData) {
	finished := false
	for !finished {
		finished = handleUserInput(data)
	}
}

func handleUserInput(data dbg.BreakData) (finished bool) {
	var input string
	for input == "" {
		input, _ = line.Prompt("(dbg) ")
		line.AppendHistory(input)
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

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)
}
