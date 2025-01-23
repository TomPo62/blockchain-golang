package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TomPo62/blockchain-golang/internal/infrastructure"
)

type Server struct {
	Node *infrastructure.P2PNode
}

// Démarre le serveur HTTP
func (s *Server) StartHTTPServer(port int) {
	http.HandleFunc("/peers", s.GetPeers)
	http.HandleFunc("/multiaddr", s.GetMultiaddr)
	http.HandleFunc("/connect", s.ConnectPeer)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP server on port %d...\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// Handler pour connecter un peer
func (s *Server) ConnectPeer(w http.ResponseWriter, r *http.Request) {
	type ConnectRequest struct {
		PeerAddr string `json:"peer_addr"`
	}

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.PeerAddr == "" {
		http.Error(w, "Peer address is required", http.StatusBadRequest)
		return
	}

	go s.Node.ConnectToPeer(req.PeerAddr)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Connecting to peer: %s\n", req.PeerAddr)
}


// Handler pour récupérer la liste des peers
func (s *Server) GetPeers(w http.ResponseWriter, r *http.Request) {
	s.Node.PeersLock.Lock()
	defer s.Node.PeersLock.Unlock()

	peers := []string{}
	for _, peerInfo := range s.Node.Peers {
		peers = append(peers, peerInfo.String())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

// Handler pour récupérer le multiaddr complet du node
func (s *Server) GetMultiaddr(w http.ResponseWriter, r *http.Request) {
	multiaddr := s.Node.GetMultiaddr()

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", multiaddr)
}
