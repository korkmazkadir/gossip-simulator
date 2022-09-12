package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/korkmazkadir/gossip-simulator/simulator"
)

func main() {

	dissmeinationType := flag.String("type", "classic", "type of dissemination ida/classic")
	fanOutF := flag.Int("d", 0, "fanout")
	faultF := flag.Float64("f", 0.0, "faulty node percent between 0 and 1")
	experimentCountF := flag.Int("e", 0, "experiment count")
	flag.Parse()

	nodeCount := 4096
	fanout := *fanOutF
	faultPercent := *faultF

	dataChunkCount := 96
	parityChunkCount := 160

	//log.Println(*dissmeinationType)

	var stats []simulator.DisseminationStats

	experimentCount := *experimentCountF

	for i := 0; i < experimentCount; i++ {

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

		stats = append(stats, stat)

		//log.Printf("----> Message Delivery Count %d\n", stat.MessageDeliveryCount)

	}

	totalRoundCount := 0
	totalDeliveryPercent := 0.0
	totalForwardCount := 0.0
	totalReceivedCount := 0.0

	failureCount := 0

	faultyNodeCount := int(float64(nodeCount) * faultPercent)

	for i := range stats {
		totalRoundCount += stats[i].Round
		totalDeliveryPercent += (float64(stats[i].MessageDeliveryCount) / float64(nodeCount))
		totalForwardCount += float64(stats[i].DeliveredChunkCount) / float64(nodeCount-faultyNodeCount)
		totalReceivedCount += float64(stats[i].ReceivedChunkCount) / float64(nodeCount)

		if stats[i].MessageDeliveryCount == 0 {
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
