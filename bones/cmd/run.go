package cmd

import (
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/display"
	"github.com/m4ntis/bones/models"
	"github.com/spf13/cobra"
)

var (
	w *cpu.Worker

	// runCmd represents the run command
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run an iNES program",
		Long:  "The bones run command is used to run NES roms, in iNES format.\n",
		Run: func(cmd *cobra.Command, args []string) {
			ctrl := &models.Controller{}
			d := display.New(ctrl)
			rom := openRom(args)
			w = cpu.NewWorker(rom, d, ctrl)

			go w.Start()
			d.Run()
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

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)
}
