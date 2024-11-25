// Package adblock is a CoreDNS plugin that prints "adblock" to stdout on every packet received.
//
// It serves as an adblock CoreDNS plugin with numerous code comments.
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

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("adblock")

// Adblock is an adblock plugin to show how to write a plugin.
type Adblock struct {
	Next    plugin.Handler
	Domains []string
}

// ServeDNS implements the plugin.Handler interface. This method gets called when adblock is used
// in a Server.
func (e Adblock) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	// Export metric with the server label set to the current server handling the request.
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	for _, domain := range e.Domains {
		if strings.HasSuffix(qname, domain) {
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

// ResponsePrinter wrap a dns.ResponseWriter and will write adblock to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "adblock" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	return r.ResponseWriter.WriteMsg(res)
}
