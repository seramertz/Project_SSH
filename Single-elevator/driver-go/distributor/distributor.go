package distributor

import (
	"Driver-go/config"
	"Driver-go/elevator"
	"Driver-go/elevio"
)


func elevatorDistributorInit(id string) config.ElevatorDistributor{
	requests := make([][]config.RequestState, 4)
	for floor := range requests{
		requests[floor] = make([]config.RequestState, 3)
		
	}
	return config.ElevatorDistributor{Requests: requests, ID: id, Floor:0, Behaviour: config.Idle}

}



//distribuing orders among the elevators
func Distributor(
	id string,
	ch_newLocalOrder chan elevio.ButtonEvent,
	ch_newLocalState chan elevator.Elevator, 
	ch_orderToLocal chan elevio.ButtonEvent, 
	ch_watchdogStuckReset chan bool , 
	ch_watchdogStuckSignal chan bool, 
	ch_clearLocalHallOrders chan bool){

	
	thisElevator := new(config.ElevatorDistributor)
	*thisElevator = elevatorDistributorInit(id)

	for{
		select{
		case newOrder := <- ch_newLocalOrder:
			assignOrder(thisElevator, newOrder)
			if thisElevator.Requests[newOrder.Floor][newOrder.Button] == config.Order{
				thisElevator.Requests[newOrder.Floor][newOrder.Button] = config.Confirmed
				setHallLights(thisElevator)
				ch_orderToLocal <- newOrder
			}
			setHallLights(thisElevator)
		case newState := <- ch_newLocalState:
			if newState.Floor != thisElevator.Floor || newState.Behave == elevator.Idle || newState.Behave == elevator.DoorOpen{
				thisElevator.Behaviour = config.Behaviour(int(newState.Behave))
				thisElevator.Floor = newState.Floor
				thisElevator.Direction = config.Direction(int(newState.Direction))
				ch_watchdogStuckReset <- false
			}
			for floor := range thisElevator.Requests{
				for button := range thisElevator.Requests[floor]{
					if !newState.Requests[floor][button] && thisElevator.Requests[floor][button] == config.Confirmed{
						thisElevator.Requests[floor][button] = config.Complete
					}
					if thisElevator.Behaviour != config.Unavailable && newState.Requests[floor][button] && thisElevator.Requests[floor][button] != config.Confirmed{
						thisElevator.Requests[floor][button] = config.Confirmed
					}
				}
				
			}
			setHallLights(thisElevator)
			removeCompletedOrders(thisElevator)
		
		case <- ch_watchdogStuckSignal:
			thisElevator.Behaviour = config.Unavailable
			for floor := range thisElevator.Requests{
				for button := 0; button < len(thisElevator.Requests[floor])-1;button++{
					thisElevator.Requests[floor][button] = config.None
				}
			}
			setHallLights(thisElevator)
			ch_clearLocalHallOrders <- true
		}
	}
}


func removeCompletedOrders(elev *config.ElevatorDistributor){
	
		for floor := range elev.Requests{
			for button := range elev.Requests[floor]{
				if elev.Requests[floor][button] == config.Complete{
					elev.Requests[floor][button] = config.None
				}
			}
		}
	
}

	

func setHallLights(elevator *config.ElevatorDistributor) {
	for button := 0 ; button < config.NumButtons - 1 ; button++{
		for floor := 0 ; floor < config.NumFloors ; floor++{
			isLight := false
			
			if elevator.Requests[floor][button] == config.Confirmed{
				isLight = true
				
			}
			elevio.SetButtonLamp(elevio.ButtonType(button), floor, isLight)
		}
	}
}

func assignOrder(elevator *config.ElevatorDistributor, order elevio.ButtonEvent) {
	elevator.Requests[order.Floor][order.Button] = config.Order
		
}