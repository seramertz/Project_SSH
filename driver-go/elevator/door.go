package elevator

import (
	"root/config"
	"root/elevio"
	"time"
)

type DoorState int

const (
	Closed DoorState = iota
	InCountDown
	Obstructed
)

func Door(
	doorClosedChannel 	chan<- bool,
	doorOpenChannel		<-chan bool,
	obstrucedChannel		chan<- bool,
) {

	elevio.SetDoorOpenLamp(false)
	obstructionChannel := make(chan bool)
	go elevio.PollObstructionSwitch(obstructionChannel)

	obstruction := false
	doorState := Closed
	timeCounter := time.NewTimer(time.Hour)
	timeCounter.Stop()

	for {
		select {
		case obstruction = <-obstructionChannel:
			if !obstruction && doorState == Obstructed {
				elevio.SetDoorOpenLamp(false)
				doorClosedChannel <- true
				doorState = Closed
			}
			if obstruction {
				obstrucedChannel <- true
			} else {
				obstrucedChannel <- false
			}

		case <-doorOpenChannel:
			if obstruction {
				obstrucedChannel <- true
			}
			switch doorState {
			case Closed:
				elevio.SetDoorOpenLamp(true)
				timeCounter = time.NewTimer(config.DoorOpenDuration)
				doorState = InCountDown
			case InCountDown:
				timeCounter = time.NewTimer(config.DoorOpenDuration)

			case Obstructed:
				timeCounter = time.NewTimer(config.DoorOpenDuration)
				doorState = InCountDown

			default:
				panic("Door state not implemented")
			}
		case <-timeCounter.C:
			if doorState != InCountDown {
				panic("Door state not implemented")
			}
			if obstruction {
				doorState = Obstructed
			} else {
				elevio.SetDoorOpenLamp(false)
				doorClosedChannel <- true
				doorState = Closed
			}
		}
	}
}
