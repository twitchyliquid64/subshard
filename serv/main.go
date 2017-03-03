package main

import (
	"log"
	"net/http"
	"os"
)

const serverHost = "subshard"

func main() {
	configuration, err := readConfig("subshard_serv.json")
	if err != nil {
		log.Printf("Error loading configuration (%s): %s\n", "subshard_serv.json", err.Error())
		os.Exit(1)
	}

	proxy, err := makeProxyServer(configuration)
	if err != nil {
		log.Printf("Error initializing server: %s\n", err.Error())
		os.Exit(2)
	}

	log.Fatal(http.ListenAndServe(configuration.Listener, proxy))
}
