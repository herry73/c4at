package main

import (
	"log"
	"net"
)

const Port = "6969"
const safeMode = false

func safeRemoteAddr(conn net.Conn) string {
	if safeMode {
		return "[REDACTED]"
	} else {
		return conn.RemoteAddr().String()
	}
}

type MessageType int

const (
	ClientConnected MessageType = iota + 1
	DeleteClient
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

type Client struct {
	conn     net.Conn
	outgoing chan string
}

func server(messages chan Message) {
	conns := map[string]net.Conn{}
	for {
		msg := <-messages
		switch msg.Type {
		case ClientConnected:
			conns[msg.Conn.RemoteAddr().String()] = msg.Conn
		case DeleteClient:
			delete(conns, msg.Conn.RemoteAddr().String())
		case NewMessage:
			senderaddr := msg.Conn.RemoteAddr().String()
			for addr, conn := range conns {
				if addr != senderaddr {
					_, err := conn.Write([]byte(msg.Text))
					if err != nil {
						//left: remove connection from list
						log.Printf("Could not send data to %s: %s", safeRemoteAddr(conn), err)
					}
				}
			}
		}
	}
}

func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			messages <- Message{
				Type: DeleteClient,
				Conn: conn,
			}
			return
		}
		messages <- Message{
			Type: NewMessage,
			Text: string(buffer[0:n]),
			Conn: conn,
		}

	}
}

func main() {

	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		// handle error
		log.Fatalf("ERROR: Could not listen to port %s: %s\n", Port, err)

	}
	log.Printf("Listening to TCP Connections on Port %s ....\n", Port)

	messages := make(chan Message)
	go server(messages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			//handle error
			log.Printf("ERROR: Could not accept a connection: %s\n", err)
		}
		log.Printf("Connection accepted from %s\n", safeRemoteAddr(conn))
		messages <- Message{
			Type: ClientConnected,
			Conn: conn,
		}
		go client(conn, messages)
	}
}
