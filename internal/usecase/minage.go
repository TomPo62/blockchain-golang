package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/TomPo62/blockchain-golang/internal/domain"
)

func MineBlock(previousBlock domain.Block, data string) domain.Block {
	var newBlock domain.Block
	newBlock.Index = previousBlock.Index + 1
	newBlock.Timestamp = time.Now().Unix()
	newBlock.PreviousHash = previousBlock.Hash
	newBlock.Data = data

	nonce := 0
	for {
		newBlock.Nonce = nonce
		hash := calculateHash(newBlock)
		if isValidHash(hash) {
			newBlock.Hash = hash
			break
		}
		nonce++
	}

	fmt.Printf("Mined new block: %+v\n", newBlock)
	return newBlock
}

func calculateHash(block domain.Block) string {
	record := strconv.FormatInt(block.Index, 10) + strconv.FormatInt(block.Timestamp, 10) + block.PreviousHash + block.Data + strconv.Itoa(block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isValidHash(hash string) bool {
	// Exemple de difficulté : le hash doit commencer par 4 zéros
	return hash[:4] == "0000"
}
