package coredns_adblock

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// init registers this plugin.
func init() { plugin.Register("adblock", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // Ignore "adblock" and give us the next token.
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error("adblock", c.ArgErr())
	}

	domains, err := Download()
	if err != nil {
		plugin.Error("adblock", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Adblock{Next: next, Domains: domains}
	})

	// All OK, return a nil error.
	return nil
}
