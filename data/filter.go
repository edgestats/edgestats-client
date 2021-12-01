package data

import "bytes"

const (
	ErrFilter = iota
	UMFilter
	P2PFilter
)

var (
	uptime = []byte("[uptime miner")
	p2p    = []byte("[p2p]")
)

func Filter(b []byte) int {
	// filter log by category
	switch true {
	case bytes.Contains(b, uptime):
		return UMFilter
	case bytes.Contains(b, p2p):
		return P2PFilter
	default:
		return ErrFilter
	}
}
