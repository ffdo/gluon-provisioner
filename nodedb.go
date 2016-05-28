package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

type NodeDb struct {
	nodes map[string]*Node
	graph *Graph
	ips   map[string]*Node
}

func NewNodeDb(nodesPath, graphPath string) (ndb *NodeDb, err error) {
	ndb = &NodeDb{}

	var nl NodeList
	err = GetJson(nodesPath, &nl)
	if err != nil {
		ndb = nil
		return
	}
	ndb.nodes = nl.Nodes

	var graph Graph
	err = GetJson(graphPath, &graph)
	if err != nil {
		ndb = nil
		return
	}
	ndb.graph = &graph

	for _, link := range graph.Batadv.Links {
		if !(link.Source < uint(len(graph.Batadv.Nodes)) && link.Target < uint(len(graph.Batadv.Nodes))) {
			err = errors.New("Node index out of range")
			return
		}
		link.SourceMac = graph.Batadv.Nodes[link.Source].Id
		link.TargetMac = graph.Batadv.Nodes[link.Target].Id

		nodeId := graph.Batadv.Nodes[link.Source].NodeId
		if nodeId == nil {
			continue
		}
		node, ok := ndb.nodes[*nodeId]
		if ok {
			node.Links = append(node.Links, link)
			link.SourceNode = node
		}

		nodeId = graph.Batadv.Nodes[link.Target].NodeId
		if nodeId == nil {
			continue
		}
		node, ok = ndb.nodes[*nodeId]
		if ok {
			node.Links = append(node.Links, link)
			link.TargetNode = node
		}
	}

	ndb.ips = make(map[string]*Node)
	for _, n := range ndb.nodes {
		for _, ip := range n.Nodeinfo.Network.Addresses {
			nip := net.ParseIP(ip)
			// Filter link-local and ULA addresses
			if nip != nil && !bytes.Equal(nip[0:2], []byte{0xfe, 0x80}) && !bytes.Equal(nip[0:2], []byte{0xfd, 0xa0}) {
				ndb.ips[nip.String()] = n
			}
		}
	}

	return
}

func (ndb *NodeDb) Dump() {

	// gdb, err := NewGatewayDb("bat0")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for mac, _ := range gdb {
	// 	fmt.Println(mac)
	// }

	for _, node := range ndb.ips {
		tags := []string{}
		if !node.Flags.Online {
			continue
		}

		if len(node.Links) == 0 {
			tags = append(tags, "no_linkinfo")
		} else {
			only_vpn := true
			// only_gateways := true
			for _, link := range node.Links {
				if !link.Vpn {
					only_vpn = false
				}

				// if !(gdb[link.SourceMac] || gdb[link.TargetMac]) {
				// 	only_gateways = false
				// }
			}

			if only_vpn {
				tags = append(tags, "only_vpn")
			}
		}

		// if only_gateways {
		// 	tags = append(tags, "only_gateways")
		// }

		fmt.Print(node.Nodeinfo.Hostname)
		for _, tag := range tags {
			fmt.Print(" ", tag)
		}
		fmt.Println(" ", len(node.Links))
	}

}

func GetJson(path string, result interface{}) (err error) {
	lowerPath := strings.ToLower(path)
	if strings.HasPrefix(lowerPath, "http://") || strings.HasPrefix(lowerPath, "https://") {

		var resp *http.Response
		resp, err = http.Get(path)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			err = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			return
		}

		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(result)

	} else {
		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		err = dec.Decode(result)
	}

	return
}
