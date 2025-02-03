package elevator

import (
	"root/config"
	"root/elevio"
)

type Orders [config.NumFloors][config.NumButtons]bool

func (orders Orders) OrderInDirection(floor int, dir Direction) bool {
	switch dir {
	case Up:
		for f := floor + 1; f < config.NumFloors; f++ {
			for b := 0; b < config.NumButtons; b++ {
				if orders[f][b] {
					return true
				}
			}
		}
		return false
	case Down:
		for f := 0; f < floor; f++ {
			for b := 0; b < config.NumButtons; b++ {
				if orders[f][b] {
					return true
				}
			}
		}
		return false
	default:
		panic("Invalid direction")
	}
}

func OrderDone(floor int, dir Direction, orders Orders, orderDoneChannel chan<- elevio.ButtonEvent) {
	if orders[floor][elevio.BT_Cab] {
		orderDoneChannel <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if oders[floor][dir] {
		orderDoneChannel <- elevio.ButtonEvent{Floor: floor, Button: dir.toBT()}
	}
}
