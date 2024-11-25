package coredns_adblock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

var url = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"

// Download the list of domains to block
func Download() ([]string, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	rows := bytes.Split(body, []byte("\n"))
	var domains []string
	for _, row := range rows {
		if bytes.HasPrefix(row, []byte("0.0.0.0")) {
			columns := bytes.Split(row, []byte(" "))
			domains = append(domains, string(columns[1])+".")
		}
	}
	return domains, nil
}
