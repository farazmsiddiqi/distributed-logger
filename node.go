package main

import (
	f "fmt"
	"os"
)


/* create a series of nodes that are capable of sending messages on a port*/

type Node struct {
	// port to route data to. This port number is likely that of the centralized logger.
	dest_port int
	// holds the raw input from generator.py.
	incoming_data string
}

func (node Node) Area() int {
	// TODO: add functionality
	// return nil
	return -1
}

func main() {

	args := os.Args[1:]
	f.Println(args)
	

	/* 
	python3 -u generator.py 0.1 | ./node node1 10.0.0.1 1234
	The first argument is the name of the node. 
	The second and third arguments are the address and port of the centralized logging server. 
	This should be the address of your VM running the centralized server (e.g., VM0) and the port.
	*/


	
}