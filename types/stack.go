package types

import (
	"errors"
	"fmt"
	"os"
)

const MAX_STACK_SIZE = 1024

var (
	ErrStackOverflow  = errors.New("Stack Overflow")  // ErrStackOverflow is returned when the stack is full
	ErrStackUnderflow = errors.New("Stack Underflow") // ErrStackUnderflow is returned when the stack is empty
)

// Stack is a list of 32 bytes elements
type Stack struct {
	elem [][32]byte // A list of 32 bytes elements
}

// var stack = Stack{
// 	elem: make([][32]byte, MAX_STACK_SIZE),
// }

func NewStack() *Stack {
	return &Stack{
		elem: make([][32]byte, MAX_STACK_SIZE),
	}
}

func (stack *Stack) Push(value [32]byte) {
	if len(stack.elem) > MAX_STACK_SIZE {
		fmt.Fprintln(os.Stderr, ErrStackOverflow.Error())
		return
	}
	stack.elem = append(stack.elem, value)
}

func (stack *Stack) Pop() [32]byte {
	if len(stack.elem) == 0 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return [32]byte{}
	}
	value := stack.elem[len(stack.elem)-1]
	stack.elem = stack.elem[:len(stack.elem)-1]
	return value
}

// func Peek() {}

func (stack *Stack) Swap() {
	if len(stack.elem) < 2 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return
	}
	stack.elem[len(stack.elem)-1], stack.elem[len(stack.elem)-2] = stack.elem[len(stack.elem)-2], stack.elem[len(stack.elem)-1]
}

func (stack *Stack) Dup() {
	if len(stack.elem) < 1 {
		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
		return
	}
	stack.elem = append(stack.elem, stack.elem[len(stack.elem)-1])
}
