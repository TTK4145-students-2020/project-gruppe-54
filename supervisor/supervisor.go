package supervisor

import (
	"time"

	"./watchdog"
)

var timeoutPeriod int = 5

//WatchOrder ... watches an order
//start as goroutine
// still missing orderComplete functionality (needs connecting channel)
func WatchOrder(isDone chan bool) {
	wd := watchdog.NewWatchdog(time.Second * time.Duration(timeoutPeriod))
	orderComplete := make(chan bool)
	select { // whichever one of order complete or wd finished first
	case <-wd.GetKickChannel():
		isDone <- true
		wd.Stop()
	case <-orderComplete:
		isDone <- false
		return
	}
}

func main() {
	//stop := make(chan bool)
	//go NewOrder(a)
	//for {

	//}
	//stop <- true
	//<-stop
}
