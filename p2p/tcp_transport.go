package p2p

import (
	"fmt"
	"net"
	"log"
	"errors"
	// "bytes"
	"sync"
)

type TCPPeer struct {
	net.Conn
	// if we dial and retrieve a conn -> outbound == true
	// if we accept and retrieve a conn -> inbound == false
	outbound bool
	wg       *sync.WaitGroup
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
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
		rpcch:            make(chan RPC, 1024),
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

// consume implements the transport interface which returns read only channel for reading the incoming message from another peer

func (t *TCPTransport) Consume() <- chan RPC {
	return t.rpcch
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil{
		return err
	}
	go t.startAcceptLoop()
	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("New incoming connection: %+v\n", conn)
		go t.handleConn(conn, false)	
	}
}

type Temp struct {}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	defer func(){
		fmt.Printf("dropping peer connection %s\n", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, outbound)
	if err = t.HandshakeFunc(peer); err != nil {
		return
	}
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	
	for {
		rpc := RPC{}
	    err := t.Decoder.Decode(conn, &rpc); 
		if err != nil {
			return
		}
		rpc.From = conn.RemoteAddr().String()
		if rpc.Stream {
			peer.wg.Add(1)
		    fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
		    
		    peer.wg.Wait()
            fmt.Printf("[%s] stream closed, continuing read loop\n", conn.RemoteAddr())
			continue
		}
		t.rpcch <- rpc
	}	
}