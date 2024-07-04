package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"errors"
	"bytes"
	"strings"
	"crypto/sha1"
	"encoding/hex"
)

const (
	defaultRootFolderName = "alnetwork"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i * blockSize, (i * blockSize) + blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		Pathname: strings.Join(paths, "/"),
		FilePath: hashStr,
	}
}

type PathKey struct {
	Pathname string
	FilePath string
}

func (p PathKey) FirstPathName() string{
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.FilePath)
}

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey{
	return PathKey {
		Pathname: key,
		FilePath: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
    }
	return &Store {
		StoreOpts: opts,
	}
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(pathKeyWithRoot)
}
func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.WriteStream(key, r)
}
func (s *Store) WriteStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}
	fullPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	f, err := os.Create(fullPathNameWithRoot)
    if err != nil {
		return 0, err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}
	
	return n, nil
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func(){
		log.Printf("deleted (%s) from disk", pathKey.Pathname)
	}()
    firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}