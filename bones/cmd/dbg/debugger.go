package dbg

import (
	"fmt"
	"os"
	"strconv"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/swerve"
)

type Debugger struct {
	n *bones.NES
	b bones.Break

	aliases map[string]int

	shell *swerve.Shell
}

func New(nes *bones.NES) *Debugger {
	vec := nes.Vectors()

	dbg := &Debugger{
		n: nes,

		aliases: map[string]int{
			"NMI":   vec[0],
			"Reset": vec[1],
			"IRQ":   vec[2],
		},
	}

	s := swerve.New("(dbg) ")
	s.Add(commands(dbg)...)
	dbg.shell = s

	return dbg
}

func (dbg *Debugger) Run() {
	fmt.Println("Type 'help' for a list of available commands.")
	dbg.b = <-dbg.n.Breaks

	dbg.shell.Run()
	os.Exit(0)
}

func (dbg *Debugger) waitBreak() {
	dbg.b = <-dbg.n.Breaks
	list(dbg.b)
}

func (dbg *Debugger) parseAddr(s string) int {
	// Test if argument is an addr alias
	addr, ok := dbg.aliases[s]
	if ok {
		return addr
	}

	addr64, _ := strconv.ParseInt(s, 16, 32)
	return int(addr64)
}

func (dbg *Debugger) validAliasCmd(p swerve.Prompt, args []string) (ok bool) {
	ok = argsLenValidator([]int{2})(p, args)
	if !ok {
		return false
	}
	ok = dbg.argsAddrValidator(cpu.RAMSize)(p, args[:1])
	if !ok {
		return false
	}

	// Return false if alias taken
	addr, ok := dbg.aliases[args[1]]
	if ok {
		p.Printf("Alias '%s' already taken for $%04x\n", args[1], addr)
		return false
	}

	return true
}

// argsAddrValidator creates an argument validator function that validates a
// single hex argument within memory bounds.
func (dbg *Debugger) argsAddrValidator(memSize int,
) func(swerve.Prompt, []string) bool {

	return func(p swerve.Prompt, args []string) (ok bool) {
		// Validate args len
		ok = argsLenValidator([]int{1})(p, args)
		if !ok {
			return false
		}

		// Test if arg is an addr alias
		_, ok = dbg.aliases[args[0]]
		if ok {
			return true
		}

		// Try parsing arg into an int in addr space range
		addr, err := strconv.ParseInt(args[0], 16, 32)
		if err != nil || addr < 0 || addr >= int64(memSize) {
			p.Printf("Error: This command takes a single hex value between 0 and 0x%x\n",
				memSize)
			return false
		}

		return true
	}
}

// argsLenValidator creates an argument validator function that validates length
// specified by a valid lens slice.
func argsLenValidator(lens []int) func(swerve.Prompt, []string) bool {
	return func(p swerve.Prompt, args []string) (ok bool) {
		for _, n := range lens {
			if len(args) == n {
				return true
			}
		}

		// Handle single valid arguments length
		if len(lens) == 1 {
			if lens[0] == 0 {
				p.Println("Error: This command takes no arguments")
				return false
			}
			if lens[0] == 1 {
				p.Println("Error: This command takes 1 argument")
				return false
			}

			p.Printf("Error: This command takes %d arguments\n", lens[0])
			return false
		}

		p.Printf("Error: This command takes %v arguments\n", lens)
		return false
	}
}

func list(b bones.Break) {
	for i, inst := range b.Code {
		if i == b.PCIdx {
			fmt.Printf("=> %04x: %s\n", inst.Addr, inst.Text)
			continue
		}
		fmt.Printf("   %04x: %s\n", inst.Addr, inst.Text)
	}

	if b.Err != nil {
		fmt.Println(b.Err)
	}
}
