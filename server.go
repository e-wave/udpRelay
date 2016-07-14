package main

import (
	"net"
	"fmt"
	"os"
)

// Constants
const(
	PORT_A string = "7777"  // port for single client A
	PORT_B string = "8888"  // port for multiple clients B
	PACKET_SIZE int = 1024  // packets size sent and received
)


func create_connection(ip, port string) *net.UDPConn{
	// UDPAddr used to create listener
	udpaddr,_ := net.ResolveUDPAddr("udp", ip + ":" + port)

	// UDP listeners one on port A and one on Port B.
	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil{
		fmt.Println("listener ", err)
		os.Exit(0)
	}

	return conn
}

func receive_bytes(conn_sender *net.UDPConn, messages chan []byte){
	// Variables
	var client_A *net.UDPAddr
	buffer := make([]byte, PACKET_SIZE)  // buffer will receive the data from A

	// Loop waiting for data from client A and sending to all clients B
	for{
		n,raddr,err1 := conn_sender.ReadFromUDP(buffer)  // Gets data from client A
		
		if err1 != nil{
			fmt.Println(err1)
			continue
		}

		// If not from same client then refuse
		if client_A != nil && raddr.String() != client_A.String(){
			continue
		}
		if client_A == nil{
			client_A = raddr
		}

		messages <- buffer[0:n]
		fmt.Println(string(buffer[0:n]))
	}
}

func send_to_B(conn_receiver *net.UDPConn, clients_B map[string]*net.UDPAddr, messages chan []byte){
	for{
		msg := <-messages
		// Loop send data from client A to all clients B
		i := 0
		for key, val := range clients_B{
			_,err2 := conn_receiver.WriteToUDP(msg, val)  
			if err2 != nil{
				fmt.Println(key, err2)
			}
			i++
		}
	}
}

func main(){
	fmt.Println("Starting server")

	// Variables
	clients_B := make(map[string]*net.UDPAddr) // map of all the connected B  clients 
	cmd_from_B := make([]byte, 10)  // command sent by clients: B CONNECT or DISCONNECT
	messages := make(chan []byte)

	// UDP listeners one on port A and one on Port B.
	conn_sender := create_connection("", PORT_A)
	conn_receiver := create_connection("", PORT_B)

	// Close connection when exiting main thread
 	defer conn_sender.Close()
	defer conn_receiver.Close()

	go receive_bytes(conn_sender, messages)
	go send_to_B(conn_receiver, clients_B, messages)

	// Checking for new client B and adding them to the map
	// map clients_B keys are string representation of the udpaddr 
	// and the value are the udpaddr object
	for{
		n,raddr,err := conn_receiver.ReadFromUDP(cmd_from_B)
		if err != nil{
			fmt.Println(err)
			continue
		}
		switch string(cmd_from_B[0:n]) {
		    case "CONNECT":
		        clients_B[raddr.String()] = raddr   // Adds the new connected client to the map
		    case "DISCONNECT":
		        delete(clients_B, raddr.String()) // Removes the disconnected client form the clientB list
		}
		//fmt.Println(clients_B)
	}
	
}
