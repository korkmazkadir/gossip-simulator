package simulator

import (
	"math/rand"
	"time"
)

type System struct {
	nodeCount         int
	fanout            int
	faultyNodePercent float64
	nodes             []*Node
}

func NewSystem(nodeCount int, fanout int, faultPercent float64) *System {

	if faultPercent > 1 || faultPercent < 0 {
		panic("fault percent must be between 0 and 1")
	}

	// creates a system
	s := &System{
		nodeCount:         nodeCount,
		fanout:            fanout,
		faultyNodePercent: faultPercent,
	}

	// creates nodes
	for i := 0; i < nodeCount; i++ {
		s.nodes = append(s.nodes, NewNode(i, s))
	}

	// shuffles the node slice
	rand.Seed(time.Now().UnixNano())
	for i := range s.nodes {
		j := rand.Intn(i + 1)
		s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
	}

	faultyNodeCount := int(float64(nodeCount) * faultPercent)
	for i := 0; i < faultyNodeCount; i++ {
		// marks nodes as faulty
		s.nodes[i].isFaulty = true
	}

	//log.Printf("system initialized: node count: %d, fanout: %d, fault per: %f ,fault count: %d", s.nodeCount, s.fanout, s.faultyNodePercent, faultyNodeCount)

	return s
}

func (s *System) Sample(except int) []*Node {

	return s.SampleN(s.fanout, except)
}

func (s *System) SampleN(n int, except int) []*Node {

	dict := make(map[int]struct{})
	dict[except] = struct{}{}

	var sample []*Node
	for len(sample) < n {

		node := s.nodes[rand.Intn(s.nodeCount)]
		_, ok := dict[node.id]
		if ok {
			continue
		}

		sample = append(sample, node)
		dict[node.id] = struct{}{}
	}

	return sample
}
