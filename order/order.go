package order

import (
	"fmt"

	ic "../internal_control"
	sv "../supervisor"
)

func listenForOrderCompleted() {
	//orderComplete := make(chan bool)
	// if complete
	//orderComplete :=  <-true
}

func sendOrder(order ic.OrderStruct) {
	// first send over network, needs work on

	isNotDone := make(chan bool, 1)
	isDone := make(chan bool, 1)

	go sv.WatchOrder(isNotDone)
	// go network.listenforCOmpletemsg(isDone)

	select {
	case <-isNotDone:
		fmt.Println("Order is not done")
		delegateOrder(order) //redelegate
	case <-isDone:
		return
	}
}

func delegateOrder(order ic.OrderStruct) {
	//chosenElev = lowestCost()
	//sendOrder(chosenElev)
}

func delegateOrders() {
	newOrders := make(chan ic.OrderStruct)
	for {
		select {
		case newOrder := <-newOrders:
			go delegateOrder(newOrder)
		}
	}
}

func receiveOrders() {

}

func OrderMain() {
	go delegateOrders()
	go receiveOrders()
}
