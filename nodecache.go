package main

import (
	"log"
	"time"
)

type NodeCache struct {
	dbChan chan *NodeDb
}

func NewNodeCache(updateInterval int, nodesPath, graphPath string) (nc *NodeCache) {
	nc = &NodeCache{
		dbChan: make(chan *NodeDb),
	}

	updateChan := make(chan *NodeDb)
	go func() {
		for {
			ndb, err := NewNodeDb(nodesPath, graphPath)
			if err != nil {
				log.Println("Error updating node cache:", err)
			} else {
				log.Println("Node cache update successful")
				updateChan <- ndb
			}
			time.Sleep(time.Duration(updateInterval) * time.Second)
		}
	}()

	go func() {
		var ndb *NodeDb

		for {
			select {
			case ndb = <-updateChan:
			case nc.dbChan <- ndb:
			}
		}
	}()

	return
}

func (nc *NodeCache) Get() (ndb *NodeDb) {
	ndb = <-nc.dbChan
	return
}
