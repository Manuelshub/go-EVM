package evm

import (
	"errors"
	"fmt"

	t "github.com/Manuelshub/go-EVM/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// ExecutionContext is the context of the EVM executiont which contains the
// state of the execution
type ExecutionContext struct {
	ProgramCounter  uint64
	CallerAddress   common.Address
	CallValue       *uint256.Int
	CalleeAddress   common.Address
	ContractAddress common.Address
	Stack           *t.Stack
	Memory          *t.Memory
	Storage         *t.Storage
	GasMeter        *t.GasMeter
	ByteCode        []byte
	Stopped         bool   // Flag to indicate if execution should stop
	ReturnData      []byte // Data returned by RETURN or REVERT
	Error           error  // Last execution error
}

// NewExecutionContext creates a new ExecutionContext
func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		ProgramCounter: 0,
		CallValue:      uint256.NewInt(0),
		Stack:          t.NewStack(),
		Memory:         t.NewMemory(),
		Storage:        t.NewStorage(),
		GasMeter:       t.NewGasMeter(10000000), // Default gas limit
		Stopped:        false,
	}
}

// Run executes the bytecode
func (ctx *ExecutionContext) Run(bytecode []byte) ([]byte, error) {
	ctx.ByteCode = bytecode
	ctx.ProgramCounter = 0
	ctx.Stopped = false

	// Main execution loop
	for !ctx.Stopped && ctx.ProgramCounter < uint64(len(ctx.ByteCode)) {
		// Fetch the current opcode
		op := t.Opcode(ctx.ByteCode[ctx.ProgramCounter])

		// Look up the instruction
		instruction, exists := InstructionTable[op]
		if !exists {
			ctx.Error = ErrInvalidOpcode
			return nil, ctx.Error
		}

		// Consume gas
		gasCost := instruction.GasCost(ctx)
		if err := ctx.GasMeter.UseGas(gasCost); err != nil {
			ctx.Error = ErrOutOfGas
			return nil, ctx.Error
		}

		// Execute the instruction
		ctx.ProgramCounter++
		if err := instruction.Execute(ctx); err != nil {
			ctx.Error = err
			return nil, err
		}
	}

	return ctx.ReturnData, ctx.Error
}

// GetOpcodeName returns the name of an opcode
func GetOpcodeName(opcode byte) string {
	op := t.Opcode(opcode)
	instruction, exists := InstructionTable[op]
	if exists {
		return instruction.Name
	}
	return fmt.Sprintf("UNKNOWN (0x%x)", opcode)
}

// ExecuteStep executes a single instruction and advances the program counter
func ExecuteStep(ctx *ExecutionContext) error {
	if ctx.ProgramCounter >= uint64(len(ctx.ByteCode)) {
		return errors.New("end of code")
	}

	// Fetch the current opcode
	op := t.Opcode(ctx.ByteCode[ctx.ProgramCounter])

	// Look up the instruction
	instruction, exists := InstructionTable[op]
	if !exists {
		ctx.Error = ErrInvalidOpcode
		return fmt.Errorf("invalid opcode: 0x%x", op)
	}

	// Consume gas
	gasCost := instruction.GasCost(ctx)
	if err := ctx.GasMeter.UseGas(gasCost); err != nil {
		ctx.Error = ErrOutOfGas
		return err
	}

	// Execute the instruction
	ctx.ProgramCounter++
	if err := instruction.Execute(ctx); err != nil {
		ctx.Error = err
		return err
	}

	return nil
}
