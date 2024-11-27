// Package adblock is a CoreDNS plugin that blocks ads using the StevenBlack
// blacklist.
// It is based on the CoreDNS example plugin.
package coredns_adblock

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"strings"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("adblock")

type Adblock struct {
	Next    plugin.Handler
	Domains []string
}

func (e Adblock) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	// Export metric with the server label set to the current server handling the request.
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	for _, domain := range e.Domains {
		if strings.HasSuffix(qname, domain) {
			blockedRequestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()
			log.Debugf("Blocking domain: %s\nRequest: %s", domain, qname)

			m := new(dns.Msg)
			m.SetRcode(r, dns.RcodeNameError) // NXDOMAIN response code

			w.WriteMsg(m)
			return dns.RcodeNameError, nil
		}

	}
	log.Debugf("Allowing Request: %s", qname)
	pw := NewResponsePrinter(w)

	// Call next plugin (if any).
	return plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
}

// Name implements the Handler interface.
func (e Adblock) Name() string { return "adblock" }

type ResponsePrinter struct {
	dns.ResponseWriter
}

func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	return r.ResponseWriter.WriteMsg(res)
}
