package cdb64

// Iterator represents a sequential iterator over a CDB database.
type Iterator struct {
	reader *Reader
	off    int64
	maxOff int64
	err    error
	buf    []byte
	piv    int
}

// Next reads the next key/value pair and advances the iterator one record.
// It returns false when the scan stops, either by reaching the end of the
// database or an error. After Next returns false, the Err method will return
// any error that occurred while iterating.
func (it *Iterator) Next() bool {
	if it.off >= it.maxOff || it.err != nil {
		return false
	}

	klen, vlen, err := it.reader.readTuple(it.off)
	if err != nil {
		it.err = err
		return false
	}

	if n := int(klen + vlen); n <= cap(it.buf) {
		it.buf = it.buf[:n]
	} else {
		it.buf = make([]byte, n)
	}

	_, err = it.reader.reader.ReadAt(it.buf, it.off+16)
	if err != nil {
		it.err = err
		return false
	}

	// Update iterator state
	it.piv = int(klen)
	it.off += 16 + int64(klen+vlen)

	return true
}

// Key returns the current key, which is valid until a subsequent
// call to Next(). You must copy they key if you plan to use it
// beyond this point.
func (it *Iterator) Key() []byte { return it.buf[:it.piv] }

// Value returns the current value, which is valid until a subsequent
// call to Next(). You must copy they value if you plan to use it
// beyond this point.
func (it *Iterator) Value() []byte { return it.buf[it.piv:] }

// Err returns the current error.
func (it *Iterator) Err() error { return it.err }
