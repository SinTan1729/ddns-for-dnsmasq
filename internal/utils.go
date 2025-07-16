package internal

import (
	"errors"
	"net"
	"net/http"
)

func getClientInfo(req *http.Request, h string) (string, string, error) {
	var errIP error

	ipString := req.Header.Get(h)
	if ipString == "" {
		ipString = req.RemoteAddr
	}

	ip, port, err := net.SplitHostPort(ipString)
	if err != nil {
		errIP = errors.New("Error reading your IP!")
	}
	if net.ParseIP(ip) == nil {
		errIP = errors.New("Request has an invalid IP!")
	}

	return ip, port, errIP
}

func newHTTPError(msg string) httpError {
	return httpError{
		Error:  true,
		Reason: msg,
	}
}
