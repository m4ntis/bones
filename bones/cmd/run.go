package cmd

import (
	"github.com/m4ntis/bones/controller"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/display"
	"github.com/m4ntis/bones/nes"
	"github.com/spf13/cobra"
)

var (
	w *cpu.Worker

	// runCmd represents the run command
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run an iNES program",
		Long:  "The run command is used to run NES roms, in iNES format.\n",
		Run: func(cmd *cobra.Command, args []string) {
			rom := openRom(args)

			var ctrl *controller.Controller
			disp := display.New(ctrl)

			n = nes.New(rom, disp, ctrl, nes.ModeRun)

			go n.Start()
			disp.Run()
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)

	// Make bones run's usage be 'bones run <romname>.nes'
	runCmd.SetUsageTemplate(`Usage:
  bones run <romname>.nes{{if gt (len .Aliases) 0}}

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
