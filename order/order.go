package order

import (
	"fmt"

	c "../configuration"
	"../hardware/driver-go/elevio"
	sv "../supervisor"
)

var numNodes int = 3
var numFloors int = 4

func listenForOrderCompleted(isDone chan bool) {
	var complete = true //Odin sett in noe her om at den ser på matrisen og hvis den får inn at det er utført blir det true
	if complete {
		isDone <- true
	}
}

func sendOrder(order elevio.ButtonEvent, ch c.Channels) {
	// send ordren over nettverket, odin fiks her

	isNotDone := make(chan bool, 1)
	isDone := make(chan bool, 1)

	go sv.WatchOrder(isNotDone)
	go listenForOrderCompleted(isDone)

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
