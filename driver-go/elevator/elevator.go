package elevator

import (
	"Driver-go/config"
	"Driver-go/elevio"
)

func InitElevator() config.Elevator {
	requests := make([][]bool, 0)
	for floor := 0; floor < config.NumFloors; floor++ {
		requests = append(requests, make([]bool, config.NumButtons))
		for button := range requests[floor] {
			requests[floor][button] = false
		}
	}
	for elevio.GetFloor() == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	return config.Elevator{
		Floor:          elevio.GetFloor(),
		Direction:      elevio.MD_Stop,
		Requests:       requests,
		LocalBehaviour: config.ElevIdle,
		TimerCount:     0,
		Obstructed:     false}
}

// Set local elevtor lights and floor indicators
func SetLocalLights(e config.Elevator) {
	elevio.SetFloorIndicator(e.Floor)
	for floor := range e.Requests {
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_Cab), floor, e.Requests[floor][elevio.BT_Cab])
	}
}
