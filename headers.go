package main

import "github.com/stefankopieczek/gossip/base"

// Utility methods for creating headers.

func Via(e *endpoint, branch string) *base.ViaHeader {
	return &base.ViaHeader{
		&base.ViaHop{
			ProtocolName:    "SIP",
			ProtocolVersion: "2.0",
			Transport:       e.transport,
			Host:            e.host,
			Port:            &e.port,
			Params: base.Params{
				"branch": &branch,
			},
		},
	}
}

func To(e *endpoint, tag string) *base.ToHeader {
	header := &base.ToHeader{
		DisplayName: &e.displayName,
		Address: &base.SipUri{
			User:      &e.username,
			Host:      e.host,
			UriParams: base.Params{},
		},
		Params: base.Params{},
	}

	if tag != "" {
		header.Params["tag"] = &tag
	}

	return header
}

func From(e *endpoint, tag string) *base.FromHeader {
	header := &base.FromHeader{
		DisplayName: &e.displayName,
		Address: &base.SipUri{
			User: &e.username,
			Host: e.host,
			UriParams: base.Params{
				"transport": &e.transport,
			},
		},
		Params: base.Params{},
	}

	if tag != "" {
		header.Params["tag"] = &tag
	}

	return header
}

func Contact(e *endpoint) *base.ContactHeader {
	return &base.ContactHeader{
		DisplayName: &e.displayName,
		Address: &base.SipUri{
			User: &e.username,
			Host: e.host,
		},
	}
}

func CSeq(seqno uint32, method base.Method) *base.CSeq {
	return &base.CSeq{
		SeqNo:      seqno,
		MethodName: method,
	}
}

func CallId(callid string) *base.CallId {
	header := base.CallId(callid)
	return &header
}

func ContentLength(l uint32) base.ContentLength {
	return base.ContentLength(l)
}
