package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/io"
	"github.com/spf13/cobra"
)

var (
	// benchCmd represents the run command
	benchCmd = &cobra.Command{
		Use:   "bench",
		Short: "Benchmark BoNES using an iNES program",
		Run: func(cmd *cobra.Command, args []string) {
			rom := openRom(cmd.Use, args)

			ctrl := new(io.Controller)
			disp := io.NewBenchDisplay()

			n := bones.New(disp, ctrl, bones.ModeRun)
			n.Load(rom)

			go n.Start()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
		},
	}
)

func init() {
	rootCmd.AddCommand(benchCmd)

	// Make bones bench's usage be 'bones bench <romname>.nes'
	benchCmd.SetUsageTemplate(`Usage:
  bones bench <romname>.nes{{if gt (len .Aliases) 0}}

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
