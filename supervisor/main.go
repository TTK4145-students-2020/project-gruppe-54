package main

import (
	"fmt"
	"time"

	order "../order"
	"./watchdog"
)

var timeoutPeriod int = 5 

type orderStruct = order.OrderStruct

// NewOrder ... start timer on new order
func NewOrder(order orderStruct) {
	t := order.Floor
	fmt.Println("floor:", t)
	startNewTimer(time.Duration(timeoutPeriod))
}

func startNewTimer(t time.Duration) {
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
