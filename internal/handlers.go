package internal

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

const version = "0.2.2"

func Version(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "DDNS for Dnsmasq v%v\n", version)
}

func WhoAmI(w http.ResponseWriter, req *http.Request) {
	var status int
	var body []byte
	config := req.Context().Value("config").(*Config)

	ip, err := getClientInfo(req, config.IPHeader)
	if err == nil {
		status = http.StatusOK
		body = []byte(ip)
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
		w.WriteHeader(status)
		fmt.Fprintln(w, string(body))
	}()

	config := req.Context().Value("config").(*Config)
	hostfile := req.Context().Value("hostfile").(*Hostfile)

	var data hostEntry
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		status = http.StatusBadRequest
		body, _ = json.Marshal(newHTTPError("The request was malformed."))
		return
	}

	ok := validAuth(req, config, &data)
	if ok {
		ip := data.IP
		if ip == "" {
			reqIP, err := getClientInfo(req, config.IPHeader)
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
		hostfile.update(data.Host, ip)
		status = http.StatusOK
		body, _ = json.Marshal(hostEntry{Host: data.Host, IP: ip})
	} else {
		status = http.StatusUnauthorized
		body, _ = json.Marshal(newHTTPError("Wrong API key was provided."))
	}

}

func GetInfo(w http.ResponseWriter, req *http.Request) {
	var status int
	var body []byte
	defer func() {
		w.WriteHeader(status)
		fmt.Fprintln(w, string(body))
	}()
	config := req.Context().Value("config").(*Config)
	hostfile := req.Context().Value("hostfile").(*Hostfile)

	var data hostEntry
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		status = http.StatusBadRequest
		body, _ = json.Marshal(newHTTPError("The request was malformed."))
		return
	}

	ok := validAuth(req, config, &data)
	if ok {
		host, ok := hostfile.hosts[data.Host]
		if ok {
			status = http.StatusOK
			body, _ = json.Marshal(host)
		} else {
			status = http.StatusNotFound
			body, _ = json.Marshal(newHTTPError("There's no entry for this host."))
		}
	} else {
		status = http.StatusUnauthorized
		body, _ = json.Marshal(newHTTPError("Wrong API key was provided."))
	}
}
