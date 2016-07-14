package main
 
import (
  "fmt"
  "net"
  "time"
)

const PACKET_SIZE int = 1024
 
func init_buffer(buffer *[]byte) int{
  i := 0
  for i < len(*buffer){
    (*buffer)[i] = 42  // ascii 42 = *
    i++
  }

  return 0
}
 
func main() {
	fmt.Println("Starting clientA")

  buffer := make([]byte, PACKET_SIZE)

	server,_ := net.ResolveUDPAddr("udp", "127.0.0.1:7777") // Send to server on port 7777
  conn, err := net.DialUDP("udp", nil, server) // Create connection to server
	if err != nil{
		fmt.Println(err)
	}

  defer conn.Close()

  init_buffer(&buffer)
  for {
    _,err := conn.Write(buffer)  // Send packets to server every 100ms
    if err != nil {
        fmt.Println(err)
    }
    time.Sleep(time.Millisecond * 100)
  }
}
