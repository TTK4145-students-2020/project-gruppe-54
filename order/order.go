package order

import (
	"fmt"

	"../hardware/driver-go/elevio"
	ic "../internal_control"
	sv "../supervisor"
)

func listenForOrderCompleted(isDone chan bool) {
	var complete = true //Odin sett in noe her om at den hører på en port elns og hvis den får inn noe blir det true
	if complete {
		isDone <- true
	}
}

func sendOrder(order elevio.ButtonEvent) {
	// send ordren over nettverket, odin fiks her

	isNotDone := make(chan bool, 1)
	isDone := make(chan bool, 1)

	go sv.WatchOrder(isNotDone)
	go listenForOrderCompleted(isDone)

	select {
	case <-isNotDone:
		fmt.Println("Order is not done")
		delegateOrder(order) //redelegate
	case <-isDone:
		return
	}
}

func delegateOrder(order elevio.ButtonEvent) {
	//chosenElev = lowestCost()
	//sendOrder(chosenElev)
	fmt.Println("delegating")
	ic.AddOrder(order)
}

func DelegateOrders(newOrders chan elevio.ButtonEvent) {
	for {
		select {
		case newOrder := <-newOrders:
			fmt.Println("New order in order.go!")
			go delegateOrder(newOrder)
		}
	}
}

func receiveOrders() {
	//odin fiks her så den hører på en port for ordre

}

func ControlOrders() {
	initOrderMatrix()
}
