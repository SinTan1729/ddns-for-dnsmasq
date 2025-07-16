package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	yaml "github.com/goccy/go-yaml"
)

type appData struct {
	IPHeader string `yaml:"ip-header"`
	Port     int    `yaml:"port"`
}

func main() {
	log.SetFlags(0)
	app := appData{IPHeader: "", Port: 4187}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		configFile, err := os.Open(configPath)
		if err == nil {
			yaml.NewDecoder(configFile).Decode(&app)
		} else {
			log.Println("Not config found at provided path. Using default values.")
		}
		configFile.Close()
	}

	http.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "app", &app)
		r = r.WithContext(ctx)
		whoami(w, r)
	})

	if app.IPHeader != "" {
		log.Printf("Using header %v for reading IP.\n", app.IPHeader)
	}
	portString := fmt.Sprintf(":%v", app.Port)
	log.Println("Server listening at", portString)
	http.ListenAndServe(portString, nil)
}
