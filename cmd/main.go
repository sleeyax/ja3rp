package main

import (
	"bufio"
	"flag"
	"github.com/sleeyax/ja3rp"
	"log"
	"net/url"
	"os"
)

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result, scanner.Err()

}

func main() {
	address := flag.String("a", "localhost:8080", "Address to listen on")
	cert := flag.String("c", "", "Path to SSL certificate")
	key := flag.String("k", "", "Path to SSL certificate key")
	destination := flag.String("d", "", "Optional URL to forward traffic to (enables reverse proxy mode if set)")
	blacklistFile := flag.String("b", "", "Path to file containing blacklisted JA3 hashes (separated by \n)")
	whitelistFile := flag.String("w", "", "Path to file containing whitelisted JA3 hashes (separated by \n)")

	flag.Parse()

	if *cert == "" || *key == "" {
		log.Fatal("SSL certificate and key required!")
	}

	var u *url.URL
	if *destination != "" {
		u, _ = url.Parse(*destination)
	}

	o := &ja3rp.ServerOptions{
		Destination: u,
	}

	if *whitelistFile != "" {
		whitelist, err := readFile(*whitelistFile)
		if err != nil {
			log.Fatal(err)
		}
		o.Whitelist = whitelist
	}

	if *blacklistFile != "" {
		blacklist, err := readFile(*blacklistFile)
		if err != nil {
			log.Fatal(err)
		}
		o.Blacklist = blacklist
	}

	s := ja3rp.NewServer(*address, *o)

	log.Println("Started listening on https://" + *address)

	log.Fatal(s.ListenAndServeTLS(*cert, *key))
}
