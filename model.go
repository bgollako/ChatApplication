package main

import "net"

type Message struct {
	From string `json:"from"`
	Msg  string `json:"msg"`
}

type Announce struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type User struct {
	Username string
	ID       string
	Conn     net.Conn
}
