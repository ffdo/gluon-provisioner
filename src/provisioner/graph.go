package main

type Link struct {
	Bidirect   bool    `json:"bidirect"`
	Source     uint    `json:"source"`
	Target     uint    `json:"target"`
	Tq         float64 `json:"tq"`
	Vpn        bool    `json:"vpn"`
	SourceMac  string  `json:"-"`
	TargetMac  string  `json:"-"`
	SourceNode *Node   `json:"-"`
	TargetNode *Node   `json:"-"`
}

type Graph struct {
	Version    uint `json:"version"`
	Multigraph bool `json:"multigraph"`
	Batadv     struct {
		Directed bool    `json:"directed"`
		Links    []*Link `json:"links"`
		Nodes    []struct {
			Id     string  `json:"id"`
			NodeId *string `json:"node_id"`
		} `json:"nodes"`
	} `json:"batadv"`
}
