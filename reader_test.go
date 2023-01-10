package cdb64_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	. "github.com/bsm/cdb64"
)

func TestReader(t *testing.T) {
	dir, err := os.MkdirTemp("", "cdb64_test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.RemoveAll(dir)

	// seed
	w, err := Create(dir + "/test.cdb")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer w.Close()

	if err := seedData(w, 1000); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// seed some more exotic entries
	if err := w.Put(nil, []byte("blank")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := w.Put([]byte("blank"), nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := w.Put([]byte("key-00000333"), []byte("duplicate")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// open reader
	r, err := w.Freeze()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer r.Close()

	t.Run("Get", func(t *testing.T) {
		if v, err := r.Get(nil); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if exp := []byte("blank"); !bytes.Equal(v, exp) {
			t.Errorf("expected %q, got %q", exp, v)
		}

		if v, err := r.Get([]byte{}); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if exp := []byte("blank"); !bytes.Equal(v, exp) {
			t.Errorf("expected %q, got %q", exp, v)
		}

		if v, err := r.Get([]byte("missing")); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if v != nil {
			t.Errorf("expected %v, got %v", nil, v)
		}

		if v, err := r.Get([]byte("blank")); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if exp := []byte{}; !bytes.Equal(v, exp) {
			t.Errorf("expected %q, got %q", exp, v)
		}

		if v, err := r.Get([]byte("key-00000005")); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if exp := []byte("val-00000005"); !bytes.Equal(v, exp) {
			t.Errorf("expected %q, got %q", exp, v)
		}

		if v, err := r.Get([]byte("key-00000333")); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if exp := []byte("val-00000333"); !bytes.Equal(v, exp) {
			t.Errorf("expected %q, got %q", exp, v)
		}
	})

	t.Run("Batch", func(t *testing.T) {
		b := r.Batch()

		if v, err := b.Get([]byte("missing")); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if v != nil {
			t.Errorf("expected %v, got %v", nil, v)
		}

		var res2 []byte
		if v, err := b.Get([]byte("key-00000005")); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if exp := []byte("val-00000005"); !bytes.Equal(v, exp) {
			t.Errorf("expected %q, got %q", exp, v)
		} else {
			res2 = v
		}

		if v, err := b.Get([]byte("key-00000333")); err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if exp := []byte("val-00000333"); !bytes.Equal(v, exp) {
			t.Errorf("expected %q, got %q", exp, v)
		}

		// res2 changes after 3rd call to Get
		if exp := []byte("val-00000333"); !bytes.Equal(res2, exp) {
			t.Errorf("expected %q, got %q", exp, res2)
		}
	})

	t.Run("Iterator", func(t *testing.T) {
		iter := r.Iterator()

		var keys, vals []string
		for iter.Next() {
			keys = append(keys, string(iter.Key()))
			vals = append(vals, string(iter.Value()))
		}
		if exp, got := 1003, len(keys); exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if exp, got := 1003, len(vals); exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}

		if exp, got := []string{"key-00000001", "key-00000003", "key-00000005"}, keys[:3]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if exp, got := []string{"val-00000001", "val-00000003", "val-00000005"}, vals[:3]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if exp, got := []string{"key-00000331", "key-00000333", "key-00000335"}, keys[165:168]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if exp, got := []string{"val-00000331", "val-00000333", "val-00000335"}, vals[165:168]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if exp, got := []string{"", "blank", "key-00000333"}, keys[1000:]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if exp, got := []string{"blank", "", "duplicate"}, vals[1000:]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
}
