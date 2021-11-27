# JA3RP (JA3 Reverse Proxy)
Ja3RP is a basic reverse proxy server that filters traffic based on [JA3](https://github.com/salesforce/ja3) fingerprints.
It can also operate as a regular HTTP server for testing purposes.

Inspired by this [ja3-server](https://github.com/CapacitorSet/ja3-server) POC.

## Installation
```
# Install library locally:
$ go get github.com/sleeyax/ja3rp

# Install binary globally:
$ go install github.com/sleeyax/ja3rp
```

## Usage
### Preparation
A JA3 hash is constructed from a TLS ClientHello packet.
For this reason the JA3RP server will need an SSL certificate in order to work.

You can generate a self-signed certificate using the following commands:
```
$ openssl req -new -subj "/C=US/ST=Utah/CN=localhost" -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr
$ openssl x509 -req -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt
```

### Package
The following example starts an HTTPS server and filters incoming traffic based on a JA3 hash.
If the hash is found in the whitelist the traffic is forwarded to the configured destination server.
Otherwise or if blacklisted the request is blocked.

```go
package main

import (
	"fmt"
	"github.com/sleeyax/ja3rp"
	"github.com/sleeyax/ja3rp/net/http"
	"log"
	"net/url"
)

func main() {
	address := "localhost:1337"
	d, _ := url.Parse("https://example.com")

	server := ja3rp.NewServer(address, ja3rp.ServerOptions{
		Destination: d,
		Whitelist: []string{
			"bd50e49d418ed1777b9a410d614440c4", // firefox
			"b32309a26951912be7dba376398abc3b", // chrome
		},
		Blacklist: []string{
			"3b5074b1b5d032e5620f69f9f700ff0e", // CURL
		},
		OnBlocked: func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Sorry, you are not in our whitelist :(")
		},
	})

	err := server.ListenAndServeTLS("certificate.crt", "certificate.key")
	
	log.Fatal(err)
}
```

### CLI
```
$ ja3rp -h
Usage: ja3rp -a <address> [-d <destination URL> -c <cert file> -k <cert key> -w <whitelist file> -b <blacklist file>]
Example: $ ja3rp -a localhost:1337 -d https://example.com -c certificate.crt -k certificate.key -w whitelist.txt -b blacklist.txt
```
Hashes should be stored in .txt files, each separated by a new line.

## Licenses
This project is licensed with the [MIT License](LICENSE).

The included (and then modified) `net/http`, `internal/profile` and `crypto` packages fall under the [go source code license](./LICENSE_GO.txt).
