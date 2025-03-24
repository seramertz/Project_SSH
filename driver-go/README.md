How to run our elevator
======================

Go to the folder driver-go
```
Start the elevator by typing "go run main.go -port=xxxx -id=x"
```
Add the portnumber and desired id to the elevator

Running just "go run main.go" starts an elevator with default port = 15657 and id=0

To change the amount of elevators or floors go to config and assign the wanted numbers. 

Summary
======================
This project contains software for controlling 'n' elevators actross 'm' floors. It is a peer-to-peer system with a fleeting master and UDP communication protocol. 


Elevio: 
---
Responsible for interfacing with the elevator hardware. It provides functions to control the elevator system. 

---
Elevator: 
---
Manages the states and behaviour of a single elevator. 

---
FSM:
---
Manages the state of the local elevator. Handles events as new orders, arriving at floors, door obstruction and door timer 
expiration. 

---
Config: 
---
Contains configurations constants and types used throughout the project. 

---
Distributer: 
---
Responsible for distributing orders among the elevators in the network. It handles communication between the local elevator and the other elevators, reassigns orders and updates the states of the elevators. 

---
Watchdog: 
---
Monitors the elevator system for signs of a stuck elevator. 

---

Assigner: 
---
Assigns orders to the different elevators by using the cost function. It also helps reassign orders if an elevator is lost. 

---
Network: 
---
Handling communication between different elevators. 