package main

import (
	"log"
	"net"
)

const Port = "6969"
const safeMode = true

func safeRemoteAddr(conn net.Conn) string {
	if safeMode {
		return "[REDACTED]"
	} else {
		return conn.RemoteAddr().String()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	message := []byte("Hallo meine Freunde! Kraut!\n")
	n, err := conn.Write(message)
	if err != nil {
		log.Printf("Could not write message to %s\n", safeRemoteAddr(conn))
		return
	}
	if n < len(message) {
		log.Printf("The message was not fully written %d/%d \n", n, len(message))
	}

}

func main() {

	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		// handle error
		log.Fatalf("ERROR: Could not listen to epic port %s: %s\n", Port, err)

	}
	log.Printf("Listening to TCP Connections on Port %s ....\n", Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			//handle error
			log.Println("ERROR: Could not accept a connection: %s\n", err)
		}
		log.Printf("Connection accepted from %s\n", safeRemoteAddr(conn))
		go handleConnection(conn)
	}
}
