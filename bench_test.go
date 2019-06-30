package cdb64_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/bsm/cdb64"
)

func BenchmarkGet(b *testing.B) {
	dir, err := ioutil.TempDir("", "cdb64-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	w, err := cdb64.Create(dir + "/test.cdb")
	if err != nil {
		b.Fatal(err)
	}
	defer w.Close()

	val := bytes.Repeat([]byte{'x'}, 2048)
	for i := 0; i < 100000; i++ {
		if err := w.Put(seedKey(i), val); err != nil {
			b.Fatal(err)
		}
	}

	r, err := w.Freeze()
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()

	rnd := rand.New(rand.NewSource(3))
	keys := make([][]byte, 0, 10000)
	for i := 0; i < 10000; i++ {
		keys = append(keys, seedKey(rnd.Intn(100000)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.Get(keys[i%len(keys)])
		if err != nil {
			b.Fatal(err)
		}
	}
}
