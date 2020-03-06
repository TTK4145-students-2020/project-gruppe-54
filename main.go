package main

import (
	"fmt"
	"os"
	"time"

	ic "./internal_control"
	order "./order"
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
	ic.InternalControl()
	order.OrderMain()
	ID := os.Args[1:2]
	fmt.Println(ID)
}
