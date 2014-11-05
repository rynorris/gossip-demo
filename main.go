package main

import (
	"fmt"

	"github.com/stefankopieczek/gossip/base"
	"github.com/stefankopieczek/gossip/log"
	"github.com/stefankopieczek/gossip/transaction"
)

var (
	// Call parameters
	transport = "UDP"
	callid    = base.CallId("thisisacall3")
	branch    = "z9hG4bK.callbranch1"

	// Caller parameters
	caller             = "ryan"
	src_host           = "centosvm-rpn"
	listen_port uint16 = 5060

	// Callee parameters
	callee          = "stefan"
	dst_host        = "172.18.115.72"
	dst_port uint16 = 5060
)

func main() {
	log.SetDefaultLogLevel(log.DEBUG)
	// Build an INVITE to send.
	invite := base.NewRequest(
		base.INVITE,
		&base.SipUri{
			User: &callee,
			Host: dst_host,
		},
		"SIP/2.0",
		[]base.SipHeader{
			&base.ViaHeader{
				&base.ViaHop{
					ProtocolName:    "SIP",
					ProtocolVersion: "2.0",
					Transport:       transport,
					Host:            dst_host,
					Port:            &dst_port,
					Params: base.Params{
						"branch": &branch,
					},
				},
			},
			&base.ToHeader{
				Address: &base.SipUri{
					User: &callee,
					Host: dst_host,
				},
			},
			&base.FromHeader{
				Address: &base.SipUri{
					User: &caller,
					Host: src_host,
					UriParams: base.Params{
						"transport": &transport,
					},
				},
			},
			&base.ContactHeader{
				Address: &base.SipUri{
					User: &caller,
					Host: src_host,
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

	tm, err := transaction.NewManager(transport, fmt.Sprintf("%v:%v", src_host, listen_port))
	if err != nil {
		panic(err)
	}

	tx := tm.Send(invite, fmt.Sprintf("%v:%v", dst_host, dst_port))

	for {
		select {
		case r := <-tx.Responses():
			log.Info(r.String())
			if r.StatusCode < 300 && r.StatusCode >= 200 {
				// Ack 200s manually.
				log.Info("Sending Ack")
				tx.Ack()
			}
		case e := <-tx.Errors():
			log.Warn(e.Error())
		}
	}
}
