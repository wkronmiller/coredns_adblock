package coredns_adblock

import (
  "testing"
)

func TestDownload(t *testing.T) {
  domains, err := Download()
  if err != nil {
    t.Error(err)
  }
  if len(domains) < 1 {
    t.Errorf("Domains list was empty")
  }
}
