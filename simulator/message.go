package simulator

type Message struct {
	chunkIndex       int
	dataChunkCount   int
	parityChunkCount int
}

func NewChunkedMessage(dataChunkCount int) []Message {

	var chunks []Message
	for i := 0; i < dataChunkCount; i++ {
		c := Message{
			chunkIndex:     i,
			dataChunkCount: dataChunkCount,
		}
		chunks = append(chunks, c)
	}

	return chunks
}

func NewChunkedMessageWithRedundancy(dataChunkCount int, parityChunkCount int) []Message {

	totalChunkCount := dataChunkCount + parityChunkCount

	var chunks []Message
	for i := 0; i < totalChunkCount; i++ {
		c := Message{
			chunkIndex:       i,
			dataChunkCount:   dataChunkCount,
			parityChunkCount: parityChunkCount,
		}
		chunks = append(chunks, c)
	}

	return chunks
}
