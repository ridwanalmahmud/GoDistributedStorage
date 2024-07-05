package main

import (
    "log"
	"fmt"
	//"bytes"
	"io/ioutil"
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
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(2 * time.Second)
	go s2.Start()
	time.Sleep(2 * time.Second)
	/*data := bytes.NewReader([]byte("my big data file here"))
	s2.Store("my_cool_picture.jpg", data)
	time.Sleep(5 * time.Millisecond)*/
	

    r, err := s2.Get("my_cool_picture.jpg")
	//r, err := s2.Get("don't have key")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}