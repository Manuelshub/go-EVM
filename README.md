# Go-EVM: Ethereum Virtual Machine Implementation in Go

A lightweight, educational implementation of the Ethereum Virtual Machine (EVM) written in Go (Golang). This project aims to provide a clear, understandable implementation of the EVM for learning and experimentation.

## Overview

Go-EVM is a from scratch implementation of the Ethereum Virtual Machine, the runtime environment for smart contracts on the Ethereum blockchain. It includes:

- A complete stack-based architecture
- Memory and storage management
- Gas metering and accounting
- Opcode implementation
- Contract execution
- Interactive debugging tools

## Features

- **Complete Opcode Support**: Implements the core EVM instruction set
- **Gas Metering**: Accurate gas calculation for operations
- **Interactive CLI**: Debug and step through contract execution
- **Storage system**: Persistent Key-value storage for contracts
- **Stack Manipulation**: Complete implementation of the 1024-element stack

## Architecture

Go-EVM is structured around several core components

## Stack

The EVM uses a stack-based archtecture for operations
```go
// Stack is a list of 32 bytes elements
type Stack struct {
    elem []*uint256.Int
}

// Push adds a new element to the stack
func (stack *Stack) Push(value *uint256.Int) error {
	if stack.Size() >= MAX_STACK_SIZE {
		return ErrStackOverflow
	}
	stack.elem = append(stack.elem, value)
	return nil
}

// Pop removes the last element from the stack
func (stack *Stack) Pop() (*uint256.Int, error) {
	if stack.Size() == 0 {
		return nil, ErrStackUnderflow
	}
	value := stack.elem[stack.Size()-1]
	stack.elem = stack.elem[:stack.Size()-1]
	return value, nil
}
```
## Memory
Linear, byte-addressable memory
```go
// Memory represents EVM memory
type Memory struct {
	data []byte
}

// Mstore is a method that writes the data to the memory at the given offset
func (mem *Memory) Mstore(offset uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	// Ensure memory is expanded to fit data
	memSlice := mem.expand(offset, uint64(len(data)))

	// Copy data to memory
	copy(memSlice, data)
}
```

## Interactive CLI

Go-EVM includes a REPL interface for interacting with the EVM:

```bash
(go-EVM) help
Available commands:
  help               - Show this help
  run <bytecode>     - Execute bytecode (hex format)
  debug <bytecode>   - Execute bytecode with step-by-step tracing
  stack              - Display current stack
  storage <key>      - Display storage value at key (hex format)
  push <value>       - Push a hex value onto the stack
  reset              - Reset the execution context
  exit, quit         - Exit the program
```

## Example Usage

Running a simple addition contract:

```bash
(go-EVM) run 0x6001600201
Running bytecode: 0x6001600201
Execution successful. No return value.
Stack: [0x3]
Gas used: 12
```

This executes:
- `PUSH1 0x01`: Push 1 onto the stack
- `PUSH1 0x02`: Push 2 onto the stack
- `ADD`: Add the top two values

## Getting Started

### Prerequisites

- Go 1.16 or higher
- `github.com/holiman/uint256` package
- `github.com/ethereum/go-ethereum/common` package

### Installation

```bash
git clone https://github.com/Manuelshub1/go-EVM.git
cd go-EVM
go build
```

### Running

```bash
./go-EVM
```

## Roadmap

- [x] Basic stack operations
- [x] Memory management
- [x] Storage operations
- [x] Control flow (jumps)
- [x] Basic arithmetic and logic
- [ ] Contract creation
- [ ] Message calling between contracts
- [ ] Complete environment operations
- [ ] Precompiled contracts
- [ ] Full compatibility with Ethereum tests

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Ethereum Yellow Paper
- Go Ethereum (geth) implementation
- Ethereum community and documentation

## Disclaimer

This is an educational project and not intended for production use. It may not implement all security checks or optimizations found in production EVMs.