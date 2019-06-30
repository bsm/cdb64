package cdb64_test

import (
	"fmt"
	"log"

	"github.com/bsm/cdb64"
)

func Example() {
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
	// Output: Hoax
}
