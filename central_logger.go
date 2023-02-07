package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strings"
	"strconv"
	"bufio"
)

const (
	TYPE = "tcp"
)

func main() {
	// main needs to take in a int arg for the port number
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Too few arguments")
		os.Exit(1)
	}

	HOST, _ := os.Hostname()
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
	buf, err := bufio.NewReader(connection).ReadBytes('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading: ", err.Error())
		os.Exit(1)
	}

	//First message from node is its name (node1 captured by finding first index of space)
	node_name := string(buf[:strings.Index(string(buf[:]), " ")])

	//Prints the "timestamp - node1 connected" message
	fmt.Printf("%f - %s connected\n", float64(time.Now().UnixNano()) / 1000000000.0, node_name) // find a way to force stdout print


	f, _ := os.Create("aux_logger" + "_" + node_name)
	// prints bandwidth from length of initial message from server + time 
	first_bandwidth_log := "bandwidth for " + node_name + ": " + strconv.Itoa(len(node_name)) + " " + time.Now().String() + "\n"
	f.WriteString(first_bandwidth_log)
	defer f.Close()

	//Repeatedly reads new events
	for {
		event_buf := make([]byte, 1024)
		_, err := connection.Read(event_buf)
		if err != nil {
			fmt.Printf("%f - %s disconnected\n", float64(time.Now().UnixNano()) / 1000000000.0, node_name)
			break
			//os.Exit(1) // we don't want to exit 
		} else {
			//expecting message from node as
			// "[time] [eventid]"
			event := string(event_buf)
			space_ind := strings.Index(event, " ")
			current_time := time.Now()
			float_current_time := float64(time.Now().UnixNano()) / 1000000000.0
			
			//LOG EVENT ARRIVAL DELAY in AUX_LOG
			float_generated_time, _ := strconv.ParseFloat(event[0:space_ind-1], 64)
			string_diff := strconv.FormatFloat(float_current_time - float_generated_time, 'E', -1, 64)
			bandwidth := len(string(event[0:space_ind-1]) + string(event[space_ind+1:space_ind+64]))
			event_arrival_delay := []string{node_name, string_diff, strconv.Itoa(bandwidth), current_time.String(), "\n"}
			
			f.WriteString(strings.Join(event_arrival_delay, ","))
			f.Sync()

			fmt.Fprintln(os.Stdout, /*generated_time:*/event[0:space_ind-1], node_name, /*eventid:*/event[space_ind+1:space_ind+64])
		}
	}

	// Close the connection when you're done with it.
	connection.Close()
}