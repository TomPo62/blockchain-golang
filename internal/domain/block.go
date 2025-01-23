package domain

type Block struct {
	Index int64
	Timestamp int64
	PreviousHash string
	Hash string
	Data string
	Nonce int
}
