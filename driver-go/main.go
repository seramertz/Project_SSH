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
    
    d := elevio.MD_Up

    for {
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            elevator.FsmRequestsButtonPress(a.Floor, a.Button)
            
        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            if a == numFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)
            elevator.FsmFloorArrival(a)
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
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
