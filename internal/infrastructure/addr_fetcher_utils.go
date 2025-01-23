package infrastructure

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func FetchBootstrapNodes(apiURLs []string) []string {
	var bootstrapNodes []string
	for _, url := range apiURLs {
		resp, err := http.Get(fmt.Sprintf("%s/multiaddr", url))
		if err != nil {
			log.Printf("Failed to fetch bootstrap node from %s: %v", url, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			continue
		}

		addr := strings.TrimSpace(string(body))
		addr = strings.Replace(addr, "127.0.0.1", url[7:], 1)
		bootstrapNodes = append(bootstrapNodes, addr)
	}

	return bootstrapNodes
}


