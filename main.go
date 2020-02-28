package main

import (
	"time"
	"fmt"
	"./internal_control"
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
	time.AfterFunc(1 * time.Second, func () {
		callback2(x, wait)
	})
	// timer := time.NewTimer(1 * time.Second)
	// go timerCallback(timer, 20, wait)
	<-wait
}
*/


func main(){
	go InternalControl()
}