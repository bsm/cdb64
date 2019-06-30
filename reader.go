package cdb64

import (
	"bytes"
	"os"
)

// Reader represents a thread-safe CDB database reader. To
// create a database, use Writer.
type Reader struct {
	file   *os.File
	header header
}

// Open opens an existing CDB database at the given path for reading.
func Open(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// read header
	buf := make([]byte, headerSize)
	if _, err := f.ReadAt(buf, 0); err != nil {
		_ = f.Close()
		return nil, err
	}

	// parse header
	var header header
	for i := 0; i < 256; i++ {
		off := i * 16
		header[i] = table{
			offset: int64(binLE.Uint64(buf[off : off+8])),
			length: int(binLE.Uint64(buf[off+8 : off+16])),
		}
	}

	return &Reader{
		file:   f,
		header: header,
	}, nil
}

// Get returns the value for a given key, or nil if it can't be found.
func (r *Reader) Get(key []byte) ([]byte, error) {
	value, _, err := r.find(key, nil)
	return value, err
}

// Batch creates a new batch operator. Please note that Batch instances are not
// thread-safe and must not be shared across goroutines.
func (r *Reader) Batch() *Batch {
	return &Batch{reader: r}
}

// Iterator creates an Iterator that can be used to iterate the database.
func (r *Reader) Iterator() *Iterator {
	return &Iterator{
		reader: r,
		off:    headerSize,
		maxOff: r.header[0].offset,
	}
}

// Close closes the reader to further reads.
func (r *Reader) Close() error {
	return r.file.Close()
}

func (r *Reader) find(key, buf []byte) ([]byte, []byte, error) {
	hash := hashKey(key)
	table := r.header[hash&0xff]
	if table.length == 0 {
		return nil, buf, nil
	}

	// Probe the given hash table, starting at the given slot.
	var value []byte
	firstSlot := hashSlot(hash, table.length)
	for slot := firstSlot; true; {
		slotOffset := table.offset + int64(16*slot)
		slotHash, offset, err := r.readTuple(slotOffset)
		if err != nil {
			return nil, buf, err
		}

		// An empty slot means the key doesn't exist.
		if slotHash == 0 {
			break
		}

		if slotHash == hash {
			value, buf, err = r.valueAt(int64(offset), key, buf)
			if err != nil {
				return nil, buf, err
			} else if value != nil {
				return value, buf, nil
			}
		}

		// advance to next slot
		if slot = (slot + 1) % table.length; slot == firstSlot {
			break
		}
	}

	return nil, buf, nil
}

func (r *Reader) valueAt(offset int64, key, buf []byte) ([]byte, []byte, error) {
	klen, vlen, err := r.readTuple(offset)
	if err != nil {
		return nil, buf, err
	}

	// We can compare key lengths before reading the key at all.
	if int(klen) != len(key) {
		return nil, buf, nil
	}

	if n := int(klen + vlen); n <= cap(buf) {
		buf = buf[:n]
	} else {
		buf = make([]byte, n)
	}
	if _, err := r.file.ReadAt(buf, offset+16); err != nil {
		return nil, buf, err
	}

	// If they keys don't match, this isn't it.
	if !bytes.Equal(buf[:klen], key) {
		return nil, buf, nil
	}

	return buf[klen:], buf, nil
}

func (r *Reader) readTuple(offset int64) (uint64, uint64, error) {
	buf := make([]byte, 16)
	if _, err := r.file.ReadAt(buf, offset); err != nil {
		return 0, 0, err
	}
	return binLE.Uint64(buf[:8]), binLE.Uint64(buf[8:]), nil
}
