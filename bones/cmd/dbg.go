package cmd

import (
	"fmt"
	"os"

	"github.com/m4ntis/bones/disass"
	"github.com/m4ntis/bones/ines"
	"github.com/spf13/cobra"
)

// dbgCmd represents the dbg command
var dbgCmd = &cobra.Command{
	Use:   "dbg",
	Short: "Debug an iNES program",
	Long:  "The bones dbg command is used to debug NES roms, in iNES format.\n",
	Run: func(cmd *cobra.Command, args []string) {
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

		d := disass.Disassemble(rom.PrgROM)
		for _, inst := range d.Code {
			fmt.Printf("0x%04x: %s\n", inst.Addr, inst.Text)
		}
	},
}

func init() {
	rootCmd.AddCommand(dbgCmd)

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
