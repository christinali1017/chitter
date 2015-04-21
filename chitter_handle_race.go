package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var idAssignmentChan = make(chan string)

func HandleConnection(conn net.Conn, mchan chan map[string]chan []byte, newMsg chan []byte) {
	b := bufio.NewReader(conn)
	msg := make(chan []byte)
	m := make(map[string]chan []byte)
	client_id := <-idAssignmentChan
	m[client_id] = msg
	mchan <- m
	go func() {
		for {
			conn.Write(<-msg)
		}
	}()
	for {

		line, err := b.ReadBytes('\n')
		if err != nil {
			conn.Close()
			m[client_id] = nil
			mchan <- m
			break
		}
		newMsg <- []byte(client_id + ": " + string(line))
	}
}

func IdManager() {
	var i uint64
	for i = 0; ; i++ {
		idAssignmentChan <- strconv.FormatUint(i, 10)
	}
}

func main() {
	newConn := make(chan map[string]chan []byte)
	newMsg := make(chan []byte)
	var line []byte
	go func() {
		m := make(map[string]chan []byte)
		allMap := make(map[string]chan []byte)
		for {
			select {
			case m = <-newConn:
				for id, ch := range m {
					if ch == nil {
						delete(allMap, id)
					} else {
						allMap[id] = ch
					}
				}
			case line = <-newMsg:
				client_id := string(line[0:strings.Index(string(line), ":")])
				message := strings.Trim(string(line[strings.Index(string(line), ":")+1:]), " ")
				line = []byte(message)
				if strings.Index(message, ":") == -1 {
					fmt.Println("to all")
					for _, msg := range allMap {
						msg <- []byte(client_id + ": " + string(line))
					}
				} else {
					cmd := line[0:strings.Index(message, ":")]
					cmdstr := strings.Trim(string(cmd[0:]), " ")
					msgstr := string(line[strings.Index(message, ":")+1:])
					msgstr = strings.Trim(msgstr, " ")
					if cmdstr == "all" {
						for _, msg := range allMap {
							msg <- []byte(client_id + ": " + msgstr)
						}
					}

					if cmdstr == "whoami" {
						allMap[client_id] <- []byte(client_id + "\n")
					} else {
						for id, msg := range allMap {
							if id == cmdstr {
								msg <- []byte(client_id + ": " + msgstr)
							}
						}
					}
				}
			}
		}
	}()
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: chitter <port-number>\n")
		os.Exit(1)
		return
	}
	port := os.Args[1]
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't connect to port")
		os.Exit(1)
	}
	go IdManager()
	fmt.Println("Listening on port", os.Args[1])
	for {
		conn, _ := server.Accept()
		go HandleConnection(conn, newConn, newMsg)
	}
}
