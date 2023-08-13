package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

var GAddress = ":7777"
var GUserbase = make(map[string]*User)

// startListening - starts listening on the GAddress,
// and accepts incoming tcp connections.
func startListening() {
	listener, err := net.Listen("tcp", GAddress)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleNewConnection(conn)
	}
}

// handleNewConnection - handles new incoming user connection.
func handleNewConnection(conn net.Conn) {
	announce := Announce{}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		log.Println("Error while registering new user", err)
		return
	}

	data = data[:n]
	err = json.Unmarshal(data, &announce)
	if err != nil {
		log.Println("Error while unmarshalling new user data", err)
		return
	}

	user := &User{
		Username: announce.Username,
		ID:       announce.ID,
		Conn:     conn,
	}
	GUserbase[user.ID] = user
	log.Println("Registered new user", user.Username)

	monitorUserConn(user)
}

// monitorUserConn - monitors the given user's tcp onnection
// It keeps listening on the user's input tcp connection
// and responds to incoming messages
// Each incoming message is broadcast to all other tcp connections other than
// the one it received it on.
func monitorUserConn(user *User) {
	for {
		data := make([]byte, 8192)
		n, err := user.Conn.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(user.Username, "has disconnected")
				delete(GUserbase, user.ID)
				return
			}
			log.Println("Error while reading from user", user.Username, err)
			continue
		}

		data = data[:n]
		for _, u := range GUserbase {
			if u.ID != user.ID {
				u.Conn.Write(data)
			}
		}
		continue
	}
}
