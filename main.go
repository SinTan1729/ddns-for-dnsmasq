package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type appData struct {
	ipHeader string
}

func whoami(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	h := ctx.Value("ipHeader").(string)
	var ipString string

	if h == "" {
		ipString = req.RemoteAddr
	} else {
		ipString = req.Header.Get(h)
	}
	ip, port, err := net.SplitHostPort(ipString)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Error getting your IP!")
		return
	}
	if net.ParseIP(ip) == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Got an invalid IP!")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Here's what I can see.\nIP: %v, port: %v", ip, port)
}

func main() {
	app := &appData{ipHeader: ""}
	http.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "ipHeader", app.ipHeader)
		r = r.WithContext(ctx)
		whoami(w, r)
	})

	http.ListenAndServe(":4187", nil)
}
