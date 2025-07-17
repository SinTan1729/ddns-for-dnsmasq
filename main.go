package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SinTan1729/ddns-for-dnsmasq/internal"
)

func main() {
	log.SetFlags(0)
	var hostfile internal.Hostfile
	hostfile.Init(strings.Trim(os.Getenv("HOSTFILE_PATH"), `"`))
	var config internal.Config
	config.Init()

	http.HandleFunc("GET /version", internal.Version)
	http.HandleFunc("GET /whoami", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "config", &config)
		r = r.WithContext(ctx)
		internal.WhoAmI(w, r)
	})
	http.HandleFunc("PUT /update", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "config", &config)
		ctx = context.WithValue(ctx, "hostfile", &hostfile)
		r = r.WithContext(ctx)
		internal.Update(w, r)
	})
	http.HandleFunc("POST /getinfo", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "config", &config)
		ctx = context.WithValue(ctx, "hostfile", &hostfile)
		r = r.WithContext(ctx)
		internal.GetInfo(w, r)
	})

	if config.IPHeader != "" {
		log.Printf("Using header \"%v\" for reading IP.\n", config.IPHeader)
	}
	portString := fmt.Sprintf(":%v", config.Port)
	log.Println("Server listening at", portString)
	http.ListenAndServe(portString, nil)
}
