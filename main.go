package main

import (
    "log"
	"fmt"
	"bytes"
	"io"
    "time"
	"github.com/Ridwan-Al-Mahmud/GoDistributedStorage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts {
		ListenAddr:    listenAddr,
     	HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts {
		EncKey:            newEncryptionKey(),
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", "")
	s3 := makeServer(":5000", ":3000", ":4000")
	
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(time.Millisecond * 500)
	go func() {
		log.Fatal(s2.Start())
	}()
	time.Sleep(time.Millisecond * 500)
	go func() {
		log.Fatal(s3.Start())
	}()

	time.Sleep(2 * time.Second)
	
	key := "my_cool_picture.jpg"
	data := bytes.NewReader([]byte("my big data file here!"))
	if err := s3.Store(key, data); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	
	r, err := s3.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	if err := s3.store.Delete(key); err != nil {
		log.Fatal(err)
	}
}