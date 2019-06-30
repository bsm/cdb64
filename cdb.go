/*
Package cdb64 provides a native implementation of cdb, a fast constant key/value
database, but without the 4GB size limitation.

For more information on cdb, see the original design doc at http://cr.yp.to/cdb.html.

This is based on https://github.com/chrislusf/cdb64 which in itself is based on
https://github.com/colinmarc/cdb.
*/
package cdb64

import (
	"encoding/binary"

	"github.com/dgryski/go-farm"
)

const headerSize = 256 * 8 * 2

type header [256]table

type table struct {
	offset int64
	length int
}

var binLE = binary.LittleEndian

func hashKey(p []byte) uint64 {
	return farm.Hash64(p)
}

func hashSlot(hash uint64, size int) int {
	return int(hash>>8) % size
}
