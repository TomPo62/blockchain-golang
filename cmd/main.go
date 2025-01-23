package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/TomPo62/blockchain-golang/internal/delivery"
	"github.com/TomPo62/blockchain-golang/internal/infrastructure"
)

func getBootstrapNodes() []string {
	apiURLs := []string{
		"http://node1:9001",
		"http://node2:9002",
	}

	bootstrapNodes := infrastructure.FetchBootstrapNodes(apiURLs)
	log.Printf("Bootstrap nodes: %v", bootstrapNodes)
	return bootstrapNodes
}

func main() {
	nodePort := 8000
	if len(os.Args) > 1 {
		p, err := strconv.Atoi(os.Args[1])
		if err == nil {
			nodePort = p
		}
	}

	log.Printf("Starting node on port %d...\n", nodePort)
	node := infrastructure.StartP2PNode(nodePort)

	// Démarre le serveur REST
	apiPort := nodePort + 1000
	server := &delivery.Server{
		Node: node,
	}
	go server.StartHTTPServer(apiPort)

	// Attends que le serveur REST soit prêt (petit délai pour la stabilité)
	time.Sleep(2 * time.Second)

	// Récupère les bootstrap nodes
	bootstrapNodes := getBootstrapNodes()

	// Connexion aux bootstrap nodes
	for _, peerAddr := range bootstrapNodes {
		go node.ConnectToPeer(peerAddr)
	}

	// Démarre la découverte des peers
	go node.DiscoverPeers()

	select {}
}
