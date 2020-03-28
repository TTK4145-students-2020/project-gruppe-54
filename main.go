package main

import (
	"flag"
	"fmt"
	"log"

	ch "./configuration"
	"./hardware/driver-go/elevio"
	ic "./internal_control"
	"./network"
	"./order"
)

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
		TakingOrder:       make(chan elevio.ButtonEvent),
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

	flag.Parse()

	nodes := *nodes_p
	floors := *floors_p
	id := *id_p

	if nodes <= 0 || floors <= 0 {
		log.Fatalf("Number of nodes and floors must be greater than zero! Received: %d nodes and %d floors.\n", nodes, floors)
	}

	fmt.Printf("\nInitialized with:\n\tID:\t%d\n\tNodes:\t%d\n\tFloors:\t%d\n\n", id, nodes, floors)

	metaData := ch.MetaData{NumNodes: nodes, NumFloors: floors, Id: id}

	chans := initChannels(metaData)
	err := network.InitNetwork(chans.MetaData)
	if err != nil {
		log.Fatalln(err)
	}

	go order.ControlOrders(chans)
	ic.InternalControl(chans)
}
