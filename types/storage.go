package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type Storage struct {
	elem map[common.Hash][]byte
}

func NewStorage() *Storage {
	return &Storage{
		elem: make(map[common.Hash][]byte),
	}
}

func (s *Storage) Sload(key common.Hash) []byte {
	value, ok := s.elem[key]
	if !ok {
		return nil
	}
	return value
}

func (s *Storage) Sstore(key common.Hash, value []byte) {
	if value == nil {
		return
	}
	s.elem[key] = value
}
