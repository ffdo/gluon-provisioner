package main

import (
	"bytes"
	"log"
	"net"
	"sync"
	"time"
)

type NodeDB struct {
	mutex              sync.Mutex
	nodesURL, graphURL string
	nodes              map[string]*Node
}

func NewNodeDB(updateInterval time.Duration, nodesURL, graphURL string) *NodeDB {
	ndb := &NodeDB{
		nodesURL: nodesURL,
		graphURL: graphURL,
	}

	go func() {
		for {
			ndb.update()
			time.Sleep(updateInterval)
		}
	}()

	return ndb
}

func (ndb *NodeDB) update() {
	var nodes Nodes
	err := GetJSON(ndb.nodesURL, &nodes)
	if err != nil {
		log.Println("Error fetching nodes:", err)
		return
	}

	var graph Graph
	err = GetJSON(ndb.graphURL, &graph)
	if err != nil {
		log.Println("Error fetching graph:", err)
		return
	}

	for _, link := range graph.Batadv.Links {
		if !(link.Source < len(graph.Batadv.Nodes) && link.Target < len(graph.Batadv.Nodes)) {
			log.Println("Node index out of range")
			return
		}

		nodeID := graph.Batadv.Nodes[link.Source].NodeID
		if nodeID == "" {
			continue
		}
		node, ok := nodes.Nodes[nodeID]
		if ok {
			node.Links = append(node.Links, link)
		}

		nodeID = graph.Batadv.Nodes[link.Target].NodeID
		if nodeID == "" {
			continue
		}
		node, ok = nodes.Nodes[nodeID]
		if ok {
			node.Links = append(node.Links, link)
		}
	}

	ips := make(map[string]*Node)
	for _, n := range nodes.Nodes {
		for _, ip := range n.Nodeinfo.Network.Addresses {
			nip := net.ParseIP(ip)
			// Filter link-local and ULA addresses
			if nip != nil && !bytes.Equal(nip[0:2], []byte{0xfe, 0x80}) && !bytes.Equal(nip[0:2], []byte{0xfd, 0xa0}) {
				ips[nip.String()] = n
			}
		}
	}

	ndb.mutex.Lock()
	defer ndb.mutex.Unlock()
	ndb.nodes = ips
	return
}

func (ndb *NodeDB) GetNode(ip net.IP) *Node {
	ndb.mutex.Lock()
	defer ndb.mutex.Unlock()
	return ndb.nodes[ip.String()]
}
