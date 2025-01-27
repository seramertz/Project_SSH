package elevator

import(
	"Project/driver-go-master/elevio"
	"fmt"
	"elevator"
	"requests"
	"timer"

)

var elevator elevio.Elevator


func init(){
	elevator = elevator.ElevatorUnIntialized()
	elevio.Init("localhost:15657", 4)
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func setAllLights(){
	for f := 0; f < elevio.NumFloors; f++ {
		for b := 0; b < elevio.NumButtons; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elevator.Requests[f][b])
		}
	}
}

func setFloorIndicator(){
	elevio.SetFloorIndicator(elevator.Floor)
}

func requestsButtonPress(btnFloor int, btnType elevio.ButtonType){
	fmt.Printf("\n\n%s(%d, %d)\n", btnType, btnFloor, elevator.Floor)
	elevator.ElevatorPrint(elevator)

	switch elevator.Behaviour {
	case elevator.EB_doorOpen:
		if requests.requestsShouldClearImmediately(elevator, btnFloor, btnType){
			requests.requestsClearAtCurrentFloor(elevator, btnFloor, btnType)
			timer.TimerStart(elevator.Config.DoorOpenDuration)
		}else{
			elevator.Requests[btnFloor][btnType] = true
		}
	case elevator.EB_moving:
		elevator.Requests[btnFloor][btnType] = true

	case elevator.EB_idle:
		elevator.Requests[btnFloor][btnType] = true
		pair := requests.requestsChooseDirection(elevator)
		elevator.Dirn = pair.Direction
		elevator.Behaviour = elevator.EB_moving

		switch pair.Behaviour{

		case elevator.EB_doorOpen:
			elevator.Behaviour = elevator.EB_doorOpen
			timer.TimerStart(elevator.Config.DoorOpenDuration)
		case elevator.EB_moving:
			elevio.SetMotorDirection(elevator.Dirn)
		case elevator.EB_idle: 
			break
		}
	}
	setAllLights()
	fmt.Printf("\n New state: \n")
	elevator.ElevatorPrint(elevator)
}

func FsmFloorArrival(newFloor int){
	fmt.Printf("\n\nFloorArrival(%d)\n", newFloor)
	elevator.ElevatorPrint(elevator)

	elevator.Floor = newFloor
	setFloorIndicator(elevator.Floor)

	switch elevator.Behaviour{

	case elevator.EB_moving:
		if requests.requestsShouldStop(elevator){
			SetMotorDirection(elevator, elevio.MD_Stop)
			SetDoorOpenLamp(elevator, true)
			elevator = requests.requestsClearAtCurrentFloor(elevator)
			timer.TimerStart(elevator.Config.DoorOpenDuration)
			setAllLights()
			elevator.Behaviour = elevator.EB_doorOpen
		}
	default:
		break
	}
	fmt.Printf("\n New state: \n")
	elevator.ElevatorPrint(elevator)
}

func FsmDoorTimeout(){
	fmt.Printf("\n\nDoorTimeout()\n")
	elevator.ElevatorPrint(elevator)

	switch elevator.Behaviour{
	case elevator.EB_doorOpen:
		pair := requests.requestsChooseDirection(elevator)
		elevator.Dirn = pair.Direction
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour{
		case elevator.EB_doorOpen:
			timer.TimerStart(elevator.Config.DoorOpenDuration)
			elevator = requests.requestsClearAtCurrentFloor(elevator)
			setAllLights()
		case elevator.EB_moving:
		case elevator.EB_idle:
			SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevio.MD_Stop)
		}
	default:
		break
	}
	setAllLights()
	fmt.Printf("\n New state: \n")
	elevator.ElevatorPrint(elevator)
}