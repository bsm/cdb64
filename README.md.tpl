# CDB64

[![Build Status](https://travis-ci.org/bsm/cdb64.png?branch=master)](https://travis-ci.org/bsm/cdb64)
[![GoDoc](https://godoc.org/github.com/bsm/cdb64?status.png)](http://godoc.org/github.com/bsm/cdb64)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsm/cdb64)](https://goreportcard.com/report/github.com/bsm/cdb64)
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

func main() {{ "Example" | code }}
```
