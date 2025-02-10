package fsm

import (
	"Driver-go/config"
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/request"
	"fmt"
	"time"
)

//Statemachine for running the main elevator
func Fsm(ch_orderChan chan elevio.ButtonEvent,ch_elevatorState chan <- elevator.Elevator,ch_clearLocalHallOrders chan bool,
	ch_arrivedAtFloors chan int,ch_obstruction chan bool,ch_timerDoor chan bool){

		elev := elevator.InitElevator()
		e := &elev
		
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevio.MD_Down)

		elevator.ElevatorPrint(*e)

		//Initialize at floor zero
		for{
			floor := <-ch_arrivedAtFloors
			if floor != 0{
				elevio.SetMotorDirection(elevio.MD_Down)
			} else{
				elevio.SetMotorDirection((elevio.MD_Stop))
				break
			}
		}
		
		ch_elevatorState <- *e

		doorTimer := time.NewTimer(time.Duration(config.DoorOpenDuration) * time.Second)
		timerUpdateState := time.NewTicker(time.Duration(config.StateUpdatePeriodsMs) * time.Millisecond)
		
		//Statemachine defining the elevators state
		for{
			fmt.Printf("in for loop")
			elevator.LightsElevator(*e)
			select{
			case order := <-ch_orderChan: //an order is placed
				fmt.Printf("in for order")
				switch {
					case e.Behave == elevator.DoorOpen:
						if e.Floor == order.Floor{
							doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
						} else{
							e.Requests[order.Floor][order.Button] = true
						}
					case e.Behave == elevator.Moving:
						e.Requests[order.Floor][order.Button] = true
					case e.Behave == elevator.Idle:
						if e.Floor == order.Floor{
							elevator.LightsElevator(*e)
							elevio.SetDoorOpenLamp(true)
							doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
							e.Behave = elevator.DoorOpen
							ch_elevatorState <- *e
						} else{
							e.Requests[order.Floor][int(order.Button)] = true
							request.RequestChooseDirection(e)
							elevio.SetMotorDirection(e.Direction)
							e.Behave = elevator.Moving
							ch_elevatorState <- *e
							break
						}
				}
			case floor := <-ch_arrivedAtFloors: //elevator has reached a floor
				fmt.Printf("in for floor")
				e.Floor = floor
				switch{
					case e.Behave == elevator.Moving:
						if request.RequestShouldStop(e){
							elevio.SetMotorDirection(elevio.MD_Stop)
							elevator.LightsElevator(*e)
							request.RequestClearAtCurrentFloor(e)
							elevio.SetDoorOpenLamp(true)
							doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
							e.Behave = elevator.DoorOpen
							ch_elevatorState <- *e
					
						}
				default:	
					break
					
				}
			case <-doorTimer.C: //door is open and timer is counting
				fmt.Printf("in for in doortimer")
				switch{
					case e.Behave == elevator.DoorOpen:
						request.RequestChooseDirection(e)
						elevio.SetMotorDirection(e.Direction)
						elevio.SetDoorOpenLamp(false)

						if e.Direction == elevio.MD_Stop{
							e.Behave = elevator.Idle
							ch_elevatorState <- *e
						} else{
							e.Behave = elevator.Moving
							ch_elevatorState <- *e
						}
					default:	
						break
				}
			case <-ch_clearLocalHallOrders: //delete the hallorders of this elevator
				fmt.Printf("in for clear local hall orders")
				request.RequestClearHall(e)
			case obstruction := <-ch_obstruction: //obstruction button 
				if e.Behave == elevator.DoorOpen && obstruction{
					doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
				}
			case <-timerUpdateState.C: //if the time is out
				ch_elevatorState <- *e
				timerUpdateState.Reset(time.Duration(config.StateUpdatePeriodsMs) * time.Millisecond)
				
			}	
	}
}