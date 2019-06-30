# CDB64

[![Build Status](https://travis-ci.org/bsm/dbx.png?branch=master)](https://travis-ci.org/bsm/dbx)
[![GoDoc](https://godoc.org/github.com/bsm/dbx?status.png)](http://godoc.org/github.com/bsm/dbx)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsm/dbx)](https://goreportcard.com/report/github.com/bsm/dbx)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is a native Go implementation of [cdb][1], a constant key/value database
with some very nice properties, but without the 4GB size limit.

Adapted from the original [design doc][1]:

> cdb is a fast, reliable, simple package for creating and reading constant databases. Its database structure provides several features:
>
> - Fast lookups: A successful lookup in a large database normally takes just two disk accesses. An unsuccessful lookup takes only one.
> - Low overhead: A database uses 4096 bytes, plus 32 bytes per record, plus the space for keys and data.
> - No random limits: cdb can handle any database up to 16 exabytes. There are no other restrictions; records don't even have to fit into memory. Databases are stored in a machine-independent format.

[1]: http://cr.yp.to/cdb.html

This repo is based on github.com/chrislusf/cdb64

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/bsm/cdb64"
)

func main() {
	w, err := cdb64.Create("/tmp/cdb64-example.cdb")
	if err != nil {
		log.Fatalln(err)
	}
	defer w.Close()

	// Write some key/value pairs.
	_ = w.Put([]byte("Alice"), []byte("Hoax"))
	_ = w.Put([]byte("Bob"), []byte("Hope"))
	_ = w.Put([]byte("Charlie"), []byte("Horse"))

	// Freeze and re-open it for reading.
	db, err := w.Freeze()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Fetch a value.
	v, err := db.Get([]byte("Alice"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(string(v))
}
```
