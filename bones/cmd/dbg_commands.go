package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func createCommands() map[string]*dbgCommand {
	cmds := []dbgCommand{
		dbgCommand{
			name:    "next",
			aliases: []string{"n"},

			cmd: func(args []string) bool {
				dw.Next()
				return true
			},

			description: "Step over to next opcode",
			usage:       "",
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
			usage:       "",
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
			usage:       "help [command]",
			hString:     "Type 'help' to get a list of commands, or help about a specific command by appending it's name",
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

			description: "Clear the screen",
			usage:       "",
			hString:     "",
		},
	}

	return initCmdsMap(cmds)
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
		fmt.Println(help)
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

	fmt.Println(cmd.description)
	if cmd.usage != "" {
		fmt.Println()
		fmt.Println("    " + cmd.usage)
	}
	if cmd.hString != "" {
		fmt.Println()
		fmt.Println(cmd.hString)
	}
}

func generateHelp() string {
	cmdDescriptions := map[string]string{}
	var sortedDescriptions []string
	var longestTitle int

	for _, cmd := range dbgCommands {
		title := generateCmdHelpTitle(cmd)
		if len(title) > longestTitle {
			longestTitle = len(title)
		}

		cmdDescriptions[title] = cmd.description
	}

	sortedDescriptions = make([]string, len(cmdDescriptions))
	i := 0
	for title, _ := range cmdDescriptions {
		sortedDescriptions[i] = title
		i++
	}
	sort.Strings(sortedDescriptions)

	help := "The following commands are available:\n"
	for _, title := range sortedDescriptions {
		help += fmt.Sprintf("    %s %s %s\n", title,
			strings.Repeat("-", longestTitle-len(title)+1), cmdDescriptions[title])
	}

	help += "Type 'help' followed by a command for full documentation"
	return help
}

func generateCmdHelpTitle(cmd *dbgCommand) string {
	if len(cmd.aliases) == 0 {
		return cmd.name
	}
	return fmt.Sprintf("%s (alias: %s)", cmd.name, strings.Join(cmd.aliases, " | "))
}
