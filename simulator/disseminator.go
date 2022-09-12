package simulator

type DisseminationStats struct {
	Round                int
	ReceivedChunkCount   int
	DeliveredChunkCount  int
	MessageDeliveryCount int
}

type Disseminator struct {
	system *System
}

func NewDisseminator(system *System) *Disseminator {
	return &Disseminator{system: system}
}

func (d *Disseminator) DisseminateClassic(dataChunkCoun int) DisseminationStats {
	messages := NewChunkedMessage(dataChunkCoun)

	sorceNodes := d.system.SampleN(1, -1)
	// necessary to select a correct source
	for sorceNodes[0].isFaulty {
		sorceNodes = d.system.SampleN(1, -1)
	}

	sorceNodes[0].DisseminateMessages(messages)
	return d.disseminate(0)
}

func (d *Disseminator) DisseminateIDA(dataChunkCoun int, parityChunkCoun int) DisseminationStats {

	messages := NewChunkedMessageWithRedundancy(dataChunkCoun, parityChunkCoun)
	sorceNodes := d.system.SampleN(16, -1)

	faultyCount := 0
	for i := range sorceNodes {
		if sorceNodes[i].isFaulty {
			faultyCount++
		}
	}

	// if faultyCount > 10 {
	// 	log.Printf("Failed... Faulty Node Count: %d\n", faultyCount)
	// }

	for i, m := range messages {
		index := i % 16
		sorceNodes[index].DisseminateMessages([]Message{m})
	}
	return d.disseminate(1)
}

func (d *Disseminator) disseminate(round int) DisseminationStats {

	nodes := d.system.nodes
	currentRound := round

	for {
		forwardCount := 0
		for _, n := range nodes {
			forwardCount += n.Forward()
		}

		if forwardCount == 0 {
			break
		}

		for _, n := range nodes {
			n.Deliver()
		}
		currentRound++
	}

	stats := DisseminationStats{Round: currentRound}

	for _, n := range nodes {

		if n.isFinished {
			stats.MessageDeliveryCount++
		}

		stats.DeliveredChunkCount += n.stats.forwardedMessageCount
		stats.ReceivedChunkCount += n.stats.receivedMessageCount
	}

	return stats
}
