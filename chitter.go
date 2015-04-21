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

func HandleConnection(conn net.Conn, m map[string]chan []byte) {
	b := bufio.NewReader(conn)
	msg := make(chan []byte)
	client_id := <-idAssignmentChan
	m[client_id] = msg
	go func() {
		for {
			conn.Write(<-msg)
		}
	}()
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			delete(m, client_id)
			conn.Close()
			break
		}
		message := string(line[0:])
		if strings.Index(message, ":") == -1 {
			for _, msg := range m {
				msg <- []byte(client_id + ": " + string(line))
			}
		} else {
			cmd := line[0:strings.Index(message, ":")]
			cmdstr := strings.Trim(string(cmd[0:]), " ")
			msgstr := string(line[strings.Index(message, ":")+1:])
			msgstr = strings.Trim(msgstr, " ")
			if cmdstr == "all" {
				for _, msg := range m {
					msg <- []byte(client_id + ": " + msgstr)
				}
			}

			if cmdstr == "whoami" {
				conn.Write([]byte(client_id + "\n"))
			} else {
				for id, msg := range m {
					if id == cmdstr {
						msg <- []byte(client_id + ": " + msgstr)
					}
				}
			}
		}

	}
}

func IdManager() {
	var i uint64
	for i = 0; ; i++ {
		idAssignmentChan <- strconv.FormatUint(i, 10)
	}
}

func main() {
	var m map[string]chan []byte = make(map[string]chan []byte)
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
		go HandleConnection(conn, m)
	}
}
