package main
 
import (
  "fmt"
  "net"
  "time"
  "os"
)

const(
	PACKET_SIZE int = 1024
	DROP_LIMIT float64 = 0.2  // 20%
)


func is_transmission_drop(data_amount int, data_amount_old int, uptime int) bool{
	avg_total:= float64(data_amount)/float64(uptime)
	avg_last_10_sec := (float64(data_amount) - float64(data_amount_old))/10.0

	if  avg_total > avg_last_10_sec{
		drop_rate := (avg_total - avg_last_10_sec) / avg_total
		if drop_rate >= DROP_LIMIT{
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Starting clientB")

	//Variables
	buffer := make([]byte, PACKET_SIZE)
	uptime := 0  // total uptime
	data_total := 0  // total amount of data received

	// UDPAddr used to create connection
	server,_ := net.ResolveUDPAddr("udp", "127.0.0.1:8888")  // Send to server on port 8888
	local, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")     // Receive from server on random port
    conn, err := net.DialUDP("udp", local, server)  // Create connection to server
	if err != nil{
		fmt.Println("dial ",err)
	}

	defer conn.Close()

	// Send command connect to server
	fmt.Println("Sending CONNECT cmd to server") 
	cmd_for_server := ("CONNECT")
	buf := []byte(cmd_for_server)
    conn.Write(buf)

    // Goroutine checking the amount of data received every 10 seconds
    // In case of a drop in the transmission rate 
    // it sends disconnect command to server to stop receiving data
    go func() { 
		for{
			data_total_old := data_total
			time.Sleep(time.Second * 10)
			uptime += 10

			if is_transmission_drop(data_total, data_total_old, uptime){
				fmt.Println("Sending DISCONNECT cmd to server") 
				cmd_for_server := ("DISCONNECT")
				buf := []byte(cmd_for_server)
			    conn.Write(buf)
			    os.Exit(0)
			}
			
		}
	}()

	// Receive data from client A through server
	for {
		n,_,err := conn.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println(err)
			continue
		}
		data_total += n
		
		fmt.Println(string(buffer[0:n]))
  	}

}
