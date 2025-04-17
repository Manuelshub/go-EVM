package types

import "github.com/ethereum/go-ethereum/common"

type TransientStorage struct {
	data map[common.Hash][]byte
}

func NewTransientStorage() *TransientStorage {
	return &TransientStorage{
		data: make(map[common.Hash][]byte), // Initialize the map
	}
}

func (ts *TransientStorage) Tload(key common.Hash) []byte {
	value, ok := ts.data[key]
	if !ok {
		return nil
	}
	return value
}

func (ts *TransientStorage) Tstore(key common.Hash, value []byte) {
	if value == nil {
		return
	}
	ts.data[key] = value
}
