package watchdog

import (
	"fmt"
	"time"
)

// Example .. Showing an example of how the watchdog can be implemented
func Example() {
	fmt.Println("Starting watchdog")
	stop := make(chan bool)
	go func() {
		wd := NewWatchdog(time.Second * 3)
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
	wd2 := NewWatchdog(time.Second * 3)
	if <-wd2.GetKickChannel() {
		fmt.Println("wd2")
		wd2.Stop()
	}
	<-stop
}
