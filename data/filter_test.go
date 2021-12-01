package data

import "testing"

var (
	UMLogExample  = []byte("... [uptime miner] ...")
	NSLogExample  = []byte("... [netsync] ...")
	P2PLogExample = []byte("... [p2p] ...")
	ECLogExample  = []byte(`{"user": "Anonymous", ...`)
	ErrLogExample = []byte("... [unmatched] ...")
)

func TestFilter(t *testing.T) {
	// setup test vars
	var log []byte
	var got int
	var want int

	// test uptime miner log
	log = UMLogExample
	want = UMFilter
	got = Filter(log)

	if got != want {
		t.Fatalf("data.Filter() returned: %v, wanted: %v", got, want)
	}

	// test p2p log
	log = P2PLogExample
	want = P2PFilter
	got = Filter(log)

	if got != want {
		t.Fatalf("data.Filter() returned: %v, wanted: %v", got, want)
	}

	// test no match
	log = ErrLogExample
	want = ErrFilter
	got = Filter(log)

	if got != want {
		t.Fatalf("data.Filter() returned: %v, wanted: %v", got, want)
	}
}
