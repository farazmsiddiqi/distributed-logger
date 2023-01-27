package main

import "fmt"


/* create a series of nodes that are capable of sending messages on a port*/

type Node struct {
	// port to route data to. This port number is likely that of the centralized logger.
	dest_port int
	// holds the raw input from generator.py.
	incoming_data string
}

func main() {
	
}