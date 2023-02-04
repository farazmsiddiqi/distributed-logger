package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strings"
	"strconv"
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
	// buf := b.Buffer
	//Read incoming data into buffer
	_, err := connection.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading: ", err.Error())
		os.Exit(1)
	}

	//First message from node is its name (node1 captured by finding first index of space)
	node_name := string(buf[:strings.Index(string(buf[:]), " ")])
	event1 := string(buf[strings.Index(string(buf[:]), " ")+1:]) //TODO: make sure there is an event1 to print 

	//Prints the "timestamp - node1 connected" message
	fmt.Printf("%f - %s connected\n", float64(time.Now().UnixNano()) / 1000000000.0, node_name) // find a way to force stdout print
	//fmt.Fprintln(os.Stdout, event1[:strings.Index(event1, " ")], node_name, event1[strings.Index(event1, " ")+1:strings.Index(event1, " ")+1+64])
	if(event1 != ""){ //event1 exists
		fmt.Fprintln(os.Stdout, event1[:strings.Index(event1, " ")], node_name, event1[strings.Index(event1, " ")+1:strings.Index(event1, " ")+1+64])
	}

	f, _ := os.Create("/home/mkolla2/MP0/aux_logger")
	defer f.Close()

	//Repeatedly reads new events
	for {
		event_buf := make([]byte, 1024)
		_, err := connection.Read(event_buf)
		if err != nil {
			fmt.Printf("%f - %s disconnected", float64(time.Now().UnixNano()) / 1000000000.0, node_name)
			break
			//os.Exit(1) // we don't want to exit 
		} else {
			//expecting message from node as
			// "[time] [eventid]"
			event := string(event_buf)
			// space_ind := bytes.IndexByte(event_buf, byte(' '))
			space_ind := strings.Index(event, " ")
			current_time := float64(time.Now().UnixNano()) / 1000000000.0
			
			//LOG EVENT ARRIVAL DELAY in AUX_LOG
			float_generated_time, _ := strconv.ParseFloat(event[0:space_ind-1], 64)
			string_diff := strconv.FormatFloat(current_time - float_generated_time, 'E', -1, 64)
			event_arrival_delay := []string{node_name, string_diff, "\n"}
			
			f.WriteString(strings.Join(event_arrival_delay, " : "))
			f.Sync()

			fmt.Fprintln(os.Stdout, /*generated_time:*/event[0:space_ind-1], node_name, /*eventid:*/event[space_ind+1:space_ind+64])
		}
	}

	// Close the connection when you're done with it.
	connection.Close()
}