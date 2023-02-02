package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	HOST = "172.22.156.102"
	TYPE = "tcp"
)

func main() {
	// main needs to take in a int arg for the port number
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Too few arguments")
		os.Exit(1)
	}

	PORT := os.Args[1]

	// listen on TCP port - accept connection from nodes
	listener, error := net.Listen(TYPE, HOST+":"+PORT)
	if error != nil {
		fmt.Fprintln(os.Stderr, "Error listening:", error.Error())
		os.Exit(1)
	}

	defer listener.Close()
	for {
		//Listen for a new connection
		connection, error := listener.Accept()
		if error != nil {
			fmt.Fprintln(os.Stderr, "Error accepting: ", error.Error())
			os.Exit(1)
		}

		//Send new connections to handler
		go handleRequest(connection)

		//TODO: potentially end process (ask TA)
	}

}

// print out all received events in format
// 1. time of event
// 2. name of node that generated event
// 3. event id

// special events: connected/disconnected
// print "1610688413.743385 - node1 connected" if node connects
// print "1610688452.211595 - node2 disconnected" if node disconnects
// / if TCP connection breaks for any reason - reading from connection caused error

func handleRequest(connection net.Conn) {
	buf := make([]byte, 1024)
	//Read incoming data into buffer
	_, err := connection.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading: ", err.Error())
		os.Exit(1)
	}

	//First message from node is its name 
	// TODO: index to first space after arg
	node_name := string(buf[5])

	//Prints the "timestamp - node1 connected" message
	fmt.Fprintln(os.Stdout, time.Now().Unix(), " - ", node_name, " connected")

	//Repeatedly reads new events
	for {
		event_buf := make([]byte, 1024)
		_, err := connection.Read(event_buf)
		if err != nil {
			fmt.Fprintln(os.Stdout, time.Now().Unix(), " - ", node_name, " disconnected")
			fmt.Fprintln(os.Stderr, "Error reading: ", err.Error())
			os.Exit(1)
		} else {
			//expecting message from node as
			// "[time] [eventid]"
			space_ind := bytes.IndexByte(event_buf, byte(' '))
			fmt.Fprintln(os.Stdout, /*time:*/event_buf[0:space_ind-1], " ", node_name, " ", /*eventid:*/event_buf[space_ind:])
		}
	}

	// Close the connection when you're done with it.
	connection.Close()
}
