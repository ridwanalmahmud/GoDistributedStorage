package p2p

import (
	"net"
)

//message holds any arbitrary data that is being sent over the network between nodes

type RPC struct {
	From    net.Addr
	Payload []byte
}