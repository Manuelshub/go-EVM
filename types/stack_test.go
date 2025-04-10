package types

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	all_tests := []struct {
		name   string
		op     func (s *Stack) any
		want   any
		wantErr bool
	}{
		{
			name: "Push",
			op: func(s *Stack) any {
				return s.Push(uint256.NewInt(1))
			},
			want: nil,
			wantErr: false,
		},
		{
			name: "Push and Pop",
			op: func(s *Stack) any {
				s.Push(uint256.NewInt(21))
				data, _ := s.Pop()
				return data.Uint64()
			},
			want: 21,
			wantErr: false,
		},
		{
			name: "Push and Peek",
			op: func(s *Stack) any {
				s.Push(uint256.NewInt(42))
				data, _ := s.Peek()
				return data.Uint64()
			},
			want: 42,
			wantErr: false,
		},
		{
			name: "Push and Stack Overflow",
			op: func(s *Stack) any {
				for i := 0; i < MAX_STACK_SIZE; i++ {
					s.Push(uint256.NewInt(uint64(i)))
				}
				return s.Push(uint256.NewInt(100))
			},
			want: ErrStackOverflow,
			wantErr: true,
		},
		{
			name: "Pop and Stack Underflow",
			op: func(s *Stack) any {
				_, err := s.Pop()
				return err
			},
			want: ErrStackUnderflow,
			wantErr: true,
		},
		{
			name: "Swap",
			op: func(s *Stack) any {
				s.Push(uint256.NewInt(32))
				s.Push(uint256.NewInt(64))
				err := s.Swap(1)
				return err
			},
			want: nil,
			wantErr: false,
		},
		{
			name: "Swap and Stack Underflow",
			op: func(s *Stack) any {
				s.Push(uint256.NewInt(32))
				err := s.Swap(1)
				return err
			},
			want: ErrStackUnderflow,
			wantErr: true,
		},
		{
			name: "Duplicate",
			op: func(s *Stack) any {
				s.Push(uint256.NewInt(32))
				err := s.Dup(1)
				return err
			},
			want: nil,
			wantErr: false,
		},
		{
			name: "Duplicate and Stack Underflow",
			op: func(s *Stack) any {
				s.Push(uint256.NewInt(32))
				err := s.Dup(2)
				return err
			},
			want: ErrStackUnderflow,
			wantErr: true,
		},
	}
	for _, tt := range all_tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewStack()
			result := tt.op(stack)
			if tt.wantErr {
				assert.Equal(t, tt.want, result)
			}
		})
	}
}