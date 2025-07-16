package internal

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

func WhoAmI(w http.ResponseWriter, req *http.Request) {
	var status int
	var body []byte
	config := req.Context().Value("config").(*Config)

	ip, port, err := getClientInfo(req, config.IPHeader)
	if err == nil {
		status = http.StatusOK
		body, _ = json.Marshal(ipInfo{IP: ip, Port: port})
	} else {
		status = http.StatusInternalServerError
		body, _ = json.Marshal(newHTTPError(err.Error()))
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, string(body))
}

func Update(w http.ResponseWriter, req *http.Request) {
	var status int
	var body []byte
	defer func() {
		status = 200
		w.WriteHeader(status)
		fmt.Fprintln(w, string(body))
	}()

	config := req.Context().Value("config").(*Config)
	hostfile := req.Context().Value("hostfile").(*Hostfile)
	apiKey := req.Header.Get("X-API-Key")

	var data hostEntry
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		status = http.StatusBadRequest
		body, _ = json.Marshal(newHTTPError("The request was malformed."))
		return
	}

	entry, ok := config.Hosts[data.Host]
	if ok {
		if entry.APIKey == apiKey {
			ip := data.IP
			if ip == "" {
				reqIP, _, err := getClientInfo(req, config.IPHeader)
				if err != nil {
					status = http.StatusInternalServerError
					body, _ = json.Marshal(newHTTPError(err.Error()))
					return
				}
				ip = reqIP
			}
			if net.ParseIP(ip) == nil {
				status = http.StatusBadRequest
				body, _ = json.Marshal(newHTTPError("Invalid IP was provided."))
				return
			}
			hostfile.update(entry.Host, ip)
			status = http.StatusOK
			body, _ = json.Marshal(hostEntry{Host: entry.Host, IP: ip})
		} else {
			status = http.StatusUnauthorized
			body, _ = json.Marshal(newHTTPError("Wrong API key was provided."))
		}
		return
	}

	status = http.StatusNotFound
	body, _ = json.Marshal(newHTTPError("Payload didn't match server config."))
}
