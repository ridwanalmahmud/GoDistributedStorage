package p2p

// Peer is an interface that represents the remote node
type Peer interface{
	Close() error
}

// Transport is anything that handles the communication between the nodes
type Transport interface{
	ListenAndAccept() error
	Consume() <-chan RPC
}