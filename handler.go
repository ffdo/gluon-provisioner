package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func PathHandler(path string, pathConfig *PathConfig, nodeDB *NodeDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		remoteIP := net.ParseIP(req.Header.Get(*xffHeader))
		if remoteIP == nil {
			msg := fmt.Sprint("Cannot parse IP address in ", *xffHeader, " header")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		var node *Node
	RulesLoop:
		for _, rule := range pathConfig.Rules {
			if len(rule.When) != 0 || rule.Careful {
				if node == nil {
					node = nodeDB.GetNode(remoteIP)
				}
				if node == nil {
					log.Println("Node not found in node DB, cannot evaluate rule.")
					continue
				}
			}

			for _, condition := range rule.When {
				if !condition.Check(node.Nodeinfo) {
					continue RulesLoop
				}
			}

			if rule.Careful && !node.HasOnlyVPNLinks() {
				log.Println("Careful enabled and node is not known as VPN only, skipping rule.")
				continue
			}

			if rule.Disabled {
				log.Println("Matching rule is disabled")
				continue
			}
			redirect(w, req.URL.Path, path, rule.Path)
			return
		}

		redirect(w, req.URL.Path, path, pathConfig.Default)
		return
	})
}
