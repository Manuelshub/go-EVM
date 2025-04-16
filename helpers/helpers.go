package helpers

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/Manuelshub/go-EVM/evm"
	"github.com/holiman/uint256"
)

// PrintHelp prints the help message for the CLI when `help` is entered
func PrintHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help               - Show this help")
	fmt.Println("  run <bytecode>     - Execute bytecode (hex format)")
	fmt.Println("  debug <bytecode>   - Execute bytecode with step-by-step tracing")
	fmt.Println("  stack              - Display current stack")
	fmt.Println("  storage <key>      - Display storage value at key (hex format)")
	fmt.Println("  push <value>       - Push a hex value onto the stack")
	fmt.Println("  reset              - Reset the execution context")
	fmt.Println("  exit, quit         - Exit the program")
}

// RunBytecode executes the given bytecode and prints the result
func RunBytecode(ctx *evm.ExecutionContext, hexString string) {
	if strings.HasPrefix(hexString, "0x") {
		hexString = hexString[2:]
	}

	bytecode, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Printf("Error decoding bytecode: %v\n", err)
		return
	}

	fmt.Printf("Running bytecode: 0x%s\n", hexString)
	result, err := ctx.Run(bytecode)
	if err != nil {
		fmt.Printf("Execution failed: %v\n", err)
	} else if result != nil {
		fmt.Printf("Execution successful. Result: 0x%s\n", hex.EncodeToString(result))
	} else {
		fmt.Println("Execution successful. No return value.")
	}

	fmt.Printf("Stack: %s\n", ctx.Stack.ToString())
	if ctx.GasMeter != nil {
		fmt.Printf("Gas used: %d\n", ctx.GasMeter.GasConsumed())
	}
}

func PushValue(ctx *evm.ExecutionContext, hexValue string) {
	if strings.HasPrefix(hexValue, "0x") {
		hexValue = hexValue[2:]
	}

	// Add leading zeros if less than 64 characters
	for len(hexValue) < 64 {
		hexValue = "0" + hexValue
	}

	// Decode the hex value
	bytes, err := hex.DecodeString(hexValue)
	if err != nil {
		fmt.Printf("Error decoding value: %v\n", err)
		return
	}

	// convert bytes to uint256
	value := uint256.NewInt(0)
	value.SetBytes(bytes)

	// push to stack
	err = ctx.Stack.Push(value)
	if err != nil {
		fmt.Printf("Error pushing to stack: %v\n", err)
		return
	}

	fmt.Printf("Pushed 0x%s to stack\n", hexValue)
	fmt.Printf("Stack: %s\n", ctx.Stack.ToString())
}

func DebugBytecode(ctx *evm.ExecutionContext, hexString string) {
	// Remove 0x prefix from string if present
	if strings.HasPrefix(hexString, "0x") {
		hexString = hexString[2:]
	}

	// Decode the hex string into bytes
	bytecode, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Printf("Error decoding bytecode: %v\n", err)
		return
	}

	fmt.Printf("Debugging bytecode: 0x%s\n", hexString)

	// Setup the context
	ctx.ByteCode = bytecode
	ctx.ProgramCounter = 0
	ctx.Stopped = false

	// Step through each instruction
	step := 1
	reader := bufio.NewReader(os.Stdin)

	for !ctx.Stopped && ctx.ProgramCounter < uint64(len(ctx.ByteCode)) {
		// Get the current opcode
		op := evm.GetOpcodeName(ctx.ByteCode[ctx.ProgramCounter])

		fmt.Printf("\nStep %d: PC=%d, Opcode=%s\n", step, ctx.ProgramCounter, op)
		fmt.Printf("Stack: %s\n", ctx.Stack.ToString())

		fmt.Printf("Press Enter to continue... (or 'q' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			fmt.Println("Debugging stopped")
			break
		}

		// Execute one instruction
		err := evm.ExecuteStep(ctx)
		if err != nil {
			fmt.Printf("Execution failed: %v\n", err)
			break
		}

		step++
	}

	if ctx.Stopped {
		fmt.Println("\nExecution stopped. Reason: STOP or RETURN")
	}

	if ctx.ReturnData != nil && len(ctx.ReturnData) > 0 {
		fmt.Printf("Return data: 0x%s\n", hex.EncodeToString(ctx.ReturnData))
	}

	fmt.Printf("Final stack: %s\n", ctx.Stack.ToString())
	fmt.Printf("Gas used: %d\n", ctx.GasMeter.GasConsumed())
}
