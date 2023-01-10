package cdb64_test

import (
	"fmt"
)

func seedKey(n int) []byte {
	return []byte(fmt.Sprintf("key-%08d", n))
}

func seedData(w interface{ Put([]byte, []byte) error }, n int) error {
	var key, val []byte
	for i := 0; i < n*2; i += 2 {
		key = append(key[:0], seedKey(i+1)...)
		val = append(val[:0], fmt.Sprintf("val-%08d", i+1)...)
		if err := w.Put(key, val); err != nil {
			return err
		}
	}
	return nil
}
