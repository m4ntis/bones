package dbg

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/cpu"
	"github.com/peterh/liner"
)

type Debugger struct {
	n *bones.NES
	b bones.Break

	// hashed maps a Command to it's name and aliases.
	hashed map[string]*Command

	// alias maps all address aliases
	alias map[string]int

	line      *liner.State
	lastInput string
}

func New(nes *bones.NES) *Debugger {
	vec := nes.Vectors()

	dbg := &Debugger{
		n: nes,

		alias: map[string]int{
			"NMI":   vec[0],
			"Reset": vec[1],
			"IRQ":   vec[2],
		},

		line: liner.NewLiner(),
	}

	dbg.hashed = hash(commands(dbg))
	return dbg
}

func (dbg *Debugger) Run() {
	fmt.Println("Type 'help' for list of commands.")
	dbg.b = <-dbg.n.Breaks

	for {
		dbg.handleUserInput()
	}
}

func (dbg *Debugger) handleUserInput() {
	var input string
	defer func() { dbg.lastInput = input }()

	for input == "" {
		input, _ = dbg.line.Prompt("(dbg) ")
		if input == "" {
			input = dbg.lastInput
		} else {
			dbg.line.AppendHistory(input)
		}
	}

	args := strings.Fields(input)

	cmd, ok := dbg.hashed[args[0]]
	if !ok {
		fmt.Printf("%s isn't a valid command, type 'help' for a list\n", input)
		return
	}

	cmd.Run(args[1:])
}

// hash creates a mapping of each command to it's name and alias.
func hash(cmds []Command) map[string]*Command {
	hashed := map[string]*Command{}

	for i := range cmds {
		hashed[cmds[i].name] = &cmds[i]

		for _, a := range cmds[i].alias {
			hashed[a] = &cmds[i]
		}
	}

	return hashed
}

func (dbg *Debugger) waitBreak() {
	dbg.b = <-dbg.n.Breaks
	List(dbg.b)
}

func (dbg *Debugger) parseAddr(txt string) int {
	// Test if argument is an addr alias
	addr, ok := dbg.alias[txt]
	if ok {
		return addr
	}

	addr64, _ := strconv.ParseInt(txt, 16, 32)
	return int(addr64)
}

func (dbg *Debugger) aliasCmd(args []string) {
	addr := dbg.parseAddr(args[0])

	dbg.alias[args[1]] = addr
	fmt.Printf("Alias set for $%04x\n", addr)
}

func (dbg *Debugger) aliasesCmd(args []string) {
	for a, addr := range dbg.alias {
		fmt.Printf("%04x: %s\n", addr, a)
	}
}

func (dbg *Debugger) breakCmd(args []string) {
	addr := dbg.parseAddr(args[0])

	dbg.n.Break(int(addr))
	fmt.Printf("Breakpoint set at $%04x\n", addr)
}

func (dbg *Debugger) breakpointsCmd(args []string) {
	for _, addr := range dbg.n.List() {
		fmt.Printf("%04x\n", addr)
	}
}

func (dbg *Debugger) clearCmd(args []string) {
	// TODO: support windows :(
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (dbg *Debugger) continueCmd(args []string) {
	dbg.n.Continue()
	dbg.waitBreak()
}

func (dbg *Debugger) deleteCmd(args []string) {
	addr := dbg.parseAddr(args[0])

	ok := dbg.n.Delete(int(addr))
	if !ok {
		fmt.Printf("No breakpoint set at $%04x\n", addr)
		return
	}

	fmt.Printf("Deleted breakpoint at $%04x\n", addr)
}

func (dbg *Debugger) deleteallCmd(args []string) {
	dbg.n.DeleteAll()
}

func (dbg *Debugger) exitCmd(args []string) {
	os.Exit(0)
}

func (dbg *Debugger) helpCmd(args []string) {
	dbg.printHelp(args)
}

func (dbg *Debugger) listCmd(args []string) {
	List(dbg.b)
}

func (dbg *Debugger) nextCmd(args []string) {
	dbg.n.Next()
	dbg.waitBreak()
}

func (dbg *Debugger) printCmd(args []string) {
	// An error check is guranteed to be unnecessary as bounds are
	// checked earlier

	addr := dbg.parseAddr(args[0])
	d, _ := dbg.n.RAM().Observe(int(addr))

	fmt.Printf("$%04x: 0x%02x\n", int(addr), d)
}

func (dbg *Debugger) regsCmd(args []string) {
	fmt.Println(strings.Trim(fmt.Sprintf("%+v", dbg.n.Reg()), "&{}"))
}

func (dbg *Debugger) vprintCmd(args []string) {
	addr := dbg.parseAddr(args[0])

	fmt.Printf("$%04x: 0x%02x\n", int(addr), dbg.n.VRAM().Read(int(addr)))
}

func (dbg *Debugger) validAliasCmd(args []string) (ok bool) {
	ok = argsLenValidator([]int{2})(args)
	if !ok {
		return false
	}
	ok = dbg.argsAddrValidator(cpu.RAMSize)(args[:1])
	if !ok {
		return false
	}

	// Return false if alias taken
	addr, ok := dbg.alias[args[1]]
	if ok {
		fmt.Printf("Alias '%s' already taken for $%04x\n",
			args[1], addr)
		return false
	}

	return true
}

// argsAddrValidator creates an argument validator function that validates a
// single hex argument within memory bounds.
func (dbg *Debugger) argsAddrValidator(memSize int) func(args []string) (ok bool) {
	return func(args []string) (ok bool) {
		// Validate args len
		ok = argsLenValidator([]int{1})(args)
		if !ok {
			return false
		}

		// Test if arg is an addr alias
		_, ok = dbg.alias[args[0]]
		if ok {
			return true
		}

		// Try parsing arg into an int in addr space range
		addr, err := strconv.ParseInt(args[0], 16, 32)
		if err != nil || addr < 0 || addr >= int64(memSize) {
			fmt.Printf("Error: This command takes a single hex value between 0 and 0x%x\n",
				memSize)
			return false
		}

		return true
	}
}

func (dbg *Debugger) printHelp(args []string) {
	if len(args) == 0 {
		fmt.Println(dbg.generateHelp())
		return
	}

	cmd, ok := dbg.hashed[args[0]]
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
func (dbg *Debugger) generateHelp() string {
	// descrs maps a cmd's title to it's descr
	descrs := map[string]string{}

	// titles holds a sorted list of titles
	titles := make([]string, 0)

	// maxTitleLen holds the length of longest title
	maxTitleLen := 0

	// Map cmd title to descr
	for _, cmd := range dbg.hashed {
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
