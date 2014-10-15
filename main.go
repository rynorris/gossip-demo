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
		username:    "ryan",
		host:        "192.168.0.9",
		port:        5060,
		transport:   "TCP",
	}
)

func main() {
	log.SetDefaultLogLevel(log.INFO)
	err := caller.Start()
	if err != nil {
		log.Warn("Failed to start caller: %v", err)
		return
	}

	err = caller.Invite(callee)
	if err != nil {
		log.Warn("Failed to start call: %v", err)
		return
	}

	// Send a BYE after 3 seconds.
	<-time.After(5 * time.Second)
	err = caller.Bye(callee)
	if err != nil {
		log.Warn("Failed to end call: %v", err)
		return
	}
}
