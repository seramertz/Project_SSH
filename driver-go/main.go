package main

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/fsm"
)


func main() {

	/*
	var port string
	flag.StringVar(&port, "port", "", "port of this peer")
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	elevio.Init("localhost:"+port, 4)
	*/

	numFloors := 4
	elevio.Init("localhost: 15657", numFloors)

	ch_newLocalOrder := make(chan elevio.ButtonEvent, 100)
	ch_clearLocalHallOrders := make(chan bool)
	ch_orderToLocal := make(chan elevio.ButtonEvent, 100)
	ch_newLocalState := make(chan elevator.Elevator, 100)
	ch_arrivedAtFloors := make(chan int)
	ch_obstruction := make(chan bool, 1)
	ch_timerDoor := make(chan bool)

	go elevio.PollFloorSensor(ch_arrivedAtFloors)
	go elevio.PollObstructionSwitch(ch_obstruction)
	go elevio.PollButtons(ch_newLocalOrder)

	go fsm.Fsm(ch_orderToLocal, ch_newLocalState, ch_clearLocalHallOrders, ch_arrivedAtFloors, ch_obstruction, ch_timerDoor)

	select {}
}
