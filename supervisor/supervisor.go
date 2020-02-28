package supervisor

import (
	"fmt"
	"time"

	order "../order"
	"./watchdog"
)

var timeoutPeriod int = 5

type orderStruct = order.OrderStruct

// NewOrder ... start timer on new order, sends the order as new order if kicked
// returns nothing, lives its own life
func WatchNewOrder(order orderStruct) {
	t := order.Floor
	fmt.Println("floor:", t)
	NewWatchdog(time.Duration(timeoutPeriod))
}

func NewWatchdog(t time.Duration) {
	wd := watchdog.NewWatchdog(time.Second * t)
	if <-wd.GetKickChannel() {
		fmt.Println("wd2")
		orderNotCompleted()
		wd.Stop()
		return
	}
}

func orderNotCompleted() {
	//	send new order to orders
	// missing functionality
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
