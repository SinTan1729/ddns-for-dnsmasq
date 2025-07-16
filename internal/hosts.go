package internal

import (
	"bufio"
	"log"
	"os"

	"github.com/wasilibs/go-re2"
)

type Hostfile struct {
	path  string
	hosts map[string]hostEntry
}

func (h Hostfile) Init(s string) {
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
	pattern := re2.MustCompile(`^\s*(\S+)\s+(\S+)\s*$`)
	for scanner.Scan() {
		matches := pattern.FindStringSubmatch(scanner.Text())
		if len(matches) > 0 {
			entry := hostEntry{Host: matches[2], IP: matches[1]}
			hosts[matches[2]] = entry
		} else {
			log.Fatalln("Malformed hostfile.")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	h.path = s
	h.hosts = hosts
}
