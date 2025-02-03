package elevator

import (
	"fmt"
	"root/config"
	"root/elevio"
	"time"
)

type ElevatorState struct {
	Obstructed bool
	Motorstop  bool
	Behaviour  Behaviour
	Floor      int
	Direction  Direction
}

type Behaviour int

const (
	Idle Behaviour = iota
	DoorOpen
	Moving
)

func (b Behaviour) ToString() string {
	return map[Behaviour]string{Idle: "idle", DoorOpen: "doorOpen", Moving: "moving"}[b]
}

func ElevatorFSM(
	newOrderChannel 		<-chan Orders,
	deliveredOrderChannel chan<- elevio.ButtonEvent,
	newStateChannel 		chan<- ElevatorState,
) {

	doorOpenChannel 		:= make(chan bool, 16)
	doorClosedChannel 	:= make(chan bool, 16)
	floorEnteredChannel 	:= make(chan int)
	obstructedChannel 	:= make(chan bool, 16)
	motorChannel 			:= make(chan bool, 16)

	go Door(doorClosedChannel, doorOpenChannel, obstructedChannel)
	go elevio.PollFloorSensor(floorEnteredChannel)

	elevio.SetMotorDirection(elevio.MD_Down)
	state := ElevatorState{Direction: Down, Behaviour: Moving}

	var orders Orders

	motorTimer := time.NewTimer(config.WatchdogTime)
	motorTimer.Stop()

	for {
		select {
		case <-doorClosedChannel:
			switch state.Behaviour {
			case DoorOpen:
				switch {
				case orders.OrderInDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(state.Direction.toMD())
					state.Behaviour = Moving
					motorTimer = time.NewTimer(config.WatchdogTime)
					motorChannel <- false
					newStateChannel <- state

				case orders[state.Floor][state.Direction.Opposite()]:
					doorOpenChannel <- true
					state.Direction = state.Direction.Opposite()
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)
					newStateChannel <- state

				case orders.OrderInDirection(state.Floor, state.Direction.Opposite()):
					state.Direction = state.Direction.Opposite()
					elevio.SetMotorDirection(state.Direction.toMD())
					state.Behaviour = Moving
					motorTimer = time.NewTimer(config.WatchdogTime)
					motorChannel <- false
					newStateChannel <- state

				default:
					state.Behaviour = Idle
					newStateChannel <- state
				}
			default:
				panic("DoorClosed in wrong state")
			}

		case state.Floor = <-floorEnteredChannel:
			elevio.SetFloorIndicator(state.Floor)
			motorTimer.Stop()
			motorChannel <- false
			switch state.Behaviour {
			case Moving:
				switch {
				case orders[state.Floor][state.Direction]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenChannel <- true
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)
					state.Behaviour = DoorOpen

				case orders[state.Floor][elevio.BT_Cab] && orders.OrderInDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenChannel <- true
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)
					state.Behaviour = DoorOpen

				case orders[state.Floor][elevio.BT_Cab] && !orders[state.Floor][state.Direction.Opposite()]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenChannel <- true
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)
					state.Behaviour = DoorOpen

				case orders.OrderInDirection(state.Floor, state.Direction):
					motorTimer = time.NewTimer(config.WatchdogTime)
					motorChannel <- false

				case orders[state.Floor][state.Direction.Opposite()]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenChannel <- true
					state.Direction = state.Direction.Opposite()
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)
					state.Behaviour = DoorOpen

				case orders.OrderInDirection(state.Floor, state.Direction.Opposite()):
					state.Direction = state.Direction.Opposite()
					elevio.SetMotorDirection(state.Direction.toMD())
					motorTimer = time.NewTimer(config.WatchdogTime)
					motorChannel <- false

				default:
					elevio.SetMotorDirection(elevio.MD_Stop)
					state.Behaviour = Idle
				}
			default:
				panic("FloorEntered in wrong state")
			}
			newStateChannel <- state

		case orders = <-newOrderChannel:
			switch state.Behaviour {
			case Idle:
				switch {
				case orders[state.Floor][state.Direction] || orders[state.Floor][elevio.BT_Cab]:
					doorOpenChannel <- true
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)
					state.Behaviour = DoorOpen
					newStateChannel <- state

				case orders[state.Floor][state.Direction.Opposite()]:
					doorOpenChannel <- true
					state.Direction = state.Direction.Opposite()
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)
					state.Behaviour = DoorOpen
					newStateChannel <- state

				case orders.OrderInDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(state.Direction.toMD())
					state.Behaviour = Moving
					newStateChannel <- state
					motorTimer = time.NewTimer(config.WatchdogTime)
					motorChannel <- false

				case orders.OrderInDirection(state.Floor, state.Direction.Opposite()):
					state.Direction = state.Direction.Opposite()
					elevio.SetMotorDirection(state.Direction.toMD())
					state.Behaviour = Moving
					newStateChannel <- state
					motorTimer = time.NewTimer(config.WatchdogTime)
					motorChannel <- false
				default:
				}

			case DoorOpen:
				switch {
				case orders[state.Floor][elevio.BT_Cab] || orders[state.Floor][state.Direction]:
					doorOpenChannel <- true
					OrderDone(state.Floor, state.Direction, orders, deliveredOrderChannel)

				}

			case Moving:

			default:
				panic("Orders in wrong state")
			}
		case <-motorTimer.C:
			if !state.Motorstop {
				fmt.Println("Lost motor power")
				state.Motorstop = true
				newStateChannel <- state
			}
		case obstruction := <-obstructedChannel:
			if obstruction != state.Obstructed {
				state.Obstructed = obstruction
				newStateChannel <- state
			}
		case motor := <-motorChannel:
			if state.Motorstop {
				fmt.Println("Regained motor power")
				state.Motorstop = motor
				newStateChannel <- state
			}
		}
	}
}
