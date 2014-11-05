package main

import (
	"github.com/stefankopieczek/gossip/log"
)

var (
	// Caller parameters
	caller = &endpoint{
		displayName: "Ryan",
		username:    "ryan",
		host:        "centosvm-rpn",
		addr:        "centosvm-rpn",
		port:        5060,
		transport:   "TCP",
	}

	// Callee parameters
	callee = &endpoint{
		displayName: "Stefan",
		username:    "stefan",
		host:        "PC4470",
		addr:        "PC4470",
		port:        5060,
		transport:   "TCP",
	}
)

func main() {
	err := caller.Start()
	if err != nil {
		log.Warn("Failed to start caller: %v", err)
		return
	}

	caller.Invite(callee)
}
