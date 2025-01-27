package elevator

import (
	"Project/driver-go-master/elevio"
)


type DirnBehaviourPair struct{
	Direction MotorDirection
	Behaviour string
}

func requestsAbove(e Elevator) bool {
	for f := e.Floor + 1; f < elevio.NumFloors; f++ {
		for b := 0; b < elevio.NumButtons; b++ {
			if e.Requests[f][b] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for b := 0; b < elevio.NumButtons; b++ {
			if e.Requests[f][b] {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for b := 0; b < elevio.NumButtons; b++ {
		if e.Requests[e.Floor][b] {
			return true
		}
	}
	return false
}

func requestsChooseDirection(e Elevator) DirnBehaviourPair{
	switch e.Dirn {	
	case elevio.MD_Up:	
		if requestsAbove(e) {				
			return DirnBehaviourPair{Direction: elevio.MD_Up, Behaviour: "Moving"}
		}
		if requestsBelow(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Down, Behaviour: "Moving"}
		}
		if requestsHere(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: "Opening"}
		}
		return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: "Idle"}
	case elevio.MD_Down:
		if requestsBelow(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Down, Behaviour: "Moving"}
		}
		if requestsAbove(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Up, Behaviour: "Moving"}
		}
		if requestsHere(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: "Opening"}
		}
		return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: "Idle"}
	case elevio.MD_Stop:		
		if requestsAbove(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Up, Behaviour: "Moving"}
		}
		if requestsBelow(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Down, Behaviour: "Moving"}
		}	
		if requestsHere(e) {
			return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: "Opening"}
		}
		return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: "Idle"}
		
		default:
			return DirnBehaviourPair{Direction: elevio.MD_Stop, Behaviour: "Idle"}
		}
"Driver-go/elevio"switch e.Dirn {
	case elevio.MD_Down:
		return e.Requests[e.Floor][elevio.BT_HallDown] || e.Requests[e.Floor][elevio.BT_Cab] || !requestsBelow(e)
	case elevio.MD_Up:
		return e.Requests[e.Floor][elevio.BT_HallUp] || e.Requests[e.Floor][elevio.BT_Cab] || !requestsAbove(e)
	case elevio.MD_Stop:
		return true
	default:
		return false
	}
}

func requestsShouldClearImmediately(e Elevator, int btn_floor, btn_type Button){
	switch(e.elevio.clearRequest){
	case CRV_all:
		return e.floor == btn_floor
	case CRV_InDirn:
		return e.floor == btn_floor && 
		(
			(e.Dirn == elevio.MD_UP && btn_type == elevio.BT_HallUp) ||
			(e.Dirn == elevio.MD_DOWN && btn_type == elevio.BT_HallDown) ||
			e.Dirn == elevio.MD_STOP || btn_type == elevio.BT_Cab
		)
	default: 
		return false
	}
}

func ClearAtCurrentFloor(e Elevator) elevio.Elevator{
	switch e.elevio.clearRequest{
		case e.CRV_all:
			for b := 0; b < elevio.NumButtons; b++ {	
				e.Requests[e.Floor][b] = false
			}
		case e.CRV_InDirn:
			for b := 0; b < elevio.NumButtons; b++ {
				if e.Dirn == elevio.MD_UP && b == elevio.BT_HallUp || e.Dirn == elevio.MD_DOWN && b == elevio.BT_HallDown || b == elevio.BT_Cab {
					e.Requests[e.Floor][b] = false
				}
			}
	}
	return e
}

