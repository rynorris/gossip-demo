package main

import (
	"fmt"

	"github.com/stefankopieczek/gossip/base"
	"github.com/stefankopieczek/gossip/log"
	"github.com/stefankopieczek/gossip/transaction"
)

type endpoint struct {
	// Sip Params
	displayName string
	username    string
	host        string

	// Transport Params
	port      uint16 // Listens on this port.
	transport string // Sends using this transport. ("tcp" or "udp")

	// Internal guts
	tm *transaction.Manager
}

func (e *endpoint) Start() error {
	tm, err := transaction.NewManager(e.transport, fmt.Sprintf("%v:%v", e.host, e.port))
	if err != nil {
		return err
	}

	e.tm = tm

	return nil
}

func (caller *endpoint) Invite(callee *endpoint) error {
	branch := "z9hG4bK.callbranch1"
	callid := base.CallId("thisisacall3")
	invite := base.NewRequest(
		base.INVITE,
		&base.SipUri{
			User: &callee.username,
			Host: callee.host,
		},
		"SIP/2.0",
		[]base.SipHeader{
			&base.ViaHeader{
				&base.ViaHop{
					ProtocolName:    "SIP",
					ProtocolVersion: "2.0",
					Transport:       caller.transport,
					Host:            caller.host,
					Port:            &caller.port,
					Params: base.Params{
						"branch": &branch,
					},
				},
			},
			&base.ToHeader{
				DisplayName: &callee.displayName,
				Address: &base.SipUri{
					User: &callee.username,
					Host: callee.host,
				},
			},
			&base.FromHeader{
				DisplayName: &caller.displayName,
				Address: &base.SipUri{
					User: &caller.username,
					Host: caller.host,
					UriParams: base.Params{
						"transport": &caller.transport,
					},
				},
			},
			&base.ContactHeader{
				DisplayName: &caller.displayName,
				Address: &base.SipUri{
					User: &caller.username,
					Host: caller.host,
				},
			},
			&base.CSeq{
				SeqNo:      1,
				MethodName: base.INVITE,
			},
			&callid,
			base.ContentLength(0),
		},
		"",
	)

	log.Info("Sending: %v", invite.Short())
	tx := caller.tm.Send(invite, fmt.Sprintf("%v:%v", callee.host, callee.port))
	for {
		select {
		case r := <-tx.Responses():
			log.Info("Received response: %v", r.Short())
			log.Debug("Full form:\n%v\n", r.String())
			switch {
			case r.StatusCode >= 300:
				// Call setup failed.
				return fmt.Errorf("callee sent negative response code %v.", r.StatusCode)
			case r.StatusCode >= 200:
				// Ack 200s manually.
				log.Info("Sending Ack")
				tx.Ack()
				return nil
			}
		case e := <-tx.Errors():
			log.Warn(e.Error())
			return e
		}
	}
}
