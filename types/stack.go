package types

import (
	"errors"

	"github.com/holiman/uint256"
)

const MAX_STACK_SIZE = 1024

var (
	ErrStackOverflow  = errors.New("Stack Overflow")  // ErrStackOverflow is returned when the stack is full
	ErrStackUnderflow = errors.New("Stack Underflow") // ErrStackUnderflow is returned when the stack is empty
)

// Stack is a list of 32 bytes elements
type Stack struct {
	elem []*uint256.Int // A list of 32 bytes elements
}

// var stack = Stack{
// 	elem: make([][32]byte, MAX_STACK_SIZE),
// }

func NewStack() *Stack {
	return &Stack{
		elem: make([]*uint256.Int, 0, MAX_STACK_SIZE), // Initialize the stack with a capacity of MAX_STACK_SIZE
	}
}

// // ToBigInt returns the n-th element of the stack as a big.Int
// func (stack *Stack) ToBigInt(n int) *big.Int {
// 	if len(stack.elem) == 0 {
// 		fmt.Fprintln(os.Stderr, ErrStackUnderflow.Error())
// 		return nil
// 	}
// 	return new(big.Int).SetBytes(stack.elem[len(stack.elem)-n])
// }

// Push adds a new element to the stack
func (stack *Stack) Push(value *uint256.Int) error {
	if len(stack.elem) >= MAX_STACK_SIZE {
		return ErrStackOverflow
	}
	stack.elem = append(stack.elem, value)
	return nil
}

// Pop removes the last element from the stack
func (stack *Stack) Pop() (*uint256.Int, error) {
	if len(stack.elem) == 0 {
		return nil, ErrStackUnderflow
	}
	value := stack.elem[len(stack.elem)-1]
	stack.elem = stack.elem[:len(stack.elem)-1]
	return value, nil
}

// Peek returns the last element from the stack without removing it
func (stack *Stack) Peek() (*uint256.Int, error) {
	if len(stack.elem) == 0 {
		return nil, ErrStackUnderflow
	}
	return stack.elem[len(stack.elem)-1], nil

}

// Swap swaps the n-th element from the top of the stack with the top element
func (stack *Stack) Swap(n int) error {
	if len(stack.elem) < 2 {
		return ErrStackUnderflow
	}
	stack.elem[len(stack.elem)-1], stack.elem[len(stack.elem)-1-n] = stack.elem[len(stack.elem)-1-n], stack.elem[len(stack.elem)-1]
	return nil
}

// Dup duplicates the n-th element from the top of the stack
func (stack *Stack) Dup(n int) error {
	if len(stack.elem) < 1 {
		return ErrStackUnderflow
	}
	stack.elem = append(stack.elem, stack.elem[len(stack.elem)-n])
	return nil
}

// ToString method returns the string representation of the stack
func (stack *Stack) ToString() string {
	var s string

	if len(stack.elem) == 0 {
		s = "[]"
		return s
	}
	s = "["
	for i := len(stack.elem) - 1; i >= 0; i-- {
		s += stack.elem[i].Hex() // Use Hex() to get the hex representation of the element
		if i != 0 {
			s += ", "
		}
	}
	s += "]"
	return s
}
