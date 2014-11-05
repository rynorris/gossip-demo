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
		port:        5060,
		transport:   "TCP",
	}

	// Callee parameters
	callee = &endpoint{
		displayName: "Stefan",
		username:    "stefan",
		host:        "PC4470",
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
		log.Warn("Failed to setup call: %v", err)
	}

	log.Info("Call successfully set up.  Exiting.")
}
