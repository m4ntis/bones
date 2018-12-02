package dbg

import (
	"fmt"
	"strconv"
)

// argsLenValidator creates an argument validator function that validates length
// specified by a valid lens slice.
func argsLenValidator(lens []int) func(args []string) (ok bool) {
	return func(args []string) (ok bool) {
		for _, n := range lens {
			if len(args) == n {
				return true
			}
		}

		// Handle single valid arguments length
		if len(lens) == 1 {
			if lens[0] == 0 {
				fmt.Println("Error: This command takes no arguments")
				return false
			}
			if lens[0] == 1 {
				fmt.Println("Error: This command takes 1 argument")
				return false
			}

			fmt.Printf("Error: This command takes %d arguments\n", lens[0])
			return false
		}

		fmt.Printf("Error: This command takes %v arguments\n", lens)
		return false
	}
}

// argsAddrValidator creates an argument validator function that validates a
// single hex argument within memory bounds.
func argsAddrValidator(memSize int) func(args []string) (ok bool) {
	return func(args []string) (ok bool) {
		// Validate args len
		ok = argsLenValidator([]int{1})(args)
		if !ok {
			return false
		}

		addr, err := strconv.ParseInt(args[0], 16, 32)
		if err != nil || addr < 0 || addr >= int64(memSize) {
			fmt.Printf("Error: This command takes a single hex value between 0 and 0x%x\n",
				memSize)
			return false
		}

		return true
	}
}
