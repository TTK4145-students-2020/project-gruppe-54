package main

import (
	"fmt"
	"time"

	"./watchdog"
)

func main() {
	fmt.Println("Starting watchdog")
	stop := make(chan bool)
	go func() {
		fmt.Println("HELLO")
		wd := watchdog.NewWatchdog(time.Second * 3)
		i := 0
		for {
			select {
			case <-wd.GetKickChannel():
				fmt.Println("kick!")
				i++
				if i == 3 {
					fmt.Println("STOPPING")
					// stop the watchdog permanently
					stop <- true
					wd.Stop()
					return
				}
			}
		}
	}()
	<-stop
}
