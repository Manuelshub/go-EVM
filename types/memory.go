package types

import (
	"encoding/hex"
)

type Memory struct {
	data []byte
}

// NewMemory is a factory method that returns a new Memory object
func NewMemory() *Memory {
	return &Memory{
		data: make([]byte, 0),
	}
}

// expand is a helper method that expands the memory to the required size if needed.
// It returns a slice of memory from offset to offset+size. If size is 0, returns nil.
// The returned slice shares the underlying array with mem.data.
func (mem *Memory) expand(offset, size uint64) []byte {
	if size == 0 {
		return nil
	}
	requiredSize := offset + size

	if len(mem.data) < int(requiredSize) {
		extra := int(requiredSize) - len(mem.data)
		mem.data = append(mem.data, make([]byte, extra)...)
	}
	return mem.data[offset : offset+size]
}

// Expand is a public helper method that expands the memory to the required size if needed.
// It returns a slice of memory from offset to offset+size. If size is 0, returns nil.
func (mem *Memory) Expand(offset, size uint64) []byte {
	return mem.expand(offset, size)
}

// Mload is a method that reads the memory at the given offset and returns
// the slice of memory at that offset
func (mem *Memory) Mload(offset uint64) []byte {
	return mem.expand(offset, 32)
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

// MstoreByte stores a single byte at the specified offset
func (mem *Memory) MstoreByte(offset uint64, data byte) {
	memSlice := mem.expand(offset, 1)
	memSlice[0] = data
}

// Size returns the current size of the memory in bytes
func (mem *Memory) Size() uint64 {
	return uint64(len(mem.data))
}

func (mem *Memory) ToString() string {
	var m string
	if len(mem.data) == 0 {
		m = "[]"
		return m
	}
	m += "["
	for _, b := range mem.data {
		m += hex.EncodeToString([]byte{b})
	}
	m += "]"
	return m
}