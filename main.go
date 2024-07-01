package main

import (
	"fmt"
	"log"
	"github.com/Ridwan-Al-Mahmud/GoDistributedStorage/p2p"
)

func OnPeer(peer p2p.Peer) error {
	return fmt.Errorf("failed to OnPeer func")
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
    }     
	
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

    go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()
	
	select{}
}
