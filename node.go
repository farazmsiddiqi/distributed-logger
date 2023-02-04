package main

import (
	f "fmt"
	"bufio"
	"os"
	"net"
)

/*
% python3 generator.py 0.1
1610688413.782391 ce783874ba65a148930de32704cd4c809d22a98359f7aed2c2085bc1bd10f096
<generation timestamp> <event>
*/

type Node struct {
	// timestamp as a fractional number of seconds since 1970
	generationTimestamp int
	// event from the generator (to this node) that needs to be communicated to the logger
	event string
}

func (node Node) sendEvent(host string, port string, nodeName string) int {
	/*
	first_loop (on connnection): Central Logger expects first message from a new node connection to be its name "node1"
	second_loop (on new msg): Subsequent messages should be formatted as "[time] [eventid]"
	*/

	conn, err := net.Dial("tcp", host+":"+port)

	// send node name first
	conn.Write([]byte(f.Sprintf("%s", nodeName)))

	if err != nil {
		f.Fprintln(os.Stderr, "fatal err: %s", err.Error())
		os.Exit(1)
	}

	// send generator data to logger via tcp conn
	sendGenData(conn)

	return -1
}


func sendGenData(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')

		if err != nil {
			f.Fprintln(os.Stderr, "fatal err: %s", err.Error())
			os.Exit(1)
		}
		conn.Write([]byte(f.Sprintf("%s\r\n", text)))
	}

}


func main() {
	/* 
	python3 -u generator.py 0.1 | ./node node1 10.0.0.1 1234
	The first argument is the name of the node. 
	The second and third arguments are the address and port of the centralized logging server. 
	This should be the address of your VM running the centralized server (e.g., VM0) and the port.
	*/

	if len(os.Args) < 4 {
		f.Fprintln(os.Stderr, "too many arguments")
		os.Exit(1)
	}

	nodeName := os.Args[1] + " "
	host := os.Args[2]
	port := os.Args[3]

	var node Node
	node.sendEvent(host, port, nodeName)
}

