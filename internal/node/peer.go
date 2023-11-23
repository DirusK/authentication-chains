/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"crypto/rsa"
	"sync"

	"authentication-chains/internal/types"
)

type (
	// Peer is a node known to the current node.
	Peer struct {
		name          string
		deviceID      rsa.PublicKey
		clusterHeadID rsa.PublicKey
		client        types.NodeClient
	}

	// Peers is a list of known nodes.
	Peers struct {
		mutex sync.RWMutex
		peers []Peer
	}
)

// NewPeers creates a new peers instance.
func NewPeers(peers ...Peer) *Peers {
	return &Peers{
		mutex: sync.RWMutex{},
		peers: peers,
	}
}

// IsEmpty checks if the peer is empty.
func (p *Peers) IsEmpty() bool {
	return len(p.peers) == 0
}

// GetPeers returns a list of peers.
func (p *Peers) GetPeers() []Peer {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	peers := make([]Peer, len(p.peers))
	copy(peers, p.peers)

	return peers
}

// SetPeers sets a list of peers.
func (p *Peers) SetPeers(peers []Peer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.peers = peers
}

// AddPeer adds a peer to the list.
func (p *Peers) AddPeer(peer Peer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.peers = append(p.peers, peer)
}
