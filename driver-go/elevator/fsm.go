package elevator

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

var elevator Elevator

func FsmInit(numFloors int) {
	elevator = ElevatorUnIntialized()
	elevio.Init("localhost:15657", numFloors)
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func setAllLights() {
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elevator.Requests[f][b])
		}
	}
}

func setFloorIndicator() {
	elevio.SetFloorIndicator(elevator.Floor)
}

func FsmRequestsButtonPress(btnFloor int, btnType elevio.ButtonType) {
	fmt.Printf("\n\n%s(%d, %d)\n", btnType, btnFloor, elevator.Floor)
	ElevatorPrint(elevator)

	switch elevator.Behaviour {
	case EB_doorOpen:
		if requestsShouldClearImmediately(elevator, btnFloor, btnType) {
			ClearAtCurrentFloor(elevator)
			TimerStart(time.Duration(elevator.Config.DoorOpenDuration) * time.Second)
		} else {
			elevator.Requests[btnFloor][btnType] = true
		}
	case EB_moving:
		elevator.Requests[btnFloor][btnType] = true

	case EB_idle:
		elevator.Requests[btnFloor][btnType] = true
		pair := requestsChooseDirection(elevator)
		elevator.Dirn = pair.Direction
		elevator.Behaviour = EB_moving

		switch pair.Behaviour {

		case EB_doorOpen:
			elevator.Behaviour = EB_doorOpen
			TimerStart(time.Duration(elevator.Config.DoorOpenDuration) * time.Second)
		case EB_moving:
			elevio.SetMotorDirection(elevator.Dirn)
		case EB_idle:
			break
		}
	}
	setAllLights()
	fmt.Printf("\n New state: \n")
	ElevatorPrint(elevator)
}

func FsmFloorArrival(newFloor int) {
	fmt.Printf("\n\nFloorArrival(%d)\n", newFloor)
	ElevatorPrint(elevator)

	elevator.Floor = newFloor
	setFloorIndicator()

	switch elevator.Behaviour {

	case EB_moving:
		if requestsShouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = ClearAtCurrentFloor(elevator)
			TimerStart(time.Duration(elevator.Config.DoorOpenDuration) * time.Second)
			setAllLights()
			elevator.Behaviour = EB_doorOpen
		}
	default:
		break
	}
	fmt.Printf("\n New state: \n")
	ElevatorPrint(elevator)
}

func FsmDoorTimeout() {
	fmt.Printf("\n\nDoorTimeout()\n")
	ElevatorPrint(elevator)

	switch elevator.Behaviour {
	case EB_doorOpen:
		pair := requestsChooseDirection(elevator)
		elevator.Dirn = pair.Direction
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour {
		case EB_doorOpen:
			TimerStart(time.Duration(elevator.Config.DoorOpenDuration) * time.Second)
			elevator = ClearAtCurrentFloor(elevator)
			setAllLights()
		case EB_moving:
		case EB_idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevio.MD_Stop)
		}
	default:
		break
	}
	setAllLights()
	fmt.Printf("\n New state: \n")
	ElevatorPrint(elevator)
}
