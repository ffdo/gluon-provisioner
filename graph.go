package main

type Graph struct {
	Batadv struct {
		Multigraph bool `json:"multigraph"`
		Nodes      []struct {
			ID     string `json:"id"`
			NodeID string `json:"node_id"`
		} `json:"nodes"`
		Directed bool    `json:"directed"`
		Links    []*Link `json:"links"`
	} `json:"batadv"`
	Version int `json:"version"`
}

type Link struct {
	Bidirect bool    `json:"bidirect"`
	Source   int     `json:"source"`
	Target   int     `json:"target"`
	TQ       float64 `json:"tq"`
	VPN      bool    `json:"vpn"`
}
