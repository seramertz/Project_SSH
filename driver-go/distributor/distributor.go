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

func elevatorDistributorInit(id string) config.ElevatorDistributer{
	requests := make([][]config.RequestState, 4)
	for floor := range requests{
		requests[floor] = make([]config.RequestState, 3)
		
	}
	return config.ElevatorDistributer{Requests: requests, ID: id, Floor:0, Behaviour: config.Idle}

}

func broadcast(elevators []*config.ElevatorDistributer, ch_transmit chan <- []config.ElevatorDistributer){
	temporaryElevators := make([]config.ElevatorDistributer, 0)
	for _, elevator := range elevators{
		temporaryElevators = append(temporaryElevators, *elevator)
	}
	ch_transmit <- temporaryElevators
	time.Sleep(50*time.Millisecond)
}

func Distributor(id string, ch_newLocalOrder chan elevio.ButtonEvent, ch_newLocalState chan elevator.Elevator, ch_msgFromNetwork chan []config.ElevatorDistributer, ch_msgToNetwork chan []config.ElevatorDistributer, ch_orderToLocal chan elevio.ButtonEvent, ch_peerUpdate chan peers.PeerUpdate, ch_watchdogStuckReset bool , ch_watchdogStuckSignal chan bool, ch_clearLocalHallOrders chan bool){
	elevators := make([]*config.ElevatorDistributer, 0)
	thisElevator := new(config.ElevatorDistributer)
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
				for floor := range elevators[localNetwork].Requests{
					for button := 0; button < len(elevators[localElevator].Requests[floor])-1;button++{
						
					}
				}
			}
	}
	
