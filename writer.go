package cdb64

import (
	"bufio"
	"os"
	"sync"
)

// Writer provides an API for creating a CDB database record by record.
//
// Close (or Freeze) must be called to finalize the database, or the resulting
// file will be invalid. Any errors during writing or finalizing are
// unrecoverable.
type Writer struct {
	file    *os.File
	entries [256][]entry

	buffer  *bufio.Writer
	offset  int64
	scratch []byte

	finalizeOnce sync.Once
}

type entry struct {
	hash   uint64
	offset int64
}

// Create opens a CDB database at the given path. If the file exists, it will
// be overwritten.
func Create(name string) (*Writer, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	// prepare temporary scratch buffer
	t := make([]byte, headerSize)

	// Skip 256 * 8 * 2 bytes for the index at the head of the file.
	if _, err := f.Write(t); err != nil {
		_ = f.Close()
		return nil, err
	}

	return &Writer{
		file:    f,
		buffer:  bufio.NewWriterSize(f, 65536),
		offset:  headerSize,
		scratch: t,
	}, nil
}

// Put adds a key/value pair to the database.
func (w *Writer) Put(key, value []byte) error {
	// Record the entry in the hash table, to be written out at the end.
	hash := hashKey(key)
	table := hash & 0xff

	w.entries[table] = append(w.entries[table], entry{
		hash:   hash,
		offset: w.offset,
	})

	// Write the key length, then value length, then key, then value.
	if err := w.writeTuple(uint64(len(key)), uint64(len(value))); err != nil {
		return err
	}
	w.offset += 16

	if _, err := w.buffer.Write(key); err != nil {
		return err
	}
	w.offset += int64(len(key))

	if _, err := w.buffer.Write(value); err != nil {
		return err
	}
	w.offset += int64(len(value))

	return nil
}

// Close finalizes the database, then closes it to further writes.
//
// Close or Freeze must be called to finalize the database, or the resulting
// file will be invalid.
func (w *Writer) Close() error {
	var err error
	w.finalizeOnce.Do(func() {
		err = w.finalize()
	})
	if e := w.file.Close(); e != nil {
		err = e
	}
	return err
}

// Freeze closes the writer, then opens it for reads.
func (w *Writer) Freeze() (*Reader, error) {
	if err := w.Close(); err != nil {
		return nil, err
	}
	return Open(w.file.Name())
}

func (w *Writer) finalize() error {
	var header header

	// Write the hashtables out, one by one, at the end of the file.
	for i := 0; i < 256; i++ {
		tableEntries := w.entries[i]
		tableSize := len(tableEntries) << 1

		header[i] = table{
			offset: w.offset,
			length: tableSize,
		}

		sorted := make([]entry, tableSize)
		for _, entry := range tableEntries {
			slot := hashSlot(entry.hash, tableSize)
			for {
				if sorted[slot].hash == 0 {
					sorted[slot] = entry
					break
				}
				slot = (slot + 1) % tableSize
			}
		}

		for _, entry := range sorted {
			if err := w.writeTuple(entry.hash, uint64(entry.offset)); err != nil {
				return err
			}
			w.offset += 16
		}
	}

	// We're done with the buffer.
	if err := w.buffer.Flush(); err != nil {
		return err
	}

	// Write out the header.
	w.scratch = w.scratch[:headerSize]
	for i, table := range header {
		off := i * 16
		binLE.PutUint64(w.scratch[off:off+8], uint64(table.offset))
		binLE.PutUint64(w.scratch[off+8:off+16], uint64(table.length))
	}

	_, err := w.file.WriteAt(w.scratch, 0)
	return err
}

func (w *Writer) writeTuple(first, second uint64) error {
	w.scratch = w.scratch[:16]
	binLE.PutUint64(w.scratch[:8], first)
	binLE.PutUint64(w.scratch[8:], second)

	_, err := w.buffer.Write(w.scratch)
	return err
}
