package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SinTan1729/ddns-for-dnsmasq/internal"
	yaml "github.com/goccy/go-yaml"
)

func main() {
	log.SetFlags(0)
	config := internal.Config{IPHeader: "", Port: 4187}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		configFile, err := os.Open(configPath)
		if err == nil {
			err = yaml.NewDecoder(configFile).Decode(&config)
			if err != nil {
				log.Fatalln("Config file is malformed. Exiting.")
			}
		} else {
			log.Println("Not config found at provided path. Using default values.")
		}
		configFile.Close()
	}
	var hostfile internal.Hostfile
	hostfile.Init(os.Getenv("HOSTFILE_PATH"))

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

	if config.IPHeader != "" {
		log.Printf("Using header %v for reading IP.\n", config.IPHeader)
	}
	portString := fmt.Sprintf(":%v", config.Port)
	log.Println("Server listening at", portString)
	http.ListenAndServe(portString, nil)
}
