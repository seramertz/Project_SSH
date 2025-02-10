package main

import (
	"flag"

	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/fsm"
	"Driver-go/config"
	"Driver-go/distributor"
	"Driver-go/watchdog"
)


func main() {

	var port string
	flag.StringVar(&port, "port", "", "port of this peer")
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	
	numFloors := 4
	elevio.Init("localhost: 15657", numFloors)

	elevio.Init("localhost:"+port, 4)

	// Distributor channels
	ch_newLocalOrder := make(chan elevio.ButtonEvent, 100)

	// Communication between distributor and 'local elevator'
	ch_clearLocalHallOrders := make(chan bool)
	ch_orderToLocal := make(chan elevio.ButtonEvent, 100)
	ch_newLocalState := make(chan elevator.Elevator, 100)
	
	// Watchdog channels
	ch_watchdogStuckReset := make(chan bool)
	ch_watchdogStuckSignal := make(chan bool)

	// 'Local elevator' channels
	ch_arrivedAtFloors := make(chan int)
	ch_obstruction := make(chan bool, 1)
	ch_timerDoor := make(chan bool)



	go elevio.PollFloorSensor(ch_arrivedAtFloors)
	go elevio.PollObstructionSwitch(ch_obstruction)
	go elevio.PollButtons(ch_newLocalOrder)

	go fsm.Fsm(ch_orderToLocal, ch_newLocalState, ch_clearLocalHallOrders, ch_arrivedAtFloors, ch_obstruction, ch_timerDoor)


	go watchdog.Watchdog(config.ElevatorStuckTolerance, ch_watchdogStuckReset, ch_watchdogStuckSignal)

	go distributor.Distributor(id, ch_newLocalOrder, ch_newLocalState, ch_orderToLocal, ch_watchdogStuckReset, ch_watchdogStuckSignal, ch_clearLocalHallOrders )
	select {}
}
