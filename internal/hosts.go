package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wasilibs/go-re2"
)

type Hostfile struct {
	path  string
	hosts map[string]hostEntry
}

func (h *Hostfile) update(name string, ip string) {
	host, ok := h.hosts[name]
	if ok {
		if host.IP != ip {
			host.IP = ip
		} else {
			// No need to change anything
			return
		}
		h.hosts[name] = host
	} else {
		h.hosts[name] = hostEntry{Host: name, IP: ip}
	}

	file, err := os.Create(h.path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	w.WriteString("# Generated automatically by DDNS for Dnsmasq\n")
	for _, host := range h.hosts {
		line := fmt.Sprintln(host.IP, host.Host)
		_, err := w.WriteString(line)
		if err != nil {
			log.Fatalln(err)
		}
	}
	w.Flush()
	log.Printf("IP for %v was updated to %v", name, ip)
}

func (h *Hostfile) Init(s string) {
	if s == "" {
		log.Fatalln("No hostfile was provided.")
	}
	file, err := os.Open(s)
	if err != nil {
		log.Fatalln("There was an error reading the hostfile.")
	}
	defer file.Close()

	hosts := make(map[string]hostEntry)

	scanner := bufio.NewScanner(file)
	pattern := re2.MustCompile(`^(\S+)\s+(\S+)\s*#?.*$`)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		matches := pattern.FindStringSubmatch(line)
		if len(matches) > 0 {
			entry := hostEntry{Host: matches[2], IP: matches[1]}
			hosts[matches[2]] = entry
		} else {
			log.Fatalln("Unsupported hostfile.")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	h.path = s
	h.hosts = hosts
}
