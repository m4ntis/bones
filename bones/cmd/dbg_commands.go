package cmd

import (
	"fmt"
	"os"
	"os/exec"
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
