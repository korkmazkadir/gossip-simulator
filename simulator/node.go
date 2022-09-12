package simulator

type stats struct {
	receivedMessageCount int

	forwardedMessageCount int
}

type Node struct {

	// id of a node
	id int

	// Contains messages that send to current node on the current node
	inbox []Message

	// Contains messages received by the node, and waiting to be forwarded
	outbox []Message

	// Delivered Messages
	delivered      map[int]Message
	deliveredSlice []Message

	// Signals the current node is faulty
	isFaulty bool

	// Signals the current node delivered all messages so it will not contrubute anymore
	isFinished bool

	system *System

	// statistics
	stats
}

func NewNode(id int, system *System) *Node {
	node := &Node{id: id, system: system, delivered: make(map[int]Message), isFaulty: false}
	return node
}

func (n *Node) DisseminateMessages(messages []Message) {
	n.outbox = append(n.outbox, messages...)
}

func (n *Node) Deliver() {

	if n.isFinished {
		n.receivedMessageCount += len(n.inbox)
		n.inbox = nil
		return
	}

	for _, message := range n.inbox {
		// checks a message is already delivered
		_, ok := n.delivered[message.chunkIndex]
		if !ok {

			if len(n.delivered) == message.dataChunkCount {
				// if a node has delivered enough message to reconstruct the original message
				// it finishes the dissemination
				n.isFinished = true
			} else {
				// the message is marked as delivered
				n.delivered[message.chunkIndex] = message
				n.deliveredSlice = append(n.deliveredSlice, message)

				// the message is put into outbox to be forwarded if it is not faulty
				if !n.isFaulty {
					n.outbox = append(n.outbox, message)
				}

				// this is for classic
				if len(n.delivered) == message.dataChunkCount {
					n.isFinished = true
				}

			}

		}
	}

	// keeps received message count
	n.receivedMessageCount += len(n.inbox)
	// empties the inbox
	n.inbox = nil

}

// Forward returns forwarded message count
func (n *Node) Forward() int {

	if n.isFaulty {
		// faulty nodes does not forward a message
		return 0
	}

	forwardedMessageCount := 0

	for _, message := range n.outbox {

		// samples a distinc peer set for each message
		peers := n.system.Sample(n.id)
		for _, peer := range peers {
			// messages put into the inbox so they will wait to be delivered
			peer.inbox = append(peer.inbox, message)
		}
		forwardedMessageCount++
	}

	n.outbox = nil
	n.forwardedMessageCount += forwardedMessageCount

	return forwardedMessageCount
}
