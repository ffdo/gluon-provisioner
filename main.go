package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

const (
	logFlags = log.Lshortfile
)

func init() {
	log.SetFlags(logFlags)
}

func main() {
	var (
		listenAddress  = flag.String("listen", "[::1]:6060", "HTTP listen address")
		configFile     = flag.String("config", "provisioner.toml", "Path to configuration file")
		nodesPath      = flag.String("nodes", "https://map.ffdo.de/data/nodes.json", "URL (or local path) to nodes.json")
		graphPath      = flag.String("graph", "https://map.ffdo.de/data/graph.json", "URL (or local path) to graph.json")
		defaultDomain  = flag.String("domain", "default", "Default domain name")
		xffHeader      = flag.String("xff", "X-Forwarded-For", "Name of header supplying remote IP")
		redirectHeader = flag.String("redirect", "X-Accel-Redirect", "Name of internal redirection header to set")
	)

	flag.Parse()

	nodeCache := NewNodeCache(60, *nodesPath, *graphPath)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		setPrefix := func(prefix string) {
			w.Header().Add(*redirectHeader, fmt.Sprintf("/%s%s", prefix, req.RequestURI))
		}

		remoteIp := net.ParseIP(req.Header.Get(*xffHeader))
		if remoteIp == nil {
			log.Println("Cannot parse IP address in", xffHeader, "header")
			setPrefix(*defaultDomain)
			return
		}

		ndb := nodeCache.Get()
		if ndb == nil {
			log.Println(remoteIp.String(), "Node/Graph DB is empty")
			setPrefix(*defaultDomain)
			return
		}

		// Look up node in alfred data
		node, ok := ndb.ips[remoteIp.String()]
		if !ok {
			log.Println(remoteIp.String(), "IP not found in alfred data")
			setPrefix(*defaultDomain)
			return
		}

		// Load configuration
		config, err := NewConfig(*configFile)
		if err != nil {
			log.Println("Error loading config file", err)
			setPrefix(*defaultDomain)
			return
		}

		// Check if node should be moved
		domain, force, ignore, err := config.getDomain(node)
		if err != nil {
			log.Println(remoteIp.String(), node.Nodeinfo.Hostname, "Error looking up target domain:", err)
			setPrefix(*defaultDomain)
			return
		}

		if domain == "" {
			log.Println(remoteIp.String(), node.Nodeinfo.Hostname, "Node should not be moved.")
			setPrefix(*defaultDomain)
			return
		} else {
			log.Printf("%s %s Node should be moved to %s (force=%v, ignore=%v)",
				remoteIp.String(), node.Nodeinfo.Hostname, domain, force, ignore)
		}

		if !force {
			// Check if mesh links allow node to be moved safely
			move, err := node.CanBeMoved()
			if err != nil || !move {
				log.Println(remoteIp.String(), node.Nodeinfo.Hostname, "Node cannot be moved:", err)
				setPrefix(*defaultDomain)
				return
			}
		}

		if ignore {
			setPrefix(*defaultDomain)
		} else {
			setPrefix(domain)
		}

	})

	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("Listening on", *listenAddress)
}
