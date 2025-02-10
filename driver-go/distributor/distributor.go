package distributor

import (
	"Driver-go/assigner"
	"Driver-go/config"
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/network/peers"
	"Driver-go/request"
	"time"
)

const localElevator = 0

func elevatorDistributorInit(id string) config.ElevatorDistributor{
	requests := make([][]config.RequestState, 4)
	for floor := range requests{
		requests[floor] = make([]config.RequestState, 3)
		
	}
	return config.ElevatorDistributor{Requests: requests, ID: id, Floor:0, Behaviour: config.Idle}

}

func broadcast(elevators []*config.ElevatorDistributor, ch_transmit chan <- []config.ElevatorDistributor){
	temporaryElevators := make([]config.ElevatorDistributor, 0)
	for _, elevator := range elevators{
		temporaryElevators = append(temporaryElevators, *elevator)
	}
	ch_transmit <- temporaryElevators
	time.Sleep(50*time.Millisecond)
}

func Distributor(id string, ch_newLocalOrder chan elevio.ButtonEvent, ch_newLocalState chan elevator.Elevator, ch_msgFromNetwork chan []config.ElevatorDistributor, ch_msgToNetwork chan []config.ElevatorDistributor, ch_orderToLocal chan elevio.ButtonEvent, ch_peerUpdate chan peers.PeerUpdate, ch_watchdogStuckReset bool , ch_watchdogStuckSignal chan bool, ch_clearLocalHallOrders chan bool){
	elevators := make([]*config.ElevatorDistributor, 0)
	thisElevator := new(config.ElevatorDistributor)
	*thisElevator = elevatorDistributorInit(id)
	elevators = append(elevators, thisElevator)

	connectTimer := time.NewTimer(time.Duration(config.ReconnectTimer)*time.Second)
	select{
	case newElevators := <- ch_msgFromNetwork:
		for _, elevator := range newElevators{
			if elev.id == elevators[localElevator].ID{
				for floor := range elev.Requests{
					if elev.Requests[floor][config.Cab] == config.Confirmed || elev.Requests[floor][config.Cab] == config.Order{
						ch_newLocalOrder <- elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(int(config.Cab))}
					}
				}
			}
		}
		break 
	case <- connectTimer.C:
		break
	}
	for{
		select{
		case newOrder := <- ch_newLocalOrder:
			assigner.AssignOrder(elevators, newOrder)
			if elevators[localElevator].Requests[newOrder.Floor][newOrder.Button] == config.Order{
				broadcast(elevators, ch_msgToNetwork)
				elevators[localElevator].Requests[newOrder.Floor][newOrder.Button] = config.Confirmed
				setHallLights(elevators)
				ch_orderToLocal <- newOrder
			}
			broadcast(elevators, ch_msgToNetwork)
			setHallLights(elevators)
		case newState := ch-newLocalState:
			if newState.Floor != elevators[localElevator].Floor || newState.Behave == elevator.Idle || newState.Behave == elevator.DoorOpen{
				elevators[localElevator].Behave = config.Behvaiour(int(newState.Behave))
				elevators[localElevator].Floor = newState.Floor
				elevators[localElevator].Direction = config.Direction(int(newState.Direction))
				ch_watchdogStuckReset <- false
			}
			for floor := range elevators[config.LocalElevator].Requests{
				for button := range elevators[config.LocalElevator].Requests[floor]{
					if !newState.Requests[floor][button] && elevators[config.LocalElevator].Requests[floor][button] == config.Confirmed{
						elevators[config.LocalElevator].Requests[floor][button] = config.Complete
					}
					if elevators[config.LocalElevator].Behave != config.Unavailable && newState.Requests[floor][button] && elevators[config.LocalElevator].Requests[floor][button] != config.Confirmed{
						elevators[config.LocalElevator].Requests[floor][button] = config.Confirmed
					}
				}
				
			}
			setHallLights(elevators)
			broadcast(elevators, ch_msgToNetwork)
			removeCompletedOrders(elevators)
			
		case newElevators := <- ch_msgFromNetwork:
			updateElevators(elevators,newElevator)
			assigner.ReassignOrders(elevators, ch_newLocalOrder)
			for _, newElev := range newElevators{
				elevExists := false
				for _, elev := range elevators{
					if elev.ID == newElev.ID{
						elevExists = true
						break
					}
				}
				if !elevExists{
					addNewElevator(&elevators,newElev)
				}
			}
			extractNewOrder := confirmNewOrder(elevators[localElevator])
			setHallLights(elevators)
			removeCompletedOrders(elevators)
			if extractNewOrder != nil{
				tempOrder := elevio.ButtonEvent{
					Button : elevio.ButtonType(extractNewOrder.Button),
					Floor : extractNewOrder.Floor}
				ch_orderToLocal <- tempOrder
				broadcast(elevators, ch_msgToNetwork)
			}
		case peer := <- ch_peerUpdate:
			if len(peer.Lost) != 0{
				for _, stringLostId := range peer.Lost{
					for _,elev := range elevators {
						if stringLostId == elev.ID{
							elev.Behaviour = config.Unavailable
						}
						assigner.ReassignOrders(elevators, ch_newLocalOrder)
						for floor := range elev.Requests{
							for button := 0; button < len(elev.Requests[floor])-1 ; button++{
								elev.Requests[floor][button] = config.None
							}
						}
					}
				}
			}
			setHallLights(elevators)
			broadcast(elevators, ch_msgToNetwork)
		case <- ch_watchdogStuckSignal:
			elevators[localElevator].Behave = config.Unavailable
			broadcast(elevators, ch_msgToNetwork)
			for floor := range elevators[localElevator].Requests{
				for button := 0; button < len(elevators[localElevator].Requests[floor])-1;button++{
					elevators[localElevator].Requests[floor][button] = config.None
				}
			}
			setHallLights(elevators)
			ch_clearLocalHallOrders <- true
		}
	}
}

