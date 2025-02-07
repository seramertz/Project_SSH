package fsm

import (
	"Driver-go/config"
	"Driver-go/elevator"
	"Driver-go/elevio"
	"driver-go/config"
	"driver-go/elevator"
	"driver-go/elevio"
	"driver-go/request"
	"time"
)

func fsm(ch_orderChan chan elevio.ButtonEvent,ch_elevatorState chan <- elevator.Elevator,ch_clearLocalHallOrders chan bool,
	ch_arrivedAtFloors chan int,ch_obstruction chan bool,ch_timerDoor chan bool){

		elevator := elevator.InitElevator()
		e := &elevator

		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevio.MD_Down)

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

		for{
			elevator.LightsElevator(*e)
			select{
			case order := <-ch_orderChan:
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
							request.RequestOrder(e)
							elevio.SetMotorDirection(e.Direction)
							e.Behave = elevator.Moving
							ch_elevatorState <- *e
							break
						}
				}
			case floor := <-ch_arrivedAtFloors:
				e.Floor = floor
				switch{
					case e.Behave == elevator.Moving:
						if request.ShouldStop(e){
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
			case <-doorTimer.C:
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
			case <-ch_clearLocalHallOrders:
				request.RequestClearHall(e)
			case obstruction := <-ch_obstruction:
				if e.Behave == elevator.DoorOpen && obstruction{
					doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
				}
			case <-timerUpdateState.C:
				ch_elevatorState <- *e
				timerUpdateState.Reset(time.Duration(config.StateUpdatePeriodsMs) * time.Millisecond)
				
	}
