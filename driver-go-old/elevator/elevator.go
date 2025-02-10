package elevator

import (
	"Driver-go/elevio"
	"Driver-go/config"
)


type Behaviour int

const(
	Idle Behaviour= iota
	DoorOpen 
	Moving 
)

type Elevator struct{
	Floor int
	Direction elevio.MotorDirection
	Requests [][]bool
	Behave Behaviour
	TimerCount int
}

//Initializes an elevator to be at floor zero
func InitElevator() Elevator{
	requests := make([][]bool, 0)
	for floor := 0; floor < config.NumFloors; floor++ {
		requests = append(requests, make([]bool, config.NumButtons))
		for button := range requests[floor] {
			requests[floor][button] = false
		}
	}
	return Elevator{
		Floor: 0,
		Direction: elevio.MD_Stop,
		Requests: requests,
		Behave: Idle,
		TimerCount: 0}
}

//Set elevtor lights and floor indicators
func LightsElevator(e Elevator){
	elevio.SetFloorIndicator(e.Floor)
	for floor := range e.Requests{
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_Cab), floor, e.Requests[floor][elevio.BT_Cab])
	}
}
