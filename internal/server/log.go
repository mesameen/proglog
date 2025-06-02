package server

import (
	"sync"

	"github.com/mesameen/proglog/internal/model"
)

type Log struct {
	mu      sync.RWMutex
	records []*model.Record
}

func NewLog() *Log {
	return &Log{
		records: make([]*model.Record, 0),
	}
}

func (c *Log) Append(record *model.Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil
}

func (c *Log) Read(offset uint64) (*model.Record, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if offset >= uint64(len(c.records)) {
		return nil, model.ErrOffsetNotFound
	}
	return c.records[offset], nil
}
