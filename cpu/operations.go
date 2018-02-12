package cpu

// Operation defines an operation that the CPU executes in one or more of it's
// opcodes.
//
// The byte values it received are it's arguments. Both arguments and return
// value might be nil, depends on the operation. There isn't a gurantee that
// the operation will check for the correct number of arguments, so make sure
// you pass in the correct amount.
//
// The operation also gets a reference to the cpu's registers, in order to be
// able to test and change them
type Operation func(*Registers, ...byte) interface{}
