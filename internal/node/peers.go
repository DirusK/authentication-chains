/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"sync"

	"authentication-chains/internal/types"
)

type (
	// Peers is a list of known nodes.
	Peers struct {
		mutex sync.RWMutex
		Peers []*Peer
	}

	// Peer is a node known to the current node.
	Peer struct {
		Name          string
		GRPCAddress   string
		DeviceID      []byte
		ClusterHeadID []byte
		Level         uint32
		Client        types.NodeClient
	}
)

// NewPeer creates a new peer instance.
func NewPeer(name string, deviceID, clusterHeadID []byte, GRPCAddress string, level uint32, client types.NodeClient) *Peer {
	return &Peer{
		Name:          name,
		DeviceID:      deviceID,
		ClusterHeadID: clusterHeadID,
		Client:        client,
		GRPCAddress:   GRPCAddress,
		Level:         level,
	}
}

func (p Peer) ToProto() *types.Peer {
	return &types.Peer{
		Name:          p.Name,
		DeviceId:      p.DeviceID,
		ClusterHeadId: p.ClusterHeadID,
		GrpcAddress:   p.GRPCAddress,
		Level:         p.Level,
	}
}

// NewPeers creates a new peers instance.
func NewPeers(peers ...*Peer) *Peers {
	return &Peers{
		mutex: sync.RWMutex{},
		Peers: peers,
	}
}

// GetAll returns a list of peers.
func (p *Peers) GetAll() []*Peer {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	peers := make([]*Peer, len(p.Peers))
	copy(peers, p.Peers)

	return peers
}

// Add adds a peer to the list.
func (p *Peers) Add(peer *Peer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i := 0; i < len(p.Peers); i++ {
		if p.Peers[i].GRPCAddress == peer.GRPCAddress {
			p.Peers[i] = peer
			return
		}
	}

	p.Peers = append(p.Peers, peer)
}

func (p *Peers) ToProto() []*types.Peer {
	peers := make([]*types.Peer, len(p.Peers))

	for i, peer := range p.Peers {
		peers[i] = peer.ToProto()
	}

	return peers
}

// Exists checks if a peer exists in the list.
func (p *Peers) Exists(peer *Peer) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, p := range p.Peers {
		if p.GRPCAddress == peer.GRPCAddress {
			return true
		}
	}

	return false
}
