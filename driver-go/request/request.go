package request

import(
	"driver-go/elevio"
	"driver-go/elevator"
	"driver-go/config"
)

func RequestAbove(e elevator.Elevator) bool{
	for floor := e.Floor+1; floor < config.NumFloors; floor++{
		for button := 0; button < config.NumButtons; button++{
			if e.Requests[floor][button]{
				return true
			}
		}
	}
	return false
}