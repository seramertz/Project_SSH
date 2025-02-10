package elevator

import (
	"Driver-go/config"
	"Driver-go/elevio"
	"fmt"
)

type Behaviour int

const (
	Idle     Behaviour = 0
	DoorOpen Behaviour = 1
	Moving   Behaviour = 2
)

type Elevator struct {
	Floor      int
	Direction  elevio.MotorDirection
	Requests   [][]bool
	Behave     Behaviour
	TimerCount int
}

// Initializes an elevator to be at floor zero
func InitElevator() Elevator {
	requests := make([][]bool, 0)
	for floor := 0; floor < config.NumFloors; floor++ {
		requests = append(requests, make([]bool, config.NumButtons))
		for button := range requests[floor] {
			requests[floor][button] = false
		}
	}
	return Elevator{
		Floor:      0,
		Direction:  elevio.MD_Stop,
		Requests:   requests,
		Behave:     Idle,
		TimerCount: 0}
}

// Set elevtor lights and floor indicators
func LightsElevator(e Elevator) {
	elevio.SetFloorIndicator(e.Floor)
	for floor := range e.Requests {
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_Cab), floor, e.Requests[floor][elevio.BT_Cab])
	}
}

func EBtoString(e Elevator)string{
	switch e.Behave {
	case Idle:
		return "Idle"
	case Moving:
		return "Moving"
	case DoorOpen:
		return "DoorOpen"
	}
	return "Unknown"
}

func EDToString(dirn elevio.MotorDirection) string {
	switch dirn {
	case elevio.MD_Up:
		return "Up"
	case elevio.MD_Down:
		return "Down"
	case elevio.MD_Stop:
		return "Stop"
	}
	return "Unknown"
}

func ElevatorPrint(e Elevator){
	fmt.Println(" +-----------------+")
	fmt.Printf(
		" |  Floor: %d       |\n |  Dirn: %s       |\n |  Behaviour: %s  |\n",
		e.Floor, EDToString(e.Direction), EBtoString(e),
	)
	fmt.Println(" +-----------------+")
	fmt.Println(" | | up | down | cab |")
	for f := config.NumFloors - 1; f >= 0; f-- {
		fmt.Printf(" | |")
		for b := 0; b < config.NumButtons; b++ {
			if e.Requests[f][b] {
				fmt.Printf("  x  ")
			} else {
				fmt.Printf("     ")
			}
		}
		fmt.Println(" |")
	}	
	fmt.Println(" +-----------------+")

}
