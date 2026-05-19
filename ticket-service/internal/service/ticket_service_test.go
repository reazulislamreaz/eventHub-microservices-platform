package service

import (
	"errors"
	"testing"
)

func TestIsDuplicateKey(t *testing.T) {
	if !isDuplicateKey(errors.New("ERROR: duplicate key value violates unique constraint")) {
		t.Fatal("expected duplicate key detection")
	}
	if isDuplicateKey(errors.New("other error")) {
		t.Fatal("expected false for unrelated error")
	}
}
