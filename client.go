package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

// registerClient - registers the client with the server
// establishes a tcp connection with the chat server
// Gets the username from terminal
// Registers the username with chat server
func registerClient() {
	conn, err := net.Dial("tcp", GAddress)
	if err != nil {
		log.Fatalln("Unable to connect to the server ", err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Connected to the server")
	fmt.Println("Please enter username for chat")
	scanner.Scan()
	username := scanner.Text()
	announce := &Announce{
		Username: username,
		ID:       uuid.NewString(),
	}
	data, err := json.Marshal(announce)
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(data)
	fmt.Println("Registered user", username, "with chat server")
	go sendMessages(announce, conn)
	go readMessages(announce, conn)
	ch := make(chan bool, 1)
	<-ch
}

// readMessages - reads messages from the connection and writes them to the terminal
// Must be launched in a separate go routine
func readMessages(announce *Announce, conn net.Conn) {
	for {
		data := make([]byte, 8192)
		n, err := conn.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("User has disconnected from chat server")
				fmt.Println("Shutting down client")
				os.Exit(0)
			}
			continue
		}

		data = data[:n]
		msg := &Message{}
		json.Unmarshal(data, msg)
		fmt.Println(msg.From, "says:")
		fmt.Println(msg.Msg)
	}
}

// sendMessages - sends messages from the terminal and sends them on the connection
// Must be launched in a separate go routine
func sendMessages(announce *Announce, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := strings.TrimSpace(scanner.Text())
		msg := Message{
			From: announce.Username,
			Msg:  text,
		}
		data, _ := json.Marshal(msg)
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Connection to chat server is down")
			fmt.Println("Shutting down client")
			os.Exit(0)
		}
	}
}
