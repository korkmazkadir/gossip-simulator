package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/korkmazkadir/gossip-simulator/simulator"
)

func main() {

	dissmeinationType := flag.String("type", "classic", "type of dissemination ida/classic")
	fanOutF := flag.Int("d", 0, "fanout")
	faultF := flag.Float64("f", 0.0, "faulty node percent between 0 and 1")
	experimentCountF := flag.Int("e", 0, "experiment count")
	dataChunkF := flag.Int("dc", 0, "data chunk count")
	parityChunkF := flag.Int("pc", 0, "parity chunk count")
	flag.Parse()

	nodeCount := 4096
	fanout := *fanOutF
	faultPercent := *faultF

	dataChunkCount := *dataChunkF
	parityChunkCount := *parityChunkF

	//log.Println(*dissmeinationType)

	experimentCount := *experimentCountF

	stats := make([]simulator.DisseminationStats, experimentCount)

	var wg sync.WaitGroup

	for j := 0; j < 10; j++ {

		wg.Add(1)
		go func(index int) {

			for i := 0; i < experimentCount/10; i++ {
				system := simulator.NewSystem(nodeCount, fanout, faultPercent)
				disseminator := simulator.NewDisseminator(system)

				var stat simulator.DisseminationStats

				switch *dissmeinationType {
				case "classic":
					stat = disseminator.DisseminateClassic(dataChunkCount)
				case "ida":
					stat = disseminator.DisseminateIDA(dataChunkCount, parityChunkCount)
				default:
					panic(fmt.Errorf("unknow dissemination type %s", *dissmeinationType))
				}

				stats[(index*(experimentCount/10))+i] = stat
			}

			wg.Done()
		}(j)
	}

	wg.Wait()

	totalRoundCount := 0
	totalDeliveryPercent := 0.0
	totalForwardCount := 0.0
	totalReceivedCount := 0.0

	failureCount := 0

	faultyNodeCount := int(float64(nodeCount) * faultPercent)

	for _, stat := range stats {
		totalRoundCount += stat.Round
		totalDeliveryPercent += (float64(stat.MessageDeliveryCount) / float64(nodeCount))
		totalForwardCount += float64(stat.ForwardedChunkCount) / float64(nodeCount-faultyNodeCount)
		totalReceivedCount += float64(stat.ReceivedChunkCount) / float64(nodeCount)

		if stat.MessageDeliveryCount == 0 {
			failureCount++
		}

	}

	avgRoundCount := float64(totalRoundCount) / float64(experimentCount)
	avgDeliveryPercent := totalDeliveryPercent / float64(experimentCount)
	avgForwardCount := float64(totalForwardCount) / float64(experimentCount)
	avgReceivedCount := float64(totalReceivedCount) / float64(experimentCount)

	log.Printf("%d Experiments: avgRoundCount %f, avgDeliveryPercent %f, avgForwardCound %f unsuccessCount %d\n", experimentCount, avgRoundCount, avgDeliveryPercent, avgForwardCount, failureCount)

	outputMap := make(map[string]interface{})
	outputMap["ExperimentType"] = *dissmeinationType
	outputMap["NodeCount"] = nodeCount
	outputMap["Fanout"] = fanout
	outputMap["FaultyNodePercent"] = faultPercent

	outputMap["DataChunkCount"] = dataChunkCount
	outputMap["ParityChunkCount"] = dataChunkCount

	outputMap["ExperimentCount"] = experimentCount

	outputMap["AverageRoundCount"] = avgRoundCount
	outputMap["AverageDeliveryPercent"] = avgDeliveryPercent
	outputMap["AverageForwardCount"] = avgForwardCount
	outputMap["AverageReceivedCount"] = avgReceivedCount
	outputMap["FailureCount"] = failureCount

	jsonBytes, err := json.Marshal(outputMap)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))

	// bytes, err := json.Marshal(stats)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(bytes))

}
