package data

import (
	"testing"
	"time"
)

var (
	P2PLogEx01 = []byte("[2021-08-28 09:10:32.888] [info] [ThetaEdgeLauncher] [2021-08-28 09:10:32] ...")
	P2PLogEx02 = []byte("[2021-08-28 09:10:32] [info] [ThetaEdgeLauncher] [2021-08-28 09:10:32] ...")
	ErrLogEx   = []byte("[2021-08-28  09:10:32] [info] [ThetaEdgeLauncher] [2021-08-28  09:10:32] ...")
)

func TestParseTime(t *testing.T) {
	// setup test vars
	var log []byte
	var got time.Time
	var want time.Time
	var err error
	var loc = time.Now().Location()
	var utc = time.Now().UTC().Location()

	// test uptime miner log
	log = UMBroadcastedEx
	want, _ = time.ParseInLocation("2006-01-02T15:04:05.999", "2021-08-28T09:00:26.951", loc)
	want = want.In(utc)

	got, err = parseTime(log)
	if err != nil {
		t.Fatalf("data.parseTime() returned: %v", err)
	}

	if got != want {
		t.Fatalf("data.parseTime() returned: %v, wanted: %v", got, want)
	}

	// test p2p log
	log = P2PLogEx01
	want, _ = time.ParseInLocation("2006-01-02T15:04:05.999", "2021-08-28T09:10:32.888", loc)
	want = want.In(utc)

	got, err = parseTime(log)
	if err != nil {
		t.Fatalf("data.parseTime() returned: %v", err)
	}

	if got != want {
		t.Fatalf("data.parseTime() returned: %v, wanted: %v", got, want)
	}

	// test p2p log no milliseconds
	log = P2PLogEx02
	want, _ = time.ParseInLocation("2006-01-02T15:04:05", "2021-08-28T09:10:32", loc)
	want = want.In(utc)

	got, err = parseTime(log)
	if err != nil {
		t.Fatalf("data.parseTime() returned: %v", err)
	}

	if got != want {
		t.Fatalf("data.parseTime() returned: %v, wanted: %v", got, want)
	}

	// test time string error
	log = ErrLogEx

	got, err = parseTime(log)
	if err == nil {
		t.Fatalf("data.parseTime() returned: %v, wanted error: %v", got, err)
	}
}

func TestGetServiceURI(t *testing.T) {
	// setup test vars
	var got string
	var want string
	var p Parser

	// test uptime broadcast service uri
	p = NewUMBroadcast()
	want = umBroadcastedServiceURL
	got = getServiceURI(p)
	if got != want {
		t.Fatalf("data.getServiceURI() returned: %v, wanted: %v", got, want)
	}

	// test p2p numpeers service uri
	p = NewP2PNumPeers()
	want = p2pNumPeersServiceURL
	got = getServiceURI(p)
	if got != want {
		t.Fatalf("data.getServiceURI() returned: %v, wanted: %v", got, want)
	}
}

func TestFuzzRequest(t *testing.T) {
	// rethink how to test this
	fuzzRequest()
}
