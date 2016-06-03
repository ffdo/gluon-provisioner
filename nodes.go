package main

import (
	"time"
)

type Nodes struct {
	Timestamp time.Time        `json:"timestamp"`
	Version   int              `json:"version"`
	Nodes     map[string]*Node `json:"nodes"`
}

type Node struct {
	Nodeinfo   *Nodeinfo `json:"nodeinfo"`
	Statistics struct {
		Clients     int     `json:"clients"`
		Gateway     string  `json:"gateway"`
		Loadavg     float64 `json:"loadavg"`
		MemoryUsage float64 `json:"memory_usage"`
		RootfsUsage float64 `json:"rootfs_usage"`
		Uptime      float64 `json:"uptime"`
		Traffic     struct {
			Tx      TrafficStats `json:"tx,omitempty"`
			Rx      TrafficStats `json:"rx,omitempty"`
			Forward TrafficStats `json:"forward,omitempty"`
			MgmtTx  TrafficStats `json:"mgmt_tx,omitempty"`
			MgmtRx  TrafficStats `json:"mgmt_rx,omitempty"`
		} `json:"traffic"`
	} `json:"statistics"`
	Flags struct {
		Gateway bool `json:"gateway"`
		Online  bool `json:"online"`
		Uplink  bool `json:"uplink"`
	} `json:"flags"`
	Lastseen  string  `json:"lastseen"`
	Firstseen string  `json:"firstseen"`
	Links     []*Link `json:"-"`
}

func (node *Node) HasOnlyVPNLinks() bool {
	if !node.Flags.Online || !node.Flags.Uplink {
		return false
	}

	for _, link := range node.Links {
		if !link.VPN {
			return false
		}
	}

	return true
}

type Nodeinfo struct {
	NodeId  string `json:"node_id"`
	Network struct {
		Mac       string   `json:"mac"`
		Addresses []string `json:"addresses"`
		Mesh      map[string]struct {
			Interfaces struct {
				Wireless []string `json:"wireless,omitempty"`
				Other    []string `json:"other,omitempty"`
				Tunnel   []string `json:"tunnel,omitempty"`
			} `json:"interfaces,omitempty"`
		} `json:"mesh,omitempty"`
		MeshInterfaces []string `json:"mesh_interfaces,omitempty"`
	} `json:"network"`
	Owner struct {
		Contact string `json:"contact,omitempty"`
	} `json:"owner,omitempty"`
	System struct {
		SiteCode string `json:"site_code,omitempty"`
	} `json:"system,omitempty"`
	Hostname string `json:"hostname"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Software struct {
		Autoupdater struct {
			Branch  string `json:"branch"`
			Enabled bool   `json:"enabled"`
		} `json:"autoupdater"`
		BatmanAdv struct {
			Compat  int    `json:"compat"`
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
		StatusPage struct {
			Api int `json:"api"`
		} `json:"status-page,omitempty"`
	} `json:"software"`
	Hardware struct {
		Model string `json:"model"`
		Nproc int    `json:"nproc,omitempty"`
	} `json:"hardware"`
}

type TrafficStats struct {
	Bytes   float64 `json:"bytes,omitempty"`
	Packets int     `json:"packets,omitempty"`
	Dropped int     `json:"dropped,omitempty"`
}
