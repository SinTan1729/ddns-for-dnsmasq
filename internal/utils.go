package internal

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/argon2id"
	yaml "github.com/goccy/go-yaml"
)

func (c *Config) Init() {
	*c = Config{IPHeader: "", Port: 4187}
	configPath := strings.Trim(os.Getenv("CONFIG_PATH"), `"`)
	if configPath != "" {
		configFile, err := os.Open(configPath)
		if err == nil {
			err = yaml.NewDecoder(configFile).Decode(c)
			if err != nil {
				log.Fatalln("Config file is malformed. Exiting.")
			}
		} else {
			log.Println("Not config found at provided path. Using default values.")
		}
		configFile.Close()
	}
	for name, host := range c.Hosts {
		hash := host.APIKey
		_, _, _, err := argon2id.DecodeHash(hash)
		if err != nil {
			log.Printf("The API key hash for %v seems to be invalid.\n", name)
			log.Fatalf("%v\nPlease fix it. Exiting for now.\n", err)
		}
	}
}

func getClientInfo(req *http.Request, h string) (string, error) {
	var err error

	ipList := req.Header.Get(h)
	if ipList == "" {
		ipList = req.RemoteAddr
	}
	hostport := strings.TrimSpace(strings.SplitN(ipList, ",", 2)[0])
	ip, _, err := net.SplitHostPort(hostport)
	if err != nil {
		// This is needed since reverse proxies don't set port
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

func validAuth(req *http.Request, c *Config, data *hostEntry) bool {
	entry, ok := c.Hosts[data.Host]
	if ok {
		match, errValidate := argon2id.ComparePasswordAndHash(req.Header.Get("X-API-Key"), entry.APIKey)
		switch errValidate {
		case nil:
			return match
		default:
			log.Printf("Got the following error while processing the API key hash for %v:\n", data.Host)
			log.Fatalf("%v\nPlease fix it. Exiting for now.\n", errValidate)
		}
	}
	return false
}
