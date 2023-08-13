package main

import (
	"flag"
	"fmt"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "server", "select client/server mode")
	flag.Parse()

	if mode == "server" {
		fmt.Println("Starting in server mode")
		startListening()
	} else if mode == "client" {
		fmt.Println("Starting in client mode")
		registerClient()
	} else {
		fmt.Println("Improper input, select mode as server/client")
	}
}
