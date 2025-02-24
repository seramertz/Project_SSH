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
func Fsm(
	ch_orderChan chan elevio.ButtonEvent,
	ch_elevatorState chan <- elevator.Elevator,
	ch_clearLocalHallOrders chan bool,
	ch_arrivedAtFloors chan int,
	ch_obstruction chan bool,
	ch_timerDoor chan bool){

		elev := elevator.InitElevator()
		e := &elev
		
		elevio.SetDoorOpenLamp(false)

		elevator.ElevatorPrint(*e)

		ch_elevatorState <- *e

		doorTimer := time.NewTimer(time.Duration(config.DoorOpenDuration) * time.Second)
		timerUpdateState := time.NewTicker(time.Duration(config.StateUpdatePeriodsMs) * time.Millisecond)
		
		obstructionActive := false
		//Statemachine defining the elevators state
		for{
			elevator.LightsElevator(*e)
			select{
			case order := <-ch_orderChan: //an order is placed
				//fmt.Printf("in for order")
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
							
							elevio.SetDoorOpenLamp(true)
							doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
							e.Behave = elevator.DoorOpen
							ch_elevatorState <- *e
						} else{
							e.Requests[order.Floor][order.Button] = true
							request.RequestChooseDirection(e)
							elevio.SetMotorDirection(e.Direction)
							e.Behave = elevator.Moving
							ch_elevatorState <- *e
							break
						}
				}
			case floor := <-ch_arrivedAtFloors: //elevator has reached a floor
				
				e.Floor = floor
				switch{
					case e.Behave == elevator.Moving:
						if request.RequestShouldStop(e){
							elevio.SetMotorDirection(elevio.MD_Stop)
							request.RequestClearAtCurrentFloor(e)
							elevio.SetDoorOpenLamp(true)
							doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
							e.Behave = elevator.DoorOpen
							ch_elevatorState <- *e
							
							
							// Handle obstruction if active
							if obstructionActive {
								fmt.Printf("Obstruction detected: %v\n", obstructionActive)
								doorTimer.Stop()
								for obstructionActive {
									obstructionActive = <-ch_obstruction
								}
								
								fmt.Printf("Obstruction cleared: %v\n", obstructionActive)
								doorTimer = time.NewTimer(time.Duration(config.DoorOpenDuration) * time.Second)
							}else{
								obstructionActive = false
							}
						}
				default:	
					break
					
				}
			case <-doorTimer.C: //door is open and timer is counting
			
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
				
				request.RequestClearHall(e)

			case obstruction := <-ch_obstruction: //obstruction button 
				if obstruction {
					obstructionActive = true
					if e.Behave == elevator.DoorOpen {
						fmt.Printf("Obstruction detected: obstruction =  %v\n", obstruction)
						doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)

						// Handle obstruction while door is open
						for obstruction {
							obstruction = <-ch_obstruction
						}
					}
				} else {
					obstructionActive = false
				}
			fmt.Printf("Obstruction cleared: obstruction = %v\n", obstruction)
			doorTimer = time.NewTimer(time.Duration(config.DoorOpenDuration) * time.Second)
			case <-timerUpdateState.C: //if the time is out
				ch_elevatorState <- *e
				timerUpdateState.Reset(time.Duration(config.StateUpdatePeriodsMs) * time.Millisecond)
				
			}	
	}
}