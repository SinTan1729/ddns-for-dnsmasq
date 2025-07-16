package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type ipInfo struct {
	IP   string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
}

type HTTPError struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason,omitempty"`
}

func whoami(w http.ResponseWriter, req *http.Request) {
	var status int
	var body []byte
	defer func() {
		w.WriteHeader(status)
		fmt.Fprintln(w, string(body))
	}()

	app := req.Context().Value("app").(*appData)
	h := app.IPHeader
	ipString := req.Header.Get(h)
	if ipString == "" {
		ipString = req.RemoteAddr
	}

	ip, port, err := net.SplitHostPort(ipString)
	if err != nil {
		status = http.StatusNotFound
		body, _ = json.Marshal(HTTPError{Error: true, Reason: "Error getting your IP!"})
		return
	}
	if net.ParseIP(ip) == nil {
		status = http.StatusInternalServerError
		body, _ = json.Marshal(HTTPError{Error: true, Reason: "Got an invalid IP!"})
		return
	}

	status = http.StatusOK
	body, _ = json.Marshal(ipInfo{IP: ip, Port: port})
}
