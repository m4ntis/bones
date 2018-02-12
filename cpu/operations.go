package cpu

// Operation defines an operation that the CPU executes in one or more of it's
// opcodes.
//
// The byte values it received are it's arguments. Arguments can be of any
// length, depending on the operation. There isn't a gurantee that the
// operation will check for the correct number of arguments, so make sure that
// you pass in the correct amount.
//
// The operation also gets a reference to the cpu so it can test and change the
// registers and RAM.
type Operation func(*CPU, ...byte)
