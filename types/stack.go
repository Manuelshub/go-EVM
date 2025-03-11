package types

import (
	"errors"
	"fmt"
	"math/big"
	"os"
)

const MAX_STACK_SIZE = 1024

var (
	ErrStackOverflow  = errors.New("Stack Overflow")  // ErrStackOverflow is returned when the stack is full
	ErrStackUnderflow = errors.New("Stack Underflow") // ErrStackUnderflow is returned when the stack is empty
)

// Stack is a list of 32 bytes elements
type Stack struct {
	elem [][]byte // A list of 32 bytes elements
}

// var stack = Stack{
// 	elem: make([][32]byte, MAX_STACK_SIZE),
// }

func NewStack() *Stack {
	return &Stack{
		elem: make([][]byte, MAX_STACK_SIZE),
	}
}

// ToBigInt returns the n-th element of the stack as a big.Int
func (stack *Stack) ToBigInt(n int) *big.Int {
	if len(stack.elem) == 0 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return nil
	}
	return new(big.Int).SetBytes(stack.elem[len(stack.elem)-n])
}

// Push adds a new element to the stack
func (stack *Stack) Push(value []byte) {
	if len(stack.elem) >= MAX_STACK_SIZE {
		fmt.Fprintln(os.Stderr, ErrStackOverflow.Error())
		return
	}
	stack.elem = append(stack.elem, value)
}

// Pop removes the last element from the stack
func (stack *Stack) Pop() []byte {
	if len(stack.elem) == 0 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return []byte{}
	}
	value := stack.elem[len(stack.elem)-1]
	stack.elem = stack.elem[:len(stack.elem)-1]
	return value
}

// Peek returns the last element from the stack without removing it
func (stack *Stack) Peek() []byte {
	if len(stack.elem) == 0 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return []byte{}
	}
	return stack.elem[len(stack.elem)-1]

}

// Swap swaps the n-th element from the top of the stack with the top element
func (stack *Stack) Swap(n int) {
	if len(stack.elem) < 2 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return
	}
	stack.elem[len(stack.elem)-1], stack.elem[len(stack.elem)-1-n] = stack.elem[len(stack.elem)-1-n], stack.elem[len(stack.elem)-1]
}

// Dup duplicates the n-th element from the top of the stack
func (stack *Stack) Dup(n int) {
	if len(stack.elem) < 1 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return
	}
	stack.elem = append(stack.elem, stack.elem[len(stack.elem)-n])
}
