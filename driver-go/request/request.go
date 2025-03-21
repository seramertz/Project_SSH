package request


import (
	"Driver-go/config"
	"Driver-go/elevator"
	"Driver-go/elevio"
)


func RequestAbove(elev elevator.Elevator) bool{
	for floor := elev.Floor+1; floor < config.NumFloors; floor++{
		for button := 0; button < config.NumButtons; button++{
			if elev.Requests[floor][button]{
				return true
			}
		}
	}
	return false
}

func RequestBelow(elev elevator.Elevator)bool {
	for floor := 0; floor < elev.Floor; floor++{
		for btn := range elev.Requests[floor]{
			if elev.Requests[floor][btn]{
				return true
			}
		}
	}
	return false
}

//Clears the requests at the current floor going in the same direction
func RequestClearAtCurrentFloor(elev *elevator.Elevator){
	elev.Requests[elev.Floor][int(elevio.BT_Cab)] = false
	switch{
	case elev.Direction == elevio.MD_Up:
		elev.Requests[elev.Floor][int(elevio.BT_HallUp)] = false
		if !RequestAbove(*elev){
			elev.Requests[elev.Floor][int(elevio.BT_HallDown)] = false
		}
	case elev.Direction == elevio.MD_Down:
		elev.Requests[elev.Floor][int(elevio.BT_HallDown)] = false
		if !RequestBelow(*elev){
			elev.Requests[elev.Floor][int(elevio.BT_HallUp)] = false
		}
	}
}

//Stop based on current floor and direction
func RequestShouldStop(elev *elevator.Elevator)bool {
	switch{
	case elev.Direction == elevio.MD_Down:
		return elev.Requests[elev.Floor][int(elevio.BT_HallDown)] || elev.Requests[elev.Floor][int(elevio.BT_Cab)] || !RequestBelow(*elev)
	case elev.Direction == elevio.MD_Up:
		return elev.Requests[elev.Floor][int(elevio.BT_HallUp)] || elev.Requests[elev.Floor][int(elevio.BT_Cab)] || !RequestAbove(*elev)
	default:
		return true
	}
}

//Chooses direction based on the the requests made 
func RequestChooseDirection(elev *elevator.Elevator){
	switch elev.Direction{
	case elevio.MD_Up:
		if RequestAbove(*elev){
			elev.Direction = elevio.MD_Up
		} else if RequestBelow(*elev){
			elev.Direction = elevio.MD_Down
		} else{
			elev.Direction = elevio.MD_Stop
		}
	case elevio.MD_Down:
		fallthrough
	
	case elevio.MD_Stop:
		if RequestBelow(*elev){
			elev.Direction = elevio.MD_Down
		} else if RequestAbove(*elev){
			elev.Direction = elevio.MD_Up
		} else{
			elev.Direction = elevio.MD_Stop
		}
	}

}



func RequestClearHall(elev *elevator.Elevator){
	for floor := 0; floor < config.NumFloors; floor++{
		for btn := 0; btn < config.NumButtons; btn++{
			elev.Requests[floor][btn] = false
		}
	}
}

