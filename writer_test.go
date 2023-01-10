package cdb64_test

import (
	"os"
	"testing"

	. "github.com/bsm/cdb64"
)

func TestWriter(t *testing.T) {
	dir, err := os.MkdirTemp("", "cdb64_test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.RemoveAll(dir)

	w, err := Create(dir + "/test.cdb")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer w.Close()

	if err := seedData(w, 1000); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
