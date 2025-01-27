package elevator

import (
	"Driver-go/elevio"
	"fmt"
)

type ElevatorBehaviour int

const (
	EB_idle ElevatorBehaviour = iota
	EB_moving
	EB_doorOpen
)

type clearRequestVariant int

const (
	CRV_all clearRequestVariant = iota
	CRV_InDirn
)

const (
	NumFloors = 4
	NumButtons = 3
)


type Elevator struct {	
	Floor int
	Dirn MotorDirection
	Requests [NumFloors][NumButtons]bool
	Behaviour ElevatorBehaviour
	Config struct {
		clearRequestVariant clearRequestVariant
		DoorOpenDuration float64
	}
}

func ElevatorUnIntialized()	Elevator {
	return Elevator{
		Floor: -1, 
		Dirn: MD_Stop,
		Behaviour: EB_idle,
	}
}

func EBtoString(eb ElevatorBehaviour)string{
	switch eb {
	case EB_idle:
		return "Idle"
	case EB_moving:
		return "Moving"
	case EB_doorOpen:
		return "DoorOpen"
	}
	return "Unknown"
}

func EDToString(dirn MotorDirection) string {
	switch dirn {
	case MD_Up:
		return "Up"
	case MD_Down:
		return "Down"
	case MD_Stop:
		return "Stop"
	}
	return "Unknown"
}


func ElevatorPrint(e Elevator){
	fmt.Println(" +-----------------+")
	fmt.Printf(
		" |  Floor: %d       |\n |  Dirn: %s       |\n |  Behaviour: %s  |\n",
		e.Floor, EDToString(e.Dirn), EBtoString(e.Behaviour),
	)
	fmt.Println(" +-----------------+")
	fmt.Println(" | | up | down | cab |")
	for f := NumFloors - 1; f >= 0; f-- {
		fmt.Printf(" | |")
		for b := 0; b < NumButtons; b++ {
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



