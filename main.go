package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Manuelshub/go-EVM/evm"
	h "github.com/Manuelshub/go-EVM/helpers"
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
			h.PrintHelp()

		case "run":
			if len(parts) < 2 {
				fmt.Println("Error: Missing bytecode. Usage: run <bytecode>")
				continue
			}
			h.RunBytecode(executionContext, parts[1])

		case "stack":
			fmt.Println(executionContext.Stack.ToString())

		case "push":
			if len(parts) < 2 {
				fmt.Println("Error: Missing value. Usage: push <hex_value>")
				continue
			}
			h.PushValue(executionContext, parts[1])

		case "debug":
			if len(parts) < 2 {
				fmt.Println("Error: Missing value. Usage: debug <bytecode>")
				continue
			}
			h.DebugBytecode(executionContext, parts[1])

		case "reset":
			executionContext = evm.NewExecutionContext()
			fmt.Println("Execution context reset")

		default:
			fmt.Println("Unknown command. Type 'help' for available commands")
		}
	}
}
