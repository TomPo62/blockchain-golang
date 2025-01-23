package usecase

import (
	"strings"

	"github.com/TomPo62/blockchain-golang/internal/domain"
)

func IsValidBlock(newBlock, oldBlock domain.Block) bool {
	if newBlock.PreviousHash != oldBlock.Hash {
		return false
	}
	if !strings.HasPrefix(newBlock.Hash, "0000") {
		return false
	}
	return true
}
