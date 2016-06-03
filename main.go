package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	logFlags = log.Lshortfile
)

var (
	updateInterval = flag.Duration("update", 60*time.Second, "Node cache update interval")
	redirectHeader = flag.String("redirect", "X-Accel-Redirect", "Name of internal redirection header to set")
	xffHeader      = flag.String("xff", "X-Forwarded-For", "Name of header supplying remote IP")
)

func init() {
	log.SetFlags(logFlags)
}

func main() {
	var (
		configFile    = flag.String("config", "gluon-provisioner.yaml", "Path to configuration file")
		listenAddress = flag.String("listen", "[::1]:6060", "HTTP listen address")
	)
	flag.Parse()

	config, err := NewConfig(*configFile)
	if err != nil {
		log.Fatalln("Error loading config file:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		remoteIP := net.ParseIP(req.Header.Get(*xffHeader))
		if remoteIP == nil {
			msg := fmt.Sprint("Cannot parse IP address in ", *xffHeader, " header")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		serveMux := config.GetServeMux(remoteIP)
		if serveMux != nil {
			serveMux.ServeHTTP(w, req)
		} else {
			http.NotFound(w, req)
		}
		return
	})

	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("Listening on", *listenAddress)
}
