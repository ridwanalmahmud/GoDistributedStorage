package p2p

import (
	"fmt"
	"net"
	// "bytes"
	// "sync"
)

type TCPPeer struct {
	conn     net.Conn
	// if we dial and retrieve a conn -> outbound == true
	// if we accept and retrieve a conn -> inbound == false
	outbound bool
}

// Close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc handshakeFunc
	Decoder       Decoder
	OnPeer        func (Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener	
	rpcch         chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport {
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// consume implements the transport interface which returns read only channel for reading the incoming message from another peer

func (t *TCPTransport) Consume() <- chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil{
		return err
	}
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("New incoming connection: %+v\n", conn)
		go t.handleConn(conn)	
	}
}

type Temp struct {}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func(){
		fmt.Printf("dropping peer connection %s\n", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, true)
	if err = t.HandshakeFunc(peer); err != nil {
		return
	}
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
	    err := t.Decoder.Decode(conn, &rpc); 
		if err != nil {
			// fmt.Printf("TCP error %s\n", err)
			// continue
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
        // fmt.Printf("Message: %+v\n", rpc)
	}
	
}