package cmd

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/io"
	"github.com/spf13/cobra"
)

var (
	// profCmd represents the prof command
	profCmd = &cobra.Command{
		Use:   "prof",
		Short: "Run an iNES program with go's profiler",
		Long:  "The prof command is used to profile bones\n",
		Run: func(cmd *cobra.Command, args []string) {
			go func() { http.ListenAndServe("localhost:6060", nil) }()

			rom := openRom(cmd.Use, args)

			ctrl := new(io.Controller)
			disp := io.NewDisplay(ctrl, displayFPS, scale)

			n = bones.New(disp, ctrl, bones.ModeRun)

			go n.Start(rom)
			disp.Run()
		},
	}
)

func init() {
	rootCmd.AddCommand(profCmd)

	// Make bones prof's usage be 'bones prof <romname>.nes'
	profCmd.SetUsageTemplate(`Usage:
  bones prof <romname>.nes{{if gt (len .Aliases) 0}}

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
