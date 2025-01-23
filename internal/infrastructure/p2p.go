package infrastructure

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type P2PNode struct {
	Host      host.Host
	Peers     map[string]peer.AddrInfo
	PeersLock sync.Mutex
}

func (node *P2PNode) ConnectToPeer(peerAddr string) {
	log.Printf("Attempting to connect to peer: %s\n", peerAddr)
	peerInfo, err := peer.AddrInfoFromString(peerAddr)
	if err != nil {
		log.Printf("Invalid peer address: %v", err)
		return
	}

	if peerInfo.ID == node.Host.ID() {
		log.Println("Skipping self-connection")
		return
	}

	node.PeersLock.Lock()
	if _, exists := node.Peers[peerInfo.ID.String()]; exists {
		node.PeersLock.Unlock()
		log.Println("Peer already connected")
		return
	}
	node.PeersLock.Unlock()

	if err := node.Host.Connect(context.Background(), *peerInfo); err != nil {
		log.Printf("Failed to connect to peer: %v", err)
		return
	}

	node.PeersLock.Lock()
	node.Peers[peerInfo.ID.String()] = *peerInfo
	node.PeersLock.Unlock()

	log.Printf("Connected to peer: %v\n", peerInfo.ID)

	go node.SharePeers(peerInfo.ID)
}

func (node *P2PNode) ReceivePeers(peerData string) {
	peerAddrs := strings.Split(peerData, "\n") // Chaque peer est séparé par une ligne
	for _, peerAddr := range peerAddrs {
		if strings.TrimSpace(peerAddr) == "" {
			continue
		}

		peerInfo, err := peer.AddrInfoFromString(peerAddr)
		if err != nil {
			log.Printf("Invalid peer data received: %v", err)
			continue
		}

		node.PeersLock.Lock()
		if _, exists := node.Peers[peerInfo.ID.String()]; !exists {
			node.Peers[peerInfo.ID.String()] = *peerInfo
			log.Printf("New peer added: %s", peerInfo.ID)
		}
		node.PeersLock.Unlock()
	}
}



func (node *P2PNode) SharePeers(peerID peer.ID) {
	node.PeersLock.Lock()
	defer node.PeersLock.Unlock()

	for _, info := range node.Peers {
		multiaddr := fmt.Sprintf("%s/p2p/%s", info.Addrs[0], info.ID) // Multiaddr correcte
		stream, err := node.Host.NewStream(context.Background(), peerID, "/blockchain/1.0.0")
		if err != nil {
			log.Printf("Error creating stream to share peers: %v", err)
			continue
		}
		defer stream.Close()

		log.Printf("Sharing peer multiaddr: %s with %s", multiaddr, peerID)
		stream.Write([]byte(multiaddr + "\n")) // Ajoute une nouvelle ligne pour distinguer les peers
	}
}


func (node *P2PNode) DiscoverPeers() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		node.PeersLock.Lock()
		for _, peerInfo := range node.Peers {
			go node.ConnectToPeer(peerInfo.String())
		}
		node.PeersLock.Unlock()
	}
}


func (node *P2PNode) GetMultiaddr() string {
	for _, addr := range node.Host.Addrs() {
		return fmt.Sprintf("%s/p2p/%s", addr, node.Host.ID())
	}
	return ""
}
