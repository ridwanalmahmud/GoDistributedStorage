package main

import (
	"bytes"
	"testing"
	"io/ioutil"
	// "fmt"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	ExpectedFileName := "6804429f74181a63c50c3d81d733a12f14a353ff"
	ExpectedPathName := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if pathKey.Pathname != ExpectedPathName {
		t.Errorf("have %s, want %s", pathKey.Pathname, ExpectedPathName)
	}
	if pathKey.FilePath != ExpectedFileName {
		t.Errorf("have %s, want %s", pathKey.FilePath, ExpectedFileName)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)
	//for i := 0; i < 10; i++ {
	//key := fmt.Sprintf("foo_%d", i)
	key := "foo"
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
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
	if ok := s.Has(key); ok {
		t.Errorf("expected to NOT have key %s", key)
	}
	//}
}

func newStore() *Store {
	opts := StoreOpts {
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}