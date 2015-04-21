package main

import(
	"net"
	"fmt"
	"os"

)

func handleConnection(conn net.Conn){
	var buffer [512]byte
	for{
		n, err := conn.Read(buffer[0:])
		if err != nil{
			fmt.Println("Error in read, from handleConnection")
		}
		fmt.Println(string(buffer[0:]))
		_, err2 := conn.Write(buffer[0:n])
		if err2 != nil{
			fmt.Println("Error in write, from handleConnection")
		}
	}
}

func main(){
	service := ":22222"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for{
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("Error in accept, from listen")
		}
		handleConnection(conn)
		conn.Close()
	}
	
}



func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
