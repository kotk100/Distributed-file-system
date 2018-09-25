package kademlia

import "./protocol"

type RequestExecutor interface {
	execute()
	setChannel(ch chan *protocol.RPC)
	setMessageId(id messageID)
}
