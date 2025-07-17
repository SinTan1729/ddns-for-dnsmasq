package internal

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func getClientInfo(req *http.Request, h string) (string, error) {
	var err error

	ipList := req.Header.Get(h)
	if ipList == "" {
		ipList = req.RemoteAddr
	}
	hostport := strings.TrimSpace(strings.SplitN(ipList, ",", 2)[0])
	ip, _, err := net.SplitHostPort(hostport)
	if err != nil {
		// This is needed since reverse proxies don't set port
		ip = hostport
		err = nil
	}

	if net.ParseIP(ip) == nil {
		err = errors.New("Request has an invalid IP!")
	}

	return ip, err
}

func newHTTPError(msg string) httpError {
	return httpError{
		Error:  true,
		Reason: msg,
	}
}

func validAuth(req *http.Request, c *Config, data *hostEntry) bool {
	entry, ok := c.Hosts[data.Host]
	if ok {
		match, errValidate := argon2id.ComparePasswordAndHash(req.Header.Get("X-API-Key"), entry.APIKey)
		switch errValidate {
		case nil:
			return match
		default:
			log.Printf("Got the following error while processing the API key hash for %v:\n", data.Host)
			log.Fatalf("%v\nPlease fix it. Exiting for now.\n", errValidate)
		}
	}
	return false
}
