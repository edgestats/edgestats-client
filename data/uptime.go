package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	umBroadcastRE = `(?m)(?P<key>\w+):\s+(?P<value>\w+)\,?`
)

const (
	umErrFilter = iota
	umReceivedBlockFilter
	umNewBlockFilter
	umNewRoundFilter
	umBroadcastedVoteFilter
)

var (
	umBroadcastedServiceURL = fmt.Sprintf("%s/stats/uptimes/broadcasts", apiAddr)
	receivedBlock           = []byte("Received block")
	newBlock                = []byte("Start new block")
	newRound                = []byte("Start new round")
	broadcastedVote         = []byte("Broadcasted vote")
)

type UMBroadcast struct {
	Block           string    `json:"block"`
	Height          int       `json:"height"`
	Addr            string    `json:"address"`
	Signature       string    `json:"signature"`
	Timestamp       int       `json:"timestamp"`
	NumPeers        int       `json:"num_peers"`
	SufficientPeers int       `json:"sufficient_peers"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewUMBroadcast() *UMBroadcast {
	return &UMBroadcast{}
}

func (um *UMBroadcast) ToJSON() ([]byte, error) {
	return json.Marshal(um)
}

func (um *UMBroadcast) Parse(b []byte) error {
	if filterUMLogType(b) != umBroadcastedVoteFilter {
		return errors.New("no match")
	}

	re := regexp.MustCompile(umBroadcastRE)
	ml := re.FindAllSubmatch(b, -1)

	var vk string
	var vh int
	var va string
	var vs string
	var vt int

	for _, i := range ml {
		switch true {
		case bytes.Equal(i[1], []byte("vote")):
			continue // perhaps check that value is EENVote?
		case bytes.Equal(i[1], []byte("Block")):
			vk = string(i[2])
		case bytes.Equal(i[1], []byte("Height")):
			v, err := strconv.Atoi(string(i[2]))
			if err != nil {
				return err
			}
			vh = v
		case bytes.Equal(i[1], []byte("Address")):
			v := string(i[2])
			va = strings.ToLower(v)
		case bytes.Equal(i[1], []byte("Signature")):
			vs = string(i[2])
		case bytes.Equal(i[1], []byte("CreationTimestamp")):
			v, err := strconv.Atoi(string(i[2]))
			if err != nil {
				return err
			}
			vt = v
		case bytes.Equal(i[1], []byte("block")):
			continue // only seen in start new block log
		case bytes.Equal(i[1], []byte("height")):
			continue // only seen in start new block log
		default: // no match found
			s := fmt.Sprintf("error no match: %s, found: %s", i[1], i[2])
			return errors.New(s)
		}
	}

	t, err := parseTime(b)
	if err != nil {
		return err
	}

	// populate node address variable
	if err := setAddr(va); err != nil {
		return err
	}

	// return error if not yet bootstrapped with peers
	// perhaps rethink this condition
	if nodePeers == 0 && sufficientPeers == 0 {
		return errors.New("no peers bootstrapped")
	}

	um.Block = vk
	um.Height = vh
	um.Addr = va
	um.Signature = vs
	um.Timestamp = vt
	um.NumPeers = nodePeers
	um.SufficientPeers = sufficientPeers
	um.CreatedAt = t

	return nil
}

func filterUMLogType(b []byte) int {
	// filter log by subcategory
	switch true {
	case bytes.Contains(b, receivedBlock):
		return umReceivedBlockFilter
	case bytes.Contains(b, newBlock):
		return umBroadcastedVoteFilter
	case bytes.Contains(b, newRound):
		return umNewRoundFilter
	case bytes.Contains(b, broadcastedVote):
		return umBroadcastedVoteFilter
	default:
		return umErrFilter
	}
}

func setAddr(addr string) error {
	nodeAddr = addr
	if nodeAddr == "" {
		return errors.New("no address bootstrapped")
	}

	return nil
}
