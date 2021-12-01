package data

import (
	"reflect"
	"testing"
	"time"
)

var (
	P2PNumPeersEx = []byte("[2021-08-28 09:10:32.888] [info] [ThetaEdgeLauncher] [2021-08-28 09:10:32]  INFO [p2p] Already has sufficient number of peers, numPeers: 16, sufficientNumPeers: 16")
)

func TestNewP2PNumPeers(t *testing.T) {
	want := &P2PNumPeers{}
	got := NewP2PNumPeers()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("data.NewP2PNumPeers() returned: %v, wanted: %v", got, want)
	}
}

func TestP2PNumPeersToJSON(t *testing.T) {}

func TestP2PNumPeersParse(t *testing.T) {
	// setup test variables
	var log []byte
	var got = NewP2PNumPeers()
	var want = NewP2PNumPeers()
	var tt time.Time
	var err error

	// sim bootstrap address
	setAddr("0x8d25fa2e7d")

	// test parse numPeers
	log = P2PNumPeersEx
	tt, _ = parseTime(log)
	want = &P2PNumPeers{
		Addr: "0x8d25fa2e7d",
		// Time:            tt,
		NumPeers:        16,
		SufficientPeers: 16,
		CreatedAt:       tt,
	}

	if err = got.Parse(log); err != nil {
		t.Fatalf("some error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("data.P2PNumPeersParse() returned: %v, wanted: %v", got, want)
	}

	// test no address bootstrapped error
	setAddr("")

	if err = got.Parse(log); err == nil {
		t.Fatalf("data.P2PNumPeersParse() returned: %v, wanted error: %v", got, err)
	}
}

func TestFilterP2PLogType(t *testing.T) {
	// setup test variables
	var log []byte
	var got int
	var want int

	// test numPeers
	log = P2PNumPeersEx
	want = p2pNumPeersFilter
	got = filterP2PLogType(log)

	if got != want {
		t.Fatalf("data.filterP2PLogType() returned: %v, wanted: %v", got, want)
	}

	// test no match filter
	log = []byte("no match log")
	want = p2pErrFilter
	got = filterP2PLogType(log)

	if got != want {
		t.Fatalf("data.filterP2PLogType() returned: %v, wanted: %v", got, want)
	}
}

func TestGetP2PLogURL(t *testing.T) {}

func TestSetPeers(t *testing.T) {
	// setup test variables
	var gotNodePeers int
	var gotSufficientPeers int
	var wantNodePeers = 16
	var wantSufficientPeers = 16

	// test zero sufficient peers
	if err := setPeers(0, 0); err == nil {
		t.Fatalf("data.setPeers() returned: %v, %v, wanted error: %v", gotNodePeers, gotSufficientPeers, err)
	}

	// test set peers
	if err := setPeers(wantNodePeers, wantSufficientPeers); err != nil {
		t.Fatalf("data.setPeers() returned error: %v", err)
	}

	gotNodePeers, gotSufficientPeers = nodePeers, sufficientPeers
	if gotNodePeers != wantNodePeers || gotSufficientPeers != wantSufficientPeers {
		t.Fatalf("data.setPeers() returned: %v, %v wanted: %v, %v", gotNodePeers, gotSufficientPeers, wantNodePeers, wantSufficientPeers)
	}
}
