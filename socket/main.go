package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	// use a command line flag to choose if we are the server or client
	server := flag.Bool("server", false, "Run as a server")
	flag.Parse()

	// the argument is our name, the second argument is the host:port to connect to
	// if we are the server you can just use localhost:port
	var conn net.Conn
	var err error
	if *server {
		fmt.Println("Running as a server!")

		listener, err := net.Listen("tcp", flag.Args()[1])

		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}

		// listen for an incoming connection
		conn, err = listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// close the listener when the application closes
		defer listener.Close()
	} else {
		fmt.Println("Running as a client!")

		conn, err = net.Dial("tcp", flag.Args()[1])
		if err != nil {
			fmt.Println("Error dialing:", err.Error())
			os.Exit(1)
		}

		// close the connection when the application closes
		defer conn.Close()
	}

	s := NewSocket(conn)

	// Create a thread to handle sending messages
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// add name: to the start of the message
			s.send <- []byte(fmt.Sprintf("%s: %s", flag.Args()[0], scanner.Text()))
		}
	}()

	for {
		msg := <-s.recv
		fmt.Println(string(msg))
	}
}
