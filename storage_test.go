package main

import (
	"bytes"
	"testing"
	"io/ioutil"
	// "fmt"
)

/*func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathName := CASPathTransformFunc(key)
	ExpectedPathName := ""
	if pathName != ExpectedPathName {
		t.Errorf("have %s, want %s", pathName, ExpectedPathName)
	}
}*/

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts {
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momsspecials"
	data := []byte("some jpg files")
	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts {
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momsspecials"
	data := []byte("some jpg files")
	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}
	r, err := s.Read(key)
	if err != nil{
		t.Error(err)
	}
	b, _ := ioutil.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}
	s.Delete(key)
}