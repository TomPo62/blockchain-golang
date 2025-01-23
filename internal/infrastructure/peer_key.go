package infrastructure

import (
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func LoadOrCreateKey(path string) (crypto.PrivKey, error) {
	// Si le fichier de clé existe, charge la clé
	if _, err := os.Stat(path); err == nil {
		keyBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		privKey, err := crypto.UnmarshalPrivateKey(keyBytes)
		if err != nil {
			return nil, err
		}
		return privKey, nil
	}

	// Sinon, génère une nouvelle clé et la sauvegarde
	privKey, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	keyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, keyBytes, 0600); err != nil {
		return nil, err
	}

	return privKey, nil
}

func StartP2PNode(port int) *P2PNode {
	privKey, err := LoadOrCreateKey("node.key")
	if err != nil {
		log.Fatalf("Failed to load or create key: %v", err)
	}

	host, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.Identity(privKey),
	)
	if err != nil {
		log.Fatalf("Failed to create P2P node: %v", err)
	}

	node := &P2PNode{
		Host:  host,
		Peers: make(map[string]peer.AddrInfo),
	}

	host.SetStreamHandler("/blockchain/1.0.0", func(stream network.Stream) {
		log.Println("New stream opened")
		defer stream.Close()

		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			log.Printf("Error reading from stream: %v", err)
			return
		}

		peerData := string(buf[:n])
		log.Printf("Received peer data: %s", peerData)

		// Met à jour la liste des peers avec les données reçues
		node.ReceivePeers(peerData)
	})

	log.Printf("Node started. Listening on: %v\n", host.Addrs())
	log.Printf("Node multiaddr: %s\n", node.GetMultiaddr())
	return node
}
