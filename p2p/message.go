package p2p

import (
	//"net"
)

const (
	IncomingMessage = 0x1
	IncomingStream = 0x2
)

//message holds any arbitrary data that is being sent over the network between nodes

type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}