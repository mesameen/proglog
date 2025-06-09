package applog

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

var (
	lenWidth = 8
)

type store struct {
	*os.File
	size uint64
	buf  *bufio.Writer
	mu   sync.RWMutex
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		log.Printf("Failed to get stat of file %s. Error: %v\n", f.Name(), err)
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

func (s *store) Append(p []byte) (uint64, uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos := s.size
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	n, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	n += lenWidth
	s.size += uint64(n)
	return uint64(n), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	data := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(data, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(data))
	if _, err := s.File.ReadAt(b, int64(pos+uint64(lenWidth))); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
