package supervisor

import (
	"time"

	order "../order"
	"./watchdog"
)

var timeoutPeriod int = 5

type orderStruct = order.OrderStruct

//WatchOrder ... watches an order
//start as goroutine
// still missing orderComplete functionality (needs connecting channel)
func WatchOrder() {
	wd := watchdog.NewWatchdog(time.Second * timeoutPeriod)
	orderComplete := make(chan bool)
	select {
	case <-wd.GetKickChannel():
		resendOrder(order)
		wd.Stop()
	case <-orderCompleted():
		return
	}
}

func resendOrder(order orderStruct) {
	newOrders := make(chan elevio.ButtonEvent)
	newOrders <- order
}
func main() {
	var a orderStruct
	a.Floor = 5
	//stop := make(chan bool)
	go NewOrder(a)
	for {

	}
	//stop <- true
	//<-stop
}
