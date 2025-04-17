package evm

import (
	"errors"
	"fmt"

	t "github.com/Manuelshub/go-EVM/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// Various execution errors
var (
	ErrStackUnderflow = t.ErrStackUnderflow
	ErrStackOverflow  = t.ErrStackOverflow
	ErrInvalidOpcode  = errors.New("invalid opcode")
	ErrOutOfGas       = errors.New("out of gas")
)

// Instruction represents a single EVM instruction
type Instruction struct {
	Execute    func(ctx *ExecutionContext) error
	GasCost    func(ctx *ExecutionContext) uint64
	Name       string
	StackPops  int
	StackPushs int
}

// InstructionTable maps opcodes to their implementations
var InstructionTable = map[t.Opcode]Instruction{
	// t.STOP is the executes the stop operation and halts the execution of the program
	// And it is the only operation that does not consume any gas.
	t.STOP: {
		Execute:    opStop,
		GasCost:    constGasFunc(t.GasTierZero),
		Name:       "STOP",
		StackPops:  0,
		StackPushs: 0,
	},
	// t.ADD is the addition operation. It operates on the last two elements of the stack
	t.ADD: {
		Execute:    opAdd,
		GasCost:    constGasFunc(t.GasTierVeryLow),
		Name:       "ADD",
		StackPops:  2,
		StackPushs: 1,
	},
	t.MUL: {
		Execute:    opMul,
		GasCost:    constGasFunc(t.GasTierLow),
		Name:       "MUL",
		StackPops:  2,
		StackPushs: 1,
	},
	t.SUB: {
		Execute:    opSub,
		GasCost:    constGasFunc(t.GasTierVeryLow),
		Name:       "SUB",
		StackPops:  2,
		StackPushs: 1,
	},
	t.DIV: {
		Execute:    opDiv,
		GasCost:    constGasFunc(t.GasTierLow),
		Name:       "DIV",
		StackPops:  2,
		StackPushs: 1,
	},
	t.MOD: {
		Execute:    opMod,
		GasCost:    constGasFunc(t.GasTierLow),
		Name:       "MOD",
		StackPops:  2,
		StackPushs: 1,
	},
	t.EXP: {
		Execute:    opExp,
		GasCost:    gasExp,
		Name:       "EXP",
		StackPops:  2,
		StackPushs: 1,
	},
	t.AND: {
		Execute:    opAnd,
		GasCost:    constGasFunc(t.GasTierVeryLow),
		Name:       "AND",
		StackPops:  2,
		StackPushs: 1,
	},
	t.OR: {
		Execute:    opOr,
		GasCost:    constGasFunc(t.GasTierVeryLow),
		Name:       "OR",
		StackPops:  2,
		StackPushs: 1,
	},
	t.XOR: {
		Execute:    opXor,
		GasCost:    constGasFunc(t.GasTierVeryLow),
		Name:       "XOR",
		StackPops:  2,
		StackPushs: 1,
	},
	t.NOT: {
		Execute:    opNot,
		GasCost:    constGasFunc(t.GasTierVeryLow),
		Name:       "NOT",
		StackPops:  1,
		StackPushs: 1,
	},
	t.MLOAD: {
		Execute:    opMload,
		GasCost:    gasMLoad,
		Name:       "MLOAD",
		StackPops:  1,
		StackPushs: 1,
	},
	t.MSTORE: {
		Execute:    opMstore,
		GasCost:    gasMStore,
		Name:       "MSTORE",
		StackPops:  2,
		StackPushs: 0,
	},
	t.MSTORE8: {
		Execute:    opMstore8,
		GasCost:    gasMStore8,
		Name:       "MSTORE8",
		StackPops:  2,
		StackPushs: 0,
	},
	t.JUMP: {
		Execute:    opJump,
		GasCost:    constGasFunc(t.GasTierMid),
		Name:       "JUMP",
		StackPops:  1,
		StackPushs: 0,
	},
	t.JUMPI: {
		Execute:    opJumpi,
		GasCost:    constGasFunc(t.GasTierHigh),
		Name:       "JUMPI",
		StackPops:  2,
		StackPushs: 0,
	},
	t.JUMPDEST: {
		Execute:    opJumpdest,
		GasCost:    constGasFunc(t.GasTierBase),
		Name:       "JUMPDEST",
		StackPops:  0,
		StackPushs: 0,
	},
	t.RETURN: {
		Execute:    opReturn,
		GasCost:    gasReturn,
		Name:       "RETURN",
		StackPops:  2,
		StackPushs: 0,
	},
	t.SLOAD: {
		Execute:    opSload,
		GasCost:    constGasFunc(t.GasTierSLoad),
		Name:       "SLOAD",
		StackPops:  1,
		StackPushs: 1,
	},
	t.SSTORE: {
		Execute:    opSstore,
		GasCost:    gasSstore,
		Name:       "SSTORE",
		StackPops:  2,
		StackPushs: 0,
	},
	t.CALLVALUE: {
		Execute:    opCallValue,
		GasCost:    constGasFunc(t.GasTierBase),
		Name:       "CALLVALUE",
		StackPops:  0,
		StackPushs: 1,
	},
}

// Initialize PUSH operations (PUSH1-PUSH32)
func init() {
	// Add PUSH1 to PUSH32 to the instruction table
	for i := 1; i <= 32; i++ {
		pushOp := t.Opcode(int(t.PUSH1) + i - 1)
		InstructionTable[pushOp] = Instruction{
			Execute:    makePush(i),
			GasCost:    constGasFunc(t.GasTierVeryLow),
			Name:       fmt.Sprintf("PUSH%d", i),
			StackPops:  0,
			StackPushs: 1,
		}
	}

	// Add DUP1 to DUP16 to the instruction table
	for i := 1; i <= 16; i++ {
		dupOp := t.Opcode(int(t.DUP1) + i - 1)
		InstructionTable[dupOp] = Instruction{
			Execute:    makeDup(i),
			GasCost:    constGasFunc(t.GasTierVeryLow),
			Name:       fmt.Sprintf("DUP%d", i),
			StackPops:  i,
			StackPushs: i + 1,
		}
	}

	// Add SWAP1 to SWAP16 to the instruction table
	for i := 1; i <= 16; i++ {
		swapOp := t.Opcode(int(t.SWAP1) + i - 1)
		InstructionTable[swapOp] = Instruction{
			Execute:    makeSwap(i),
			GasCost:    constGasFunc(t.GasTierVeryLow),
			Name:       fmt.Sprintf("SWAP%d", i),
			StackPops:  i + 1,
			StackPushs: i + 1,
		}
	}
}

// Helper function for fixed gas costs
func constGasFunc(gas uint64) func(*ExecutionContext) uint64 {
	return func(*ExecutionContext) uint64 {
		return gas
	}
}

// ===== Instruction Implementations =====

// STOP halts execution
func opStop(ctx *ExecutionContext) error {
	ctx.Stopped = true
	return nil
}

// ADD implements x + y
func opAdd(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Add the two numbers
	result := uint256.NewInt(0)
	result.Add(x, y)

	return ctx.Stack.Push(result)
}

// MUL implements x * y
func opMul(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Multiply the two numbers
	result := uint256.NewInt(0)
	result.Mul(x, y)

	return ctx.Stack.Push(result)
}

// SUB implements x - y
func opSub(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Subtract the two numbers
	result := uint256.NewInt(0)
	result.Sub(x, y)

	return ctx.Stack.Push(result)
}

// DIV implements x / y
func opDiv(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Handle division by zero
	result := uint256.NewInt(0)
	if y.IsZero() {
		return ctx.Stack.Push(result)
	}

	// Divide the two numbers
	result.Div(x, y)

	return ctx.Stack.Push(result)
}

// Gas cost function for E
func gasExp(ctx *ExecutionContext) uint64 {
	// Check if we can access the stack safely
	if ctx.Stack.Size() < 2 {
		return 0
	}

	// The exponent (y in x^y) is the 2nd item on the stack
	exponent, err := ctx.Stack.GetItem(1) // 1 is the second item from the top (0-indexed)
	if err != nil {
		return 0
	}

	exponentBytesLen := uint64((exponent.BitLen() + 7) / 8)

	// Base cost + cost per byte in the exponent
	return t.GasTierLow + t.GasTierLow*exponentBytesLen
}

// MOD implements x % y
func opMod(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Handle modulo by zero
	result := uint256.NewInt(0)
	if y.IsZero() {
		return ctx.Stack.Push(result)
	}

	// Modulo the two numbers
	result.Mod(x, y)

	return ctx.Stack.Push(result)
}

// EXP implements x^y (x to the power of y)
func opExp(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	base, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	exponent, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Calculate the exponentiation
	result := uint256.NewInt(0)
	if exponent.IsZero() {
		result = uint256.NewInt(1)
	} else if base.IsZero() {
		result = uint256.NewInt(0)
	} else {
		result.Exp(base, exponent)
	}

	return ctx.Stack.Push(result)
}

// AND implements x & y (bitwise AND)
func opAnd(ctx *ExecutionContext) error {
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Bitwise AND
	result := uint256.NewInt(0)
	result.And(x, y)

	return ctx.Stack.Push(result)
}

// OR implements x | y (bitwise OR)
func opOr(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Bitwise OR
	result := uint256.NewInt(0)
	result.Or(x, y)

	return ctx.Stack.Push(result)
}

// XOR implements x ^ y (bitwise XOR)
func opXor(ctx *ExecutionContext) error {
	// Pop the last two elements off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	y, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Bitwise XOR
	result := uint256.NewInt(0)
	result.Xor(x, y)

	return ctx.Stack.Push(result)
}

// NOT implements ~x (bitwise NOT)
func opNot(ctx *ExecutionContext) error {
	// Pop the last element off the stack for this operation
	x, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Bitwise NOT
	result := uint256.NewInt(0)
	result.Not(x)

	return ctx.Stack.Push(result)
}

// ===== Memory Operations =====

// Gas cost calculation for memory operations
func memoryGasCost(ctx *ExecutionContext, additionalSize uint64) uint64 {
	oldSize := ctx.Memory.Size()
	newSize := oldSize

	if additionalSize > 0 {
		newSize = additionalSize
		if newSize < oldSize {
			newSize = oldSize
		}
	}

	// If memory doesn't expand, no additional cost
	if newSize <= oldSize {
		return t.GasTierVeryLow // Base cost for the operation
	}

	// Calculate memory expansion cost
	memoryCost := t.CalculateMemoryGasCost(oldSize, newSize)
	return t.GasTierVeryLow + memoryCost
}

// Gas cost for MLOAD
func gasMLoad(ctx *ExecutionContext) uint64 {
	// Check if we can access the stack
	if ctx.Stack.Size() < 1 {
		return 0
	}

	offset, err := ctx.Stack.GetItem(0)
	if err != nil {
		return 0
	}

	// Gas cost is base cost + memory expansion cost if applicable
	// We need to load 32 bytes from the offset
	return memoryGasCost(ctx, offset.Uint64()+32)
}

// Gas cost for MSTORE
func gasMStore(ctx *ExecutionContext) uint64 {
	// Check if we can access the stack
	if ctx.Stack.Size() < 2 {
		return 0
	}

	offset, err := ctx.Stack.GetItem(0)
	if err != nil {
		return 0
	}

	// Gas cost is base cost + memory expansion cost if applicable
	// We need to store 32 bytes at the offset
	return memoryGasCost(ctx, offset.Uint64()+32)
}

// Gas cost for MSTORE8
func gasMStore8(ctx *ExecutionContext) uint64 {
	// Check if we can access the stack
	if ctx.Stack.Size() < 2 {
		return 0
	}

	offset, err := ctx.Stack.GetItem(0)
	if err != nil {
		return 0
	}

	// Gas cost is base cost + memory expansion cost if applicable
	// We need to store 1 byte at the offset
	return memoryGasCost(ctx, offset.Uint64()+1)
}

// MLOAD implements load word from memory
func opMload(ctx *ExecutionContext) error {
	// Pop the last element off the stack for this operation
	offset, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Load 32 bytes from memory at offset
	value := ctx.Memory.Mload(offset.Uint64())
	if value == nil {
		value = make([]byte, 32)
	}

	// Convert bytes to uint256
	result := uint256.NewInt(0)
	result.SetBytes(value)

	return ctx.Stack.Push(result)
}

// MSTORE implements save word to memory
func opMstore(ctx *ExecutionContext) error {
	offset, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	value, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Convert uint256 to bytes (32 bytes, big-endian)
	bytes := make([]byte, 32)
	value.WriteToSlice(bytes)

	// Store bytes in memory
	ctx.Memory.Mstore(offset.Uint64(), bytes)

	return nil
}

// MSTORE8 implements save single byte to memory
func opMstore8(ctx *ExecutionContext) error {
	offset, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	value, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Only take the least significant byte
	lsb := byte(value.Uint64() & 0xff)

	// Store single byte in memory
	ctx.Memory.MstoreByte(offset.Uint64(), lsb)

	return nil
}

// ===== Control Flow Operations =====

// validateJumpDest checks if the destination is a valid JUMPDEST
func validateJumpDest(ctx *ExecutionContext, dest uint64) bool {
	if dest >= uint64(len(ctx.ByteCode)) {
		return false
	}

	// Check if the destination is a JUMPDEST opcode
	return t.Opcode(ctx.ByteCode[dest]) == t.JUMPDEST
}

// JUMP implements unconditional jump
func opJump(ctx *ExecutionContext) error {
	// Pop the destination from the stack
	dest, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Convert the destination to uint64
	jumpDest := dest.Uint64()

	// Validate jump destination
	if !validateJumpDest(ctx, jumpDest) {
		return errors.New("invalid jump destination")
	}

	// Set the program counter to the destination
	// Subtract 1 because the PC will be incremented after this instruction
	ctx.ProgramCounter = jumpDest
	return nil
}

// JUMPI implements conditional jump
func opJumpi(ctx *ExecutionContext) error {
	// Pop the destination and condition from the stack
	dest, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	cond, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// If condition is true, jump to destination
	if !cond.IsZero() {
		// Convert the destination to uint64
		jumpDest := dest.Uint64()

		// Validate jump destination
		if !validateJumpDest(ctx, jumpDest) {
			return errors.New("invalid jump destination")
		}

		// Set the program counter to the destination
		// Subtract 1 because the PC will be incremented after this instruction
		ctx.ProgramCounter = jumpDest
	}
	return nil
}

// JUMPDEST marks a valid destination for jumps
func opJumpdest(ctx *ExecutionContext) error {
	// This opcode is simply a marker, no operation needed
	return nil
}

// Gas cost for RETURN
func gasReturn(ctx *ExecutionContext) uint64 {
	// Check if we can access the stack
	if ctx.Stack.Size() < 2 {
		return 0
	}

	// Get the offset and size from the stack
	offset, err := ctx.Stack.GetItem(0)
	if err != nil {
		return 0
	}

	size, err := ctx.Stack.GetItem(1)
	if err != nil {
		return 0
	}

	// Cost is base cost + memory expansion cost
	return memoryGasCost(ctx, offset.Uint64()+size.Uint64())
}

// RETURN stops execution and returns data from memory
func opReturn(ctx *ExecutionContext) error {
	// Pop offset and size from stack
	offset, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	size, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// If size is 0, return empty data
	if size.IsZero() {
		ctx.ReturnData = []byte{}
		ctx.Stopped = true
		return nil
	}

	// Get the memory slice at the offset
	memOffset := offset.Uint64()
	memSize := size.Uint64()

	// Expand memory if needed
	mem := ctx.Memory.Expand(memOffset, memSize)

	// Copy data to return
	ctx.ReturnData = make([]byte, memSize)
	copy(ctx.ReturnData, mem)

	// Stop execution
	ctx.Stopped = true
	return nil
}

// ===== Push Operations =====

// makePush creates a function to handle PUSH operations
func makePush(size int) func(ctx *ExecutionContext) error {
	return func(ctx *ExecutionContext) error {
		// Check if there are enough bytes in the bytecode
		if ctx.ProgramCounter+uint64(size) > uint64(len(ctx.ByteCode)) {
			return errors.New("push: insufficient data in bytecode")
		}

		// Read 'size' bytes from bytecode
		bytes := ctx.ByteCode[ctx.ProgramCounter : ctx.ProgramCounter+uint64(size)]

		// Convert to uint256
		value := uint256.NewInt(0)
		value.SetBytes(bytes)

		// Push value to stack
		err := ctx.Stack.Push(value)
		if err != nil {
			return err
		}

		// Skip the pushed bytes (they're now processed)
		ctx.ProgramCounter += uint64(size)

		return nil
	}
}

// ===== Stack Manipulation Operations =====

// makeDup creates a function to handle DUP operations
func makeDup(n int) func(ctx *ExecutionContext) error {
	return func(ctx *ExecutionContext) error {
		// DUP n duplicates the n-th stack item (0-indexed)
		return ctx.Stack.Dup(n)
	}
}

// makeSwap creates a function to handle SWAP operations
func makeSwap(n int) func(ctx *ExecutionContext) error {
	return func(ctx *ExecutionContext) error {
		// SWAP n swaps the top and (n+1)th stack items
		return ctx.Stack.Swap(n)
	}
}

// ===== Storage Operations =====

// Gas cost for SSTORE
func gasSstore(ctx *ExecutionContext) uint64 {
	// Check if we can access the stack
	if ctx.Stack.Size() < 2 {
		return 0
	}

	key, err := ctx.Stack.GetItem(0)
	if err != nil {
		return 0
	}

	val, err := ctx.Stack.GetItem(1)
	if err != nil {
		return 0
	}

	// Convert key to common.Hash
	keyHash := common.BytesToHash(key.Bytes())

	// Check if value already exists in storage
	currentValue := ctx.Storage.Sload(keyHash)
	currentValUint := uint256.NewInt(0)

	if currentValue != nil {
		currentValUint.SetBytes(currentValue)
	}

	// Cost depends on whether we're setting a new value, updating, or clearing
	if currentValue == nil && !val.IsZero() {
		// Creating a new storage entry
		return t.GasStorageSet
	} else if !currentValUint.IsZero() && val.IsZero() {
		// Clearing an existing entry (refund will be added separately)
		return t.GasStorageUpdate
	} else {
		// Updating an existing entry
		return t.GasStorageUpdate
	}
}

// SLOAD implements load word from storage
func opSload(ctx *ExecutionContext) error {
	// Pop the key from the stack
	key, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Convert key to common.Hash
	keyHash := common.BytesToHash(key.Bytes())

	// Load value from storage
	value := ctx.Storage.Sload(keyHash)

	// Convert to uint256 and push to stack
	result := uint256.NewInt(0)
	if value != nil {
		result.SetBytes(value)
	}

	return ctx.Stack.Push(result)
}

// SSTORE implements store word to storage
func opSstore(ctx *ExecutionContext) error {
	// Pop key and value from stack
	key, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	value, err := ctx.Stack.Pop()
	if err != nil {
		return err
	}

	// Convert key to common.Hash
	keyHash := common.BytesToHash(key.Bytes())

	// Check if value already exists in storage
	currentValue := ctx.Storage.Sload(keyHash)
	currentValUint := uint256.NewInt(0)
	if currentValue != nil {
		currentValUint.SetBytes(currentValue)
	}

	// Handle gas refund if we're clearing a storage slot
	if !currentValUint.IsZero() && value.IsZero() {
		ctx.GasMeter.RefundGas(t.GasStorageRefund)
	}

	// Store value
	valueBytes := value.Bytes()
	ctx.Storage.Sstore(keyHash, valueBytes[:])

	return nil
}

func opCallValue(ctx *ExecutionContext) error {
	// Get the value from the call
	value := ctx.CallValue
	if value == nil {
		return errors.New("CallValue is nil")
	}

	push_value := uint256.NewInt(0).Set(value)
	return ctx.Stack.Push(push_value)
}
