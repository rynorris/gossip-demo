package main

import (
	"github.com/stefankopieczek/gossip/base"
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
	log.SetDefaultLogLevel(log.DEBUG)
	err := caller.Start()
	if err != nil {
		log.Warn("Failed to start caller: %v", err)
		return
	}

	log.Info("Listening for incoming requests...")
	tx := <-caller.tm.Requests()
	r := tx.Origin()
	log.Info("Received request: %v", r.Short())
	log.Debug("Full form:\n%v\n", r.String())

	// Send a 200 OK
	resp := base.NewResponse(
		"SIP/2.0",
		200,
		"OK",
		[]base.SipHeader{},
		"",
	)

	base.CopyHeaders("Via", tx.Origin(), resp)
	base.CopyHeaders("From", tx.Origin(), resp)
	base.CopyHeaders("To", tx.Origin(), resp)
	base.CopyHeaders("Call-Id", tx.Origin(), resp)
	base.CopyHeaders("CSeq", tx.Origin(), resp)
	resp.AddHeader(
		&base.ContactHeader{
			DisplayName: &caller.displayName,
			Address: &base.SipUri{
				User: &caller.username,
				Host: caller.host,
			},
		},
	)

	log.Info("Sending 200 OK")
	tx.Respond(resp)

	ack := <-tx.Ack()

	log.Info("Received ACK")
	log.Debug("Full form:\n%v\n", ack.String())
}
