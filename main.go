package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"runtime"

	ch "./configuration"
	"./hardware/driver-go/elevio"
	ic "./internal_control"
	"./network"
	"./order"
)

func numRunningGoroutines() {
	for {
		fmt.Printf("\n\n\n\n\nNumber of goroutines: %d\n\n\n\n\n", runtime.NumGoroutine())
		time.Sleep(5 * time.Second)
	}
}

func initMetaDataServer(metaData ch.MetaData) <-chan ch.MetaData {

	metaDataChan := make(chan ch.MetaData, 1)

	go func() {
		for {
			metaDataChan <- metaData
		}
	}()

	return metaDataChan
}

func initChannels(metaData ch.MetaData) ch.Channels {
	chans := ch.Channels{
		DelegateOrder:     make(chan elevio.ButtonEvent),
		OrderCompleted:    make(chan elevio.ButtonEvent),
		TakingOrder:       make(chan elevio.ButtonEvent, 100),
		TakeExternalOrder: make(chan elevio.ButtonEvent),
		MetaData:          initMetaDataServer(metaData),
	}
	return chans
}

func main() {
	nodes_p := flag.Int("nodes", 3, "Number of available nodes connected to the network")
	floors_p := flag.Int("floors", 4, "Number of floors for each node")
	// Should check over network that this ID is vacant
	id_p := flag.Int("id", 0, "ID of this node")
	elevPort_p := flag.String("elev_port", "15657", "The port of the elevator to connect to (for sim purposes)")

	flag.Parse()

	nodes := *nodes_p
	floors := *floors_p
	id := *id_p
	elevPort := *elevPort_p

	if nodes <= 0 || floors <= 0 {
		log.Fatalf("Number of nodes and floors must be greater than zero! Received: %d nodes and %d floors.\n", nodes, floors)
	}

	fmt.Printf("\nInitialized with:\n\tID:\t\t%d\n\tNodes:\t\t%d\n\tFloors:\t\t%d\n\tElevPort:\t%s\n\n", id, nodes, floors, elevPort)

	metaData := ch.MetaData{NumNodes: nodes, NumFloors: floors, Id: id, ElevPort: elevPort}

	chans := initChannels(metaData)
	err := network.InitNetwork(chans.MetaData)
	if err != nil {
		log.Fatalln(err)
	}

	go numRunningGoroutines()
	go order.ControlOrders(chans)
	ic.InternalControl(chans)
}
