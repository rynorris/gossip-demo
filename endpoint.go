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
	addr      string // Listens on and sends from this address.
	port      uint16 // Listens on this port.
	transport string // Sends using this transport. ("tcp" or "udp")

	// Internal guts
	tm *transaction.Manager
}

func (e *endpoint) Start() error {
	tm, err := transaction.NewManager(e.transport, fmt.Sprintf("%v:%v", e.addr, e.port))
	if err != nil {
		return err
	}

	e.tm = tm

	return nil
}

func (caller *endpoint) Invite(callee *endpoint) {
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
					Host:            caller.addr,
					Port:            &caller.port,
					Params: base.Params{
						"branch": &branch,
					},
				},
			},
			&base.ToHeader{
				Address: &base.SipUri{
					User: &callee.username,
					Host: callee.host,
				},
			},
			&base.FromHeader{
				Address: &base.SipUri{
					User: &caller.username,
					Host: caller.host,
					UriParams: base.Params{
						"transport": &caller.transport,
					},
				},
			},
			&base.ContactHeader{
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

	tx := caller.tm.Send(invite, fmt.Sprintf("%v:%v", callee.host, callee.port))
	for {
		select {
		case r := <-tx.Responses():
			log.Info(r.String())
			if r.StatusCode < 300 && r.StatusCode >= 200 {
				// Ack 200s manually.
				log.Info("Sending Ack")
				tx.Ack()
				return
			}
		case e := <-tx.Errors():
			log.Warn(e.Error())
			return
		}
	}
}
