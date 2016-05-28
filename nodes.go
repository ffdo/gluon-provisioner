package main

import (
	"errors"
	"log"
)

type TrafficStats struct {
	Bytes   uint `json:"bytes,omitempty"`
	Packets uint `json:"packets,omitempty"`
	Dropped uint `json:"dropped,omitempty"`
}

type NodeList struct {
	Timestamp string           `json:"timestamp"`
	Version   uint             `json:"version"`
	Nodes     map[string]*Node `json:"nodes"`
}

type Node struct {
	Firstseen string `json:"firstseen"`
	Lastseen  string `json:"lastseen"`
	Flags     struct {
		Gateway bool `json:"gateway"`
		Online  bool `json:"online"`
	} `json:"flags"`
	Nodeinfo struct {
		Hardware struct {
			Model string `json:"model"`
			Nproc uint   `json:"nproc,omitempty"`
		} `json:"hardware"`
		Hostname string `json:"hostname"`
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		Network struct {
			Addresses      []string `json:"addresses"`
			Mac            string   `json:"mac"`
			MeshInterfaces []string `json:"mesh_interfaces,omitempty"`
			Mesh           map[string]struct {
				Interfaces struct {
					Tunnel   []string `json:"tunnel,omitempty"`
					Wireless []string `json:"wireless,omitempty"`
					Other    []string `json:"other,omitempty"`
				} `json:"interfaces,omitempty"`
			} `json:"mesh,omitempty"`
		} `json:"network"`
		NodeId string `json:"node_id"`
		Owner  struct {
			Contact string `json:"contact"`
		} `json:"owner"`
		System struct {
			SiteCode string `json:"site_code,omitempty"`
		} `json:"system,omitempty"`
		Software struct {
			Autoupdater struct {
				Branch  string `json:"branch"`
				Enabled bool   `json:"enabled"`
			} `json:"autoupdater"`
			BatmanAdv struct {
				Compat  uint   `json:"compat"`
				Version string `json:"version"`
			} `json:"batman-adv"`
			Fastd struct {
				Version string `json:"version"`
				Enabled bool   `json:"enabled"`
			} `json:"fastd"`
			Firmware struct {
				Base    string `json:"base"`
				Release string `json:"release"`
			} `json:"firmware"`
		} `json:"software"`
	} `json:"nodeinfo"`
	Statistics struct {
		Clients     uint    `json:"clients"`
		Gateway     string  `json:"gateway"`
		Loadavg     float64 `json:"loadavg"`
		MemoryUsage float64 `json:"memory_usage"`
		RootfsUsage float64 `json:"rootfs_usage"`
		Uptime      float64 `json:"uptime"`
		Traffic     struct {
			Tx      TrafficStats `json:"tx,omitempty"`
			Rx      TrafficStats `json:"rx,omitempty"`
			MgmtTx  TrafficStats `json:"mgmt_tx,omitempty"`
			MgmtRx  TrafficStats `json:"mgmt_rx,omitempty"`
			Forward TrafficStats `json:"forward,omitempty"`
		} `json:"traffic"`
	} `json:"statistics"`
	Links []*Link `json:"-"`
}

func (node *Node) CanBeMoved() (move bool, err error) {
	if !node.Flags.Online {
		err = errors.New("Node is offline")
		return
	}

	if len(node.Links) == 0 {
		err = errors.New("No link data found")
		return
	}

	gdb, err := NewGatewayDb("bat0")
	if err != nil {
		return
	}

	onlyGateways := true
	for _, link := range node.Links {

		sourceName := "unknown"
		if link.SourceNode != nil {
			sourceName = link.SourceNode.Nodeinfo.Hostname
		}

		if link.SourceNode != node {
			if _, ok := gdb[link.SourceMac]; ok {
				sourceName = "GW:" + sourceName
			} else {
				// onlyGateways = false
			}
		}

		targetName := "unknown"
		if link.TargetNode != nil {
			targetName = link.TargetNode.Nodeinfo.Hostname
		}

		if link.TargetNode != node {
			if _, ok := gdb[link.TargetMac]; ok {
				targetName = "GW:" + targetName
			} else {
				// onlyGateways = false
			}
		}

		// Trust VPN flag from ffmap-backend
		if !link.Vpn {
			onlyGateways = false
		}

		log.Printf("    Links for %s (%s): %s (%s) -> %s (%s) (VPN flag: %v)", node.Nodeinfo.Network.Mac, node.Nodeinfo.Hostname, link.SourceMac, sourceName, link.TargetMac, targetName, link.Vpn)
	}

	// ToDo: Upgrade routers that have only one mesh link ?

	if !onlyGateways {
		err = errors.New("Node is meshing over wifi")
		return
	}

	move = true
	return
}
