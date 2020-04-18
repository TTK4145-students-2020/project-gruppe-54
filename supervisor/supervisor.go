package supervisor

import (
	"time"

	"./watchdog"
)

var timeoutPeriod time.Duration = 5 * time.Second

//WatchOrder ... watches an order
func WatchOrder(isNotDone chan bool) {
	wd := watchdog.NewWatchdog(timeoutPeriod)
	select { // whichever one of order complete or wd finished first
	case <-wd.GetKickChannel():
		isNotDone <- true
		wd.Stop()
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
