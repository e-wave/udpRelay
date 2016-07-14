package main

import (
	"net"
  "fmt"
  "os"
)

// Constants
const PORT_A string = "7777"  // port for single client A
const PORT_B string = "8888"  // port for multiple clients B
const PACKET_SIZE int = 1024  // packets size sent and received

func main(){
	fmt.Println("Starting server")

	// Variables
	var client_A *net.UDPAddr
	// clients_B is a slice keeping all the clients B connected to the server
	clients_B := make([]*net.UDPAddr, 0)  
	buffer := make([]byte, PACKET_SIZE)  // buffer will receive the data from A
	cmd_from_B := make([]byte, 10)  // command sent by clients B CONNECT or DISCONNECT
	// UDPAddr used to create listeners
	sender,_ := net.ResolveUDPAddr("udp", ":" + PORT_A)
	receiver,_ := net.ResolveUDPAddr("udp", "127.0.0.1:" + PORT_B)

	// UDP listeners one on port A and one on Port B.
	conn_sender, err := net.ListenUDP("udp", sender)
	if err != nil{
		fmt.Println("sender listener ", err)
		os.Exit(0)
	}

	conn_receiver, err := net.ListenUDP("udp", receiver)
	if err != nil{
		fmt.Println("receiver listener ", err)
		os.Exit(0)
	}

 	defer conn_sender.Close()
	defer conn_receiver.Close()


	// Goroutine checking for new client B and appending them to the slice
	go func() { 
		for{
			n,raddr,err := conn_receiver.ReadFromUDP(cmd_from_B)
			if err != nil{
				fmt.Println(err)
				continue
			}
			switch string(cmd_from_B[0:n]) {
			    case "CONNECT":
			        clients_B = append(clients_B, raddr)
			    case "DISCONNECT":
			        for i,client := range clients_B{
			        	if client.String() == raddr.String(){
			        		clients_B = append(clients_B[:i], clients_B[i+1:]...)  // Removes the disconnected client form the clientB list
			        		break
			        	}
			        }
			}

		}
	}()


	// Loop waiting for data from client A and sending to all clients B
	for{
		_,raddr,err1 := conn_sender.ReadFromUDP(buffer)  // Gets data from client A
		
		if err1 != nil{
			fmt.Println(err1)
			continue
		}

		// if not from same client then refuse
		if client_A != nil && raddr.String() != client_A.String(){
			continue
		}
		if client_A == nil{
			client_A = raddr
		}

		//fmt.Println(string(buffer[0:n]))

		// Loop send data from client A to all clients B
		i := 0
		for i < len(clients_B){
			_,err2 := conn_receiver.WriteToUDP(buffer, clients_B[i])  
			if err2 != nil{
				fmt.Println(err2)
			}
			i++
		}
	}


	
}
