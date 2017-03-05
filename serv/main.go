package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

const serverHost = "subshard"
const testingConfigPath = "subshard_serv.json"

var configPathVar string

func processFlags() {
	flag.StringVar(&configPathVar, "config", "/etc/subshard/subshard-serv.json", "Path to configuration file")
	flag.Parse()

	// If the main config doesnt exist, check the working directory too for subshard_serv.json.
	if _, err := os.Stat(configPathVar); err != nil && os.IsNotExist(err) {
		if _, err := os.Stat(testingConfigPath); err == nil {
			log.Printf("WARNING: config at %s does not exist, falling back to %s which does exist\n", configPathVar, testingConfigPath)
			configPathVar = "subshard_serv.json"
		}
	}
}

func main() {
	processFlags()
	var wg sync.WaitGroup

	for { //Loop restarts on SIGHUP
		gTLSConfig = nil
		gConfiguration = nil
		gValidBlacklistEntries = nil
		configuration, err := readConfig(configPathVar)
		if err != nil {
			log.Printf("Error loading configuration (%s): %s\n", configPathVar, err.Error())
			os.Exit(1)
		}

		// create proxy server object, set up routing rules from config
		proxy, err := makeProxyServer(configuration)
		if err != nil {
			log.Printf("Error initializing server configuration: %s\n", err.Error())
			os.Exit(2)
		}

		listener, err := setupListener(configuration)
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		go func() { //Listener go-routine terminates when listener is closed
			wg.Add(1)
			defer wg.Done()
			for {
				err := http.Serve(listener, proxy)
				if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
					log.Println("Serve Err: ", err)
				}
				return
			}
		}()

		sig := waitInterrupt()
		listenerCloseErr := listener.Close()
		wg.Wait() //wait for server/listener goroutine to terminate
		if sig == syscall.SIGHUP {
			log.Println("Got SIGHUP, reloading")
			gConfigReloads++
			continue
		}

		if sig == syscall.SIGINT {
			if listenerCloseErr != nil {
				log.Println("Listener close err: ", listenerCloseErr)
			}
		}
		os.Exit(0)
	}
}

func waitInterrupt() os.Signal {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return <-sig
}
