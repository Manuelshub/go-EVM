package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/Manuelshub/go-EVM/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

func main() {
	fmt.Println("Go-EVM - A simple Ethereum Virtual Machine implementation")
	fmt.Println("Type 'help' for available commands")

	executionContext := evm.NewExecutionContext()

	for {
		fmt.Printf("(go-EVM) ")

		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "" {
			continue
		}

		parts := strings.Split(command, " ")
		cmd := parts[0]

		switch cmd {
		case "exit", "quit":
			os.Exit(0)

		case "help":
			printHelp()

		case "run":
			if len(parts) < 2 {
				fmt.Println("Error: Missing bytecode. Usage: run <bytecode>")
				continue
			}
			runBytecode(executionContext, parts[1])

		case "stack":
			fmt.Println(executionContext.Stack.ToString())

		case "push":
			if len(parts) < 2 {
				fmt.Println("Error: Missing value. Usage: push <hex_value>")
				continue
			}
			pushValue(executionContext, parts[1])

		case "reset":
			executionContext = evm.NewExecutionContext()
			fmt.Println("Execution context reset")

		case "storage":
			if len(parts) > 1 {
				displayStorageAt(executionContext, parts[1])
			} else {
				fmt.Println("Error: Missing key. Usage: storage <hex_key>")
			}

		case "debug":
			if len(parts) < 2 {
				fmt.Println("Error: Missing bytecode. Usage: debug <bytecode>")
				continue
			}
			debugBytecode(executionContext, parts[1])

		default:
			fmt.Println("Unknown command. Type 'help' for available commands")
		}
	}
}

func printHelp() {
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

func runBytecode(ctx *evm.ExecutionContext, hexString string) {
	// Remove 0x prefix if present
	if strings.HasPrefix(hexString, "0x") {
		hexString = hexString[2:]
	}

	// Decode the hex string to bytes
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

func pushValue(ctx *evm.ExecutionContext, hexValue string) {
	// Remove 0x prefix if present
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

	// Convert bytes to uint256
	value := uint256.NewInt(0)
	value.SetBytes(bytes)

	// Push to stack
	err = ctx.Stack.Push(value)
	if err != nil {
		fmt.Printf("Error pushing to stack: %v\n", err)
		return
	}

	fmt.Printf("Pushed 0x%s to stack\n", hexValue)
	fmt.Printf("Stack: %s\n", ctx.Stack.ToString())
}

// Function to display value at a specific storage key
func displayStorageAt(ctx *evm.ExecutionContext, hexKey string) {
	// Remove 0x prefix if present
	if strings.HasPrefix(hexKey, "0x") {
		hexKey = hexKey[2:]
	}

	// Add leading zeros if less than 64 characters
	for len(hexKey) < 64 {
		hexKey = "0" + hexKey
	}

	// Decode the hex key
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		fmt.Printf("Error decoding key: %v\n", err)
		return
	}

	// Convert to common.Hash
	keyHash := common.BytesToHash(keyBytes)

	// Get value from storage
	value := ctx.Storage.Sload(keyHash)

	if value == nil || len(value) == 0 {
		fmt.Printf("Storage[0x%s] = 0x0\n", hexKey)
	} else {
		fmt.Printf("Storage[0x%s] = 0x%s\n", hexKey, hex.EncodeToString(value))
	}
}

// Function to execute bytecode with debugging
func debugBytecode(ctx *evm.ExecutionContext, hexString string) {
	// Remove 0x prefix if present
	if strings.HasPrefix(hexString, "0x") {
		hexString = hexString[2:]
	}

	// Decode the hex string to bytes
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
		// Current opcode
		op := evm.GetOpcodeName(ctx.ByteCode[ctx.ProgramCounter])

		fmt.Printf("\nStep %d: PC=%d, Opcode=%s\n", step, ctx.ProgramCounter, op)
		fmt.Printf("Stack: %s\n", ctx.Stack.ToString())

		fmt.Print("Press ENTER to continue (or 'q' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
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
