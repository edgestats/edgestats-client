package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	p2pNumPeersRE = `(?m)(?P<key>\w+):\s+(?P<value>\w+)\,?`
)

const (
	p2pErrFilter = iota
	p2pNumPeersFilter
	p2pInternalAddrFilter
)

var (
	p2pNumPeersServiceURL = fmt.Sprintf("%s/stats/uptimes/peers", apiAddr)
	numPeers              = []byte("numPeers")
)

var (
	nodeAddr        string
	nodePeers       int
	sufficientPeers int
)

type P2PNumPeers struct {
	Addr            string    `json:"address"`
	NumPeers        int       `json:"num_peers"`
	SufficientPeers int       `json:"sufficient_peers"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewP2PNumPeers() *P2PNumPeers {
	return &P2PNumPeers{}
}

func (p2p *P2PNumPeers) ToJSON() ([]byte, error) {
	return json.Marshal(p2p)
}

func (p2p *P2PNumPeers) Parse(b []byte) error {
	if filterP2PLogType(b) != p2pNumPeersFilter {
		return errors.New("no match")
	}

	re := regexp.MustCompile(p2pNumPeersRE)
	ml := re.FindAllSubmatch(b, -1)

	var np int
	var sp int

	for _, i := range ml {
		switch true {
		case bytes.Equal(i[1], []byte("numPeers")):
			v, err := strconv.Atoi(string(i[2]))
			if err != nil {
				return err
			}
			np = v
		case bytes.Equal(i[1], []byte("sufficientNumPeers")):
			v, err := strconv.Atoi(string(i[2]))
			if err != nil {
				return err
			}
			sp = v
		default:
			s := fmt.Sprintf("error no match: %s, found: %s", i[1], i[2])
			return errors.New(s)
		}
	}

	t, err := parseTime(b)
	if err != nil {
		return err
	}

	// populate numPeers and sufficientPeers variables
	if err := setPeers(np, sp); err != nil {
		return err
	}

	// return error if node not yet bootstrapped with address
	if nodeAddr == "" {
		return errors.New("no address bootstrapped")
	}

	p2p.Addr = nodeAddr
	p2p.NumPeers = np
	p2p.SufficientPeers = sp
	p2p.CreatedAt = t

	return nil
}

func filterP2PLogType(b []byte) int {
	// filter log by subcategory
	switch true {
	case bytes.Contains(b, numPeers):
		return p2pNumPeersFilter
	default:
		return p2pErrFilter
	}
}

func setPeers(num, suff int) error {
	nodePeers = num
	sufficientPeers = suff
	if sufficientPeers == 0 {
		return errors.New("no peers bootstrapped")
	}

	return nil
}
