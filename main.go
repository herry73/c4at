package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"
)

const (
	Port        = "6969"
	safeMode    = true
	MessageRate = 1.08
	BanLimit    = 60.0
)

func sensitive(message string) string {
	if safeMode {
		return "[REDACTED]"
	} else {
		return message
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
	bannedList := map[string]time.Time{}
	for {
		msg := <-messages
		switch msg.Type {
		case ClientConnected:
			addr := msg.Conn.RemoteAddr().(*net.TCPAddr)
			bannedAt, banned := bannedList[addr.IP.String()]
			now := time.Now()
			if banned {
				if time.Since(bannedAt).Seconds() >= BanLimit {
					delete(bannedList, addr.IP.String())
					banned = false
				}
			}
			if !banned {
				log.Printf("Client (%s) has been connected\n", sensitive(addr.String()))
				clients[msg.Conn.RemoteAddr().String()] = &Client{
					Conn:        msg.Conn,
					LastMessage: time.Now(),
				}
			} else {
				msg.Conn.Write([]byte(fmt.Sprintf("You are banned lil bro: %f secs left\n", BanLimit-(now.Sub(bannedAt).Seconds()))))
				msg.Conn.Close()
			}
		case ClientDisconnected:
			addrStr := msg.Conn.RemoteAddr().String()
			if _, ok := clients[addrStr]; ok {
				log.Printf("Client %s has been disconnected\n", sensitive(addrStr))
				delete(clients, addrStr)
			}
		case NewMessage:
			now := time.Now()
			Authoraddr := msg.Conn.RemoteAddr().String()

			author := clients[Authoraddr]
			if author != nil {
				if utf8.Valid([]byte(msg.Text)) {
					if now.Sub(clients[Authoraddr].LastMessage).Seconds() >= MessageRate {
						log.Printf("Client %s sent message: %s\n", sensitive(Authoraddr), msg.Text)
						author.LastMessage = now
						author.StrikeCounter = 0
						for _, client := range clients {
							if client.Conn.RemoteAddr().String() != Authoraddr {
								/*_, err :=*/ client.Conn.Write([]byte(msg.Text))
								// if err != nil {
								//  log.Printf("Could not send data to %s\n", sensitive(client.Conn.RemoteAddr().String()))
								// }
							}
						}
					} else {
						author.StrikeCounter += 1
						if author.StrikeCounter >= 10 {
							msg.Conn.Write([]byte("You are banned lil bro\n"))
							bannedList[msg.Conn.RemoteAddr().(*net.TCPAddr).IP.String()] = now
							author.Conn.Close()
						}
					}
				} else {
					author.StrikeCounter += 1
					if author.StrikeCounter >= 10 {
						msg.Conn.Write([]byte("You are banned lil bro\n"))
						bannedList[msg.Conn.RemoteAddr().(*net.TCPAddr).IP.String()] = now
						author.Conn.Close()
					}
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
			conn.Close()
			messages <- Message{
				Type: ClientDisconnected,
				Conn: conn,
			}
			return
		}
		text := string(buffer[0:n])
		messages <- Message{
			Type: NewMessage,
			Text: text,
			Conn: conn,
		}
	}
}

func main() {

	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		// handle error
		log.Fatalf("ERROR: Could not listen to port %s: %s\n", Port, sensitive(err.Error()))

	}
	log.Printf("Listening to TCP Connections on Port %s ....\n", Port)

	messages := make(chan Message)
	go server(messages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			//handle error
			log.Printf("ERROR: Could not accept a connection: %s\n", sensitive(err.Error()))
		}
		log.Printf("Connection accepted from %s\n", sensitive(conn.RemoteAddr().String()))
		messages <- Message{
			Type: ClientConnected,
			Conn: conn,
		}
		go client(conn, messages)
	}
}
