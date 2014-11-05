package main

import (
	"github.com/stefankopieczek/gossip/base"
	"github.com/stefankopieczek/gossip/log"
	"github.com/stefankopieczek/gossip/transaction"
)

func main() {
	log.SetDefaultLogLevel(log.DEBUG)
	// Build an INVITE to send.
	stefan := "stefan"
	ryan := "ryan"
	callid := base.CallId("thisisacall3")
	branch := "z9hG4bKcallbranch1"
	tcp := "tcp"
	var port uint16 = 5060
	invite := base.NewRequest(
		base.INVITE,
		&base.SipUri{
			User: &stefan,
			Host: "172.18.115.72",
		},
		"SIP/2.0",
		[]base.SipHeader{
			&base.ViaHeader{
				&base.ViaHop{
					ProtocolName:    "SIP",
					ProtocolVersion: "2.0",
					Transport:       "TCP",
					Host:            "172.18.115.75",
					Port:            &port,
					Params: base.Params{
						"branch": &branch,
					},
				},
			},
			&base.ToHeader{
				Address: &base.SipUri{
					User: &stefan,
					Host: "172.18.115.72",
				},
			},
			&base.FromHeader{
				Address: &base.SipUri{
					User: &ryan,
					Host: "172.18.115.75",
					UriParams: base.Params{
						"transport": &tcp,
					},
				},
			},
			&base.ContactHeader{
				Address: &base.SipUri{
					User: &ryan,
					Host: "172.18.115.75",
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

	tm, err := transaction.NewManager("tcp", "172.18.115.75:5060")
	if err != nil {
		panic(err)
	}

	tx := tm.Send(invite, "172.18.115.72:5060")

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
