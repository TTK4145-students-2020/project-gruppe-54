package order

import (
	"fmt"

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
			if cmp.Equal(receivedDiffMsg.Order, order) && receivedDiffMsg.Diff == msgs.REMOVE {
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
		delegateOrder(order, ch) //redelegate
	case <-isDone:
		return
	}
}

func delegateOrder(order elevio.ButtonEvent, ch c.Channels) {
	//chosenElev = lowestCost()
	//sendOrder(chosenElev)
	fmt.Println("delegating")
	ch.TakeExternalOrder <- order //example of sending order to itself

}

//ControlOrders ... Delegates new orders it receives on channel newOrders
func ControlOrders(ch c.Channels) {
	go checkForExternalOrders(ch) //Continously check for new orders given to this elev
	for {
		select {
		case newOrder := <-ch.DelegateOrder:
			//fmt.Println("New order in order.go!")
			go delegateOrder(newOrder, ch)
		case <-ch.TakingOrder: // the external order has been taken
			//UpdateMatrix()          // this needs functionality
			//ch.TakingOrder <- false // reset takingorder
			//println("Taking order")

		}
	}
}
func checkForExternalOrders(ch c.Channels) {
	//odin fiks her
	//check matrix
	//update matrix so that others receive that the order is being processed
	//if new order:
	//ch.TakeExternalOrder <- order

}