func removeCompletedOrders(elevators []*config.ElevatorDistributor){
	for _, elev := range elevators{
		for floor := range elev.Requests{
			for button := range elev.Requests[floor]{
				if elev.Requests[floor][button] = config.None
			}
		}
	}
}

func chechLocalOrderComplete(elev *config.ElevatorDistributor, localElev elevator.Elevator){
	for floor := range elev.Requests{
		for button := range elev.Requests[floor]{
			if !localElev.Requests[floor][button] && elev.Requests[floor][button] == config.Confirmed{
				elev.Requests[floor][button] = config.Complete

			}
			if localElev.Requests[floor][button] && elev.Requests[floor][button] != config.Confirmed && elev.Behaviour != config.Unavailable {
				elev.Requests[floor][button] = config.Confirmed
			}
		}
	}
}


func updateElevators(elevators []*config.ElevatorDistributor, newElevators []config.ElevatorDistributor){
	if elevators[localElevator].ID != newElevators[localElevator].ID{
		for _,elev := range elevators{
			if elev.ID == newElevators[localElevator].ID{
				for floor := range elev.Requests{
					for button := range elev.Requests[floor]{
						if !(elev.Requests[floor][button] == config.Confirmed && newElevators[localElevator].Requests[floor][button] == config.Order){
							elev.Requests[floor][button] = newElevators[localElevator].Requests[floor][button]
						}
						elev.Floor = newElevators[localElevator].Floor
						elev.Direction = newElevators[localElevator].Direction
						elev.Behave = newElevators[localElevator].Behave
					}
				}
			}
		}
		for _, newElev := range newElevators{
			if newElev.ID == elevators[localElevator].ID{
				for floor := range newElev.Requests{
					for button := range newElev.Requests[floor]{
						if elevators[localElevator].Behave != config.Unavailable{
							if newElev.Requests[floor][button] == config.Order {
								(*elevators[localElevator]).Requests[floor][button] = config.Order
							}
						}
					}
				}
			}
		}
	}
}


func addNewElevator (elevators *[] config.ElevatorDistributor, newElevator config.ElevatorDistributor) {
	tempElev := new(config.ElevatorDistributor)
	*tempElev = elevatorDistributorInit(newElevator.ID)
	(*tempElev).Behave = newElevator.Behave
	(*tempElev).Direction = newElevator.Direction
	(*tempElev).Floor = newElevator.Floor
	
	for floor := range tempElev.Requests{
		for button := range tempElev.Requests[floor]{
			tempElev.Requests[floor][button] = newElevator.Requests[floor][button]
		}
	}
	*elevators = append(*elevators, tempElev)
}



func confirmedNewOrder(elev *config.ElevatorDistributor) *config.Requests{
	for floor := range elev.Requests {
		for button := 0 ; button < len(elev.Requests[floor]); button++{
			if elev.Requests[floor][button] == config.Order{
				elev.Requests[floor][button] = config.Comfirmed 
				tempOrder := new(config.Requests)
				*tempOrder = config.Requests{
					Floor: floor,
					Button: config.ButtonType(button)}
					return tempOrder
				}
			}
		}
	}
	return nil
}

func setHallLights(elevators []*config.ElevatorDistributor) {
	for button := 0 ; button < config.NumButtons - 1 ; button++{
		for floor := 0 ; floor < config.NumFloors ; floor++{
			isLight := false
			for _, elev := range elevators{
				if elev.Requests[floor][button] == config.Confirmed{
					isLight = true
				}
			}
			elevio.SetButtonLamp(elevio.ButtonType(button), floor, isLight)
		}
	}
}