package main

import (
	"fmt"
	"os"
	"time"

	"./hardware/driver-go/elevio"
	ic "./internal_control"
	"./order"
)

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
	newOrders := make(chan elevio.ButtonEvent)

	go order.DelegateOrders(newOrders)
	ic.InternalControl(newOrders)
	ID := os.Args[1:2]
	fmt.Println(ID)
}
