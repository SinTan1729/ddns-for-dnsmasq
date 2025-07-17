package internal

import (
	"errors"
	"net"
	"net/http"
	"strings"
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

func validAuth(req *http.Request, c *Config, data *hostEntry) (HostConfig, bool) {
	entry, ok := c.Hosts[data.Host]
	if ok && entry.APIKey == req.Header.Get("X-API-Key") {
		return entry, true
	}
	return entry, false
}
