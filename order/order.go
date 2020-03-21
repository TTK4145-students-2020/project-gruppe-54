package order

import (
	"fmt"
	"math"

	c "../configuration"
	"../hardware/driver-go/elevio"
	"../network/msgs"
	sv "../supervisor"
	"github.com/google/go-cmp/cmp"
)

var numNodes int = 3
var numFloors int = 4

func listenForOrderCompleted(order elevio.ButtonEvent, isDone chan bool, giveUp chan bool) {
	orderComplete := false //Odin sett in noe her om at den ser på matrisen og hvis den får inn at det er utført blir det true

	// Listen for diff in order tensor
	orderTensorDiffMsg := msgs.OrderTensorDiffMsg{}
	receivedDiff := make(chan msgs.OrderTensorDiffMsg)
	killPolling := make(chan bool)

	defer func() {
		killPolling <- true
		if orderComplete {
			isDone <- true
		}
	}()

	// Start polling for diff to tensor
	go func() {
		for {
			err := orderTensorDiffMsg.Listen()
			if err != nil {
				continue
			}
			receivedDiff <- orderTensorDiffMsg
			if <-killPolling {
				return
			}
		}
	}()
ListenLoop:
	for {
		select {
		case receivedDiffMsg := <-receivedDiff:
			if cmp.Equal(receivedDiffMsg.Order, order) && receivedDiffMsg.Diff == msgs.DIFF_REMOVE {
				orderComplete = true
				break ListenLoop
			}
		case <-giveUp:
			orderComplete = false
			break ListenLoop
		}
	}
}

func sendOrder(order elevio.ButtonEvent, ch c.Channels) {
	orderMsg := msgs.OrderMsg{Order: order}
	orderMsg.Send()

	isNotDone := make(chan bool, 1)
	isDone := make(chan bool, 1)

	go sv.WatchOrder(isNotDone)
	go listenForOrderCompleted(order, isDone, isNotDone)
	select {
	case <-isNotDone:
		fmt.Println("Order is not done")
		// FIXME: Dirty fix
		delegateOrder(order, ch, 3) //redelegate
	case <-isDone:
		// Nothing to do
	}
}

func delegateOrder(order elevio.ButtonEvent, ch c.Channels, numNodes int) {
	// TODO:
	// 1. Send order to all PCs on network
	// 2. Wait for cost in return, needs to count or return after some time
	// 3. If it deems itself as most fit, it passes the order to itself. Anyhow, a watchdog is created to ensure the order is completed

	// 1.
	orderMsg := msgs.OrderMsg{Order: order}
	orderMsg.Send()

	// 2.
	costs := collectCosts(numNodes)

	//chosenElev = lowestCost()
	//sendOrder(chosenElev)
	fmt.Println("delegating")
	ch.TakeExternalOrder <- order //example of sending order to itself

}

func collectCosts(numNodes int) []uint {
	costs := make([]uint, numNodes)
	for i := 0; i < numNodes; i++ {
		costs[i] = uint(math.Inf(0)) // Initialize all costs to infinite
	}

	return costs
}

//ControlOrders ... Delegates new orders it receives on channel newOrders
func ControlOrders(ch c.Channels, metaDataChan <-chan c.MetaData) {
	metaData := <-metaDataChan
	numNodes, _, ID := metaData.NumNodes, metaData.NumFloors, metaData.Id
	// ID := metaData.Id

	// updateOrderTensorCh, currentOrderTensorCh := InitOrderTensor(numNodes, numFloors)

	go checkForExternalOrders(ch, ID) //Continously check for new orders given to this elev
	for {
		select {
		case newOrder := <-ch.DelegateOrder:
			//fmt.Println("New order in order.go!")
			go delegateOrder(newOrder, ch, numNodes)
		case <-ch.TakingOrder: // the external order has been taken
			//UpdateMatrix()          // this needs functionality
			//ch.TakingOrder <- false // reset takingorder
			//println("Taking order")

		}
	}
}

func checkForExternalOrders(ch c.Channels, ID int) {
	newOrder := msgs.OrderMsg{}
	for {
		err := newOrder.Listen()
		if err != nil {
			continue
		}
		cost := calculateCost(newOrder.Order)
		costMsg := msgs.CostMsg{Cost: cost, Id: ID}
		costMsg.Send()
	}
}
