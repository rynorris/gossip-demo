package main

import (
	"time"

	"github.com/stefankopieczek/gossip/log"
)

var (
	// Caller parameters
	caller = &endpoint{
		displayName: "Ryan",
		username:    "ryan",
		host:        "192.168.0.2",
		port:        5060,
		transport:   "UDP",
	}

	// Callee parameters
	callee = &endpoint{
		displayName: "Ryan's PC",
		username:    "stefan",
		host:        "192.168.0.9",
		port:        5060,
		transport:   "UDP",
	}
)

func main() {
	log.SetDefaultLogLevel(log.DEBUG)
	err := caller.Start()
	if err != nil {
		log.Warn("Failed to start caller: %v", err)
		return
	}

	// Receive an incoming call.
	caller.ServeInvite()

	<-time.After(2 * time.Second)

	// Send the BYE
	caller.Bye(callee)
}
