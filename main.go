package main

import (
	"fmt"
	"os"
	"time"

	ch "./configuration"
	"./hardware/driver-go/elevio"
	ic "./internal_control"
	"./order"
)

func initChannels() ch.Channels {
	chans := ch.Channels{
		DelegateOrder:     make(chan elevio.ButtonEvent),
		OrderCompleted:    make(chan bool),
		TakingOrder:       make(chan bool),
		TakeExternalOrder: make(chan elevio.ButtonEvent),
	}
	return chans
}

func timerCallback(timer *time.Timer, x int, wait chan int) {
	<-(*timer).C
	fmt.Printf("Got %d after 1 seconds\n", x)
	wait <- 0
}

func callback2(x int, wait chan int) {
	fmt.Printf("Callback %d\n", x)
	wait <- 0
}

/*
func main() {
	wait := make(chan int)
	x := 1
	time.AfterFunc(1*time.Second, func() {
		callback2(x, wait)
	})
	// timer := time.NewTimer(1 * time.Second)
	// go timerCallback(timer, 20, wait)
	<-waita
}

func main() {
	watchdog.Example()
}
*/

func main() {
	//fmt.Println("testing internal control")
	//order.initOrderMatrix(numNodes, numFloors)
	var chans = initChannels()
	go order.ControlOrders(chans)
	ic.InternalControl(chans)
	ID := os.Args[1:2]
	fmt.Println(ID)
}
