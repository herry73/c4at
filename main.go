package main

import (
	"log"
	"net"
	"time"
)

const (
	Port        = "6969"
	safeMode    = true
	MessageRate = 0.5
)

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
	ClientDisconnected
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

type Client struct {
	Conn          net.Conn
	LastMessage   time.Time
	StrikeCounter int
}

func server(messages chan Message) {
	clients := map[string]*Client{}
	for {
		msg := <-messages
		switch msg.Type {
		case ClientConnected:
			log.Printf("Client (%s) has been connected", safeRemoteAddr(msg.Conn))
			clients[msg.Conn.RemoteAddr().String()] = &Client{
				Conn:        msg.Conn,
				LastMessage: time.Now(),
			}
		case ClientDisconnected:
			log.Printf("Client %s has been disconnected", safeRemoteAddr(msg.Conn))
			delete(clients, msg.Conn.RemoteAddr().String())
		case NewMessage:
			now := time.Now()
			addr := msg.Conn.RemoteAddr().String()
			author := clients[addr]
			if now.Sub(clients[addr].LastMessage).Seconds() >= MessageRate {
				log.Printf("Client %s sent message %s", safeRemoteAddr(msg.Conn), msg.Text)
				author.LastMessage = now
				author.StrikeCounter = 0
				for _, client := range clients {
					if client.Conn.RemoteAddr().String() != addr {
						_, err := client.Conn.Write([]byte(msg.Text))
						if err != nil {
							log.Printf("")
						}
					}
				}
			} else {
				author.StrikeCounter += 1
				if author.StrikeCounter >= 10 {

				}
			}
		}
	}
}

func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 64)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Could not read from %s: %s", safeRemoteAddr(conn), err)
			conn.Close()
			messages <- Message{
				Type: ClientDisconnected,
				Conn: conn,
			}
			return
		}
		text := string(buffer[0:n])
		if text == ":quit" {

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
