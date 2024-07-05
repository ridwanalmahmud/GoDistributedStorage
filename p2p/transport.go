package p2p

import (
	"net"
)

// Peer is an interface that represents the remote node
type Peer interface{
	Send([]byte) error
	CloseStream()
	net.Conn
}

// Transport is anything that handles the communication between the nodes
type Transport interface{
	Addr()            string
	ListenAndAccept() error
	Consume()         <-chan RPC
	Close()           error   
	Dial(string)      error
}