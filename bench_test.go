package cdb64_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/bsm/cdb64"
)

func BenchmarkReader_Get(b *testing.B) {
	benchmarkGet(b, func(b *testing.B, r *cdb64.Reader, keys [][]byte) {
		for i := 0; i < b.N; i++ {
			_, err := r.Get(keys[i%len(keys)])
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkBatch_Get(b *testing.B) {
	benchmarkGet(b, func(b *testing.B, r *cdb64.Reader, keys [][]byte) {
		batch := r.Batch()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := batch.Get(keys[i%len(keys)])
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkGet(b *testing.B, fn func(*testing.B, *cdb64.Reader, [][]byte)) {
	const numKeys = 100000

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

	val := bytes.Repeat([]byte{'x'}, 512)
	for i := 0; i < numKeys; i++ {
		if err := w.Put(seedKey(i), val); err != nil {
			b.Fatal(err)
		}
	}

	r, err := w.Freeze()
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()

	// seed keys to query with 80% hit rate
	rnd := rand.New(rand.NewSource(3))
	keys := make([][]byte, 0, numKeys/4)
	for i := 0; i < numKeys/4; i++ {
		keys = append(keys, seedKey(rnd.Intn(numKeys/4*5)))
	}

	b.ResetTimer()
	fn(b, r, keys)
}
