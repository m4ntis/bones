package dbg

import (
	"fmt"
	"sort"
	"strings"
)

func printHelp(args []string) {
	if len(args) == 0 {
		fmt.Println(generateHelp())
		return
	}

	cmd, ok := Cmds[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", args[0])
		return
	}

	fmt.Println(cmd.descr)
	if cmd.usage != "" {
		fmt.Println()
		fmt.Println("    " + cmd.usage)
	}
	if cmd.hstr != "" {
		fmt.Println()
		fmt.Println(cmd.hstr)
	}
}

// generateHelp generates an alphabetically sorted, multi-line help string for
// the command list, based on their name, aliases and descr.
func generateHelp() string {
	// descrs maps a cmd's title to it's descr
	descrs := map[string]string{}

	// titles holds a sorted list of titles
	titles := make([]string, 0)

	// maxTitleLen holds the length of longest title
	maxTitleLen := 0

	// Map cmd title to descr
	for _, cmd := range Cmds {
		title := cmd.title()
		if len(title) > maxTitleLen {
			maxTitleLen = len(title)
		}

		descrs[title] = cmd.descr
	}

	// Sort titles
	for title := range descrs {
		titles = append(titles, title)
	}
	sort.Strings(titles)

	help := "The following commands are available:\n"
	for _, title := range titles {
		help += fmt.Sprintf("    %s %s %s\n",
			title,
			strings.Repeat("-", maxTitleLen-len(title)+1),
			descrs[title])
	}

	help += "Type 'help' followed by a command for full documentation"
	return help
}
