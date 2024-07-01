package main

import (
	"bytes"
	"testing"
)

func StoreTest(t *testing.T) {
	opts := StoreOpts {
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg files"))
	if err := s.writeStream("my special picture", data); err != nil {
		t.Error(err)
	}
}