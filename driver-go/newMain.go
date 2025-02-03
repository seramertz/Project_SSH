package main

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/elevio"
)


func main() {

	elevio.Init("localhost:"+strconv.Itoa(Port), config.NumFloors)
  
	fmt.Println("Elevator initialized with id", id, "on port", Port)
	fmt.Println("System has", config.NumFloors, "floors and", config.NumElevators, "elevators.")

	newOrderChannel 		:= make(chan elevator.Orders, config.Buffer)
	deliveredOrderChannel	:= make(chan elevio.ButtonEvent, config.Buffer)
	newStateChannel 		:= make(chan elevator.State, config.Buffer)

	go elevator.Elevator(
		newOrderC,
		deliveredOrderC,
		newStateC)

  elevator.elevatorFSM(newOrderChannel, deliveredOrderChannel, newStateChannel)
	
	
}
