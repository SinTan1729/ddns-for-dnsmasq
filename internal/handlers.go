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
	appData := req.Context().Value("appData").(*AppData)

	ip, port, err := getClientInfo(req, appData.IPHeader)
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

	appData := req.Context().Value("appData").(*AppData)
	apiKey := req.Header.Get("X-API-Key")

	var data updatePayload
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		status = http.StatusBadRequest
		body, _ = json.Marshal(newHTTPError("The request was malformed."))
		return
	}

	for _, entry := range appData.Hosts {
		if entry.Host == data.Host {
			if entry.APIKey == apiKey {
				ip := data.IP
				if ip == "" {
					reqIP, _, err := getClientInfo(req, appData.IPHeader)
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
				status = http.StatusOK
				body, _ = json.Marshal(updatePayload{Host: entry.Host, IP: ip})
			} else {
				status = http.StatusUnauthorized
				body, _ = json.Marshal(newHTTPError("Wrong API key was provided."))
			}
			return
		}
	}

	status = http.StatusNotFound
	body, _ = json.Marshal(newHTTPError("Payload didn't match server config."))
}
