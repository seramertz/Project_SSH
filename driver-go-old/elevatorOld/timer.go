package elevator

import (
	"time"	
)

var timerEndTime time.Time
var timerActive bool

func getTime() time.Time {
	return time.Now()
}

func TimerStart(duration time.Duration) {
	timerEndTime = getTime().Add(duration)
	timerActive = true
}

func TimerStop() {
	timerActive = false
}

func TimerTimeOut() bool {
	return timerActive && getTime().After(timerEndTime)
}
