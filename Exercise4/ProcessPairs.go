// Combined Program using UDP communication
package main

import (
	"fmt"
    "net"
    "os/exec"
    "strconv"
    "time"
)

const (
    sendAddr       = "255.255.255.255:20009"
    receiveAddr    = ":20009"
    heartbeatMsg   = "heartbeat"
    heartbeatSleep = 500
)

// Function to start a backup process that will become primary if needed.
func startBackupProcess() {
    err := exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
    if err != nil {
        fmt.Println("Failed to start backup process:", err)
    }
}

// Function to resolve UDP address and handle errors.
func resolveUDPAddr(addr string) *net.UDPAddr {
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        fmt.Println("Failed to resolve UDP address:", err)
        return nil
    }
    return udpAddr
}

// Function to handle sending heartbeats.
func sendHeartbeat(conn *net.UDPConn, count int) error {
    msg := heartbeatMsg + ":" + strconv.Itoa(count)
    _, err := conn.Write([]byte(msg))
    if err != nil {
        return fmt.Errorf("failed to send heartbeat: %w", err)
    }
    fmt.Printf("Primary count: %d\n", count)
    return nil
}

// The primary process sends heartbeats to the backup.
func primaryProcess(count int) {
    sendUDPAddr := resolveUDPAddr(sendAddr)
    if sendUDPAddr == nil {
        return
    }
    conn, err := net.DialUDP("udp", nil, sendUDPAddr)
    if err != nil {
        fmt.Println("Failed to dial UDP:", err)
        return
    }
    defer conn.Close()

    for {
        if err := sendHeartbeat(conn, count); err != nil {
            fmt.Println(err)
            return
        }
        count++
        time.Sleep(heartbeatSleep * time.Millisecond)
    }
}

// The backup process listens for heartbeats from the primary.
func backupProcess() {
    count := 1
    fmt.Printf("---------BACKUP PHASE---------\n")
    receiveUDPAddr := resolveUDPAddr(receiveAddr)
    if receiveUDPAddr == nil {
        return
    }
    conn, err := net.ListenUDP("udp", receiveUDPAddr)
    if err != nil {
        fmt.Println("Failed to listen on UDP:", err)
        return
    }
    defer conn.Close()

    for {
        buffer := make([]byte, 1024)
        conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 5 * time.Millisecond))
        n, _, err := conn.ReadFromUDP(buffer)

        if err != nil {
            if e, ok := err.(net.Error); ok && e.Timeout() {
                fmt.Println("Backup did not receive heartbeat, becoming primary.")
                // This is where the backup takes over and becomes Primary
                conn.Close()
                startBackupProcess()
                primaryProcess(count + 1)
                return
            } else {
                fmt.Println("Error reading from UDP:", err)
                return
            }
        }

        msg := string(buffer[:n])
        if msg[:len(heartbeatMsg)] == heartbeatMsg {
            countStr := msg[len(heartbeatMsg)+1:]
            recievedCount, _ := strconv.Atoi(countStr)
            count = recievedCount
        }
    }
}

func main() {
    backupProcess()
}