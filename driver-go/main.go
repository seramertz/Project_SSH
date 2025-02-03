package main

import "Driver-go/elevio"
import "Driver-go/elevator"
import "fmt"
import "time"

func main(){

    numFloors := 4
    pollRate := 20 * time.Millisecond

    //elevio.Init("localhost:15657", numFloors)
    
    //var d elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)
    elevator.FsmInit(numFloors)

    prevRequestButtonPress := make([][]int, numFloors)

    for i := range prevRequestButtonPress {
        prevRequestButtonPress[i] = make([]int, elevator.NumButtons)
    }
    
    //prevFloorSensor := -1

    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    
    dirn := elevio.MD_Up
    

    for {
        select {
        case buttonPressed := <- drv_buttons:
            fmt.Printf("%+v\n", buttonPressed)
            elevio.SetButtonLamp(buttonPressed.Button, buttonPressed.Floor, true)
            elevator.FsmRequestsButtonPress(buttonPressed.Floor, buttonPressed.Button)
            
        case floorSensor := <- drv_floors:
            fmt.Printf("%+v\n", floorSensor)
            if floorSensor == numFloors-1 {
                dirn = elevio.MD_Down
            } else if floorSensor == 0 {
                dirn = elevio.MD_Up
            }
            elevio.SetMotorDirection(dirn)
            elevator.FsmFloorArrival(floorSensor)
            
            
        case obstruction := <- drv_obstr:
            fmt.Printf("%+v\n", obstruction)
            if obstruction {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(dirn)
            }
            
        case stopPressed := <- drv_stop:
            fmt.Printf("%+v\n", stopPressed)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        default: 
            select {
            case <-time.After(pollRate):
                elevator.FsmDoorTimeout()
            default:
                time.Sleep(pollRate)
            }
            time.Sleep(pollRate)
        }    
    }
    
}
