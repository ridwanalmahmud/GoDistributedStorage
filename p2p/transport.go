package p2p

// Peer is an interface that represents the remote node
type Peer interface{
	
}

// Transport is anything that handles the communication between the nodes
type Transport interface{
	ListenAndAccept() error
}