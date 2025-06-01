package server

import (
	"fmt"
	"sync"
)

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")

type Log struct {
	mu      sync.RWMutex
	records []*Record
}

func NewLog() *Log {
	return &Log{
		records: make([]*Record, 0),
	}
}

func (c *Log) Append(record *Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil
}

func (c *Log) Read(offset uint64) (*Record, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if offset >= uint64(len(c.records)) {
		return nil, ErrOffsetNotFound
	}
	return c.records[offset], nil
}
