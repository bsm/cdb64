package cdb64

// Batch instances allow to query the Reader more efficiently
// with a sequence of Get requests. Unlike Readers, Batches
// are not thread-safe and must not be shared across goroutines.
type Batch struct {
	reader *Reader
	buf    []byte
}

// Get returns the value for a given key, or nil if it can't be found.
// Returned values are re-used throughout the life-time of the Batch and
// are therefore only valid until the next call to Get.
func (b *Batch) Get(key []byte) ([]byte, error) {
	val, buf, err := b.reader.find(key, b.buf)
	b.buf = buf
	return val, err
}
