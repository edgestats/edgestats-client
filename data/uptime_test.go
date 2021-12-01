package data

import (
	"reflect"
	"testing"
	"time"
)

var (
	UMReceivedEx    = []byte("[2021-08-28 09:17:11.831] [info] [ThetaEdgeLauncher] [2021-08-28 09:17:11]  INFO [uptime miner] Received block: 0xfdca353dd0dcb8d193b1d731e065db547dcdfc6b0af20efdabb1eeff0f430cf2, height: 11759201, epoch: 11841048")
	UMNewBlockEx    = []byte("[2021-08-28 09:17:11.836] [info] [ThetaEdgeLauncher] [2021-08-28 09:17:11]  INFO [uptime miner] Start new block, block: 0xbecf60, height: 11759201, vote: EENVote{Block: 0xfdca353dd0dcb8d193b1d731e065db547dcdfc6b0af20efdabb1eeff0f430cf2, Height: 11759201, Address: 0x8d25fa2e7d, Signature: E1A06D0AE697786A8, CreationTimestamp: 1630156631}")
	UMNewRoundEx    = []byte("[2021-08-28 09:10:18.757] [info] [ThetaEdgeLauncher] [2021-08-28 09:10:18]  INFO [uptime miner] Start new round: 7")
	UMBroadcastedEx = []byte("[2021-08-28 09:00:26.951] [info] [ThetaEdgeLauncher] [2021-08-28 09:00:26]  INFO [uptime miner] Broadcasted vote: EENVote{Block: 0x6d0ae6972cd670a8f7dfd628ef516051d0fd699906c55f80cff12540bd3786a8, Height: 11759001, Address: 0x8d25fa2e7d, Signature: E1A06D0AE697786A8, CreationTimestamp: 1630155386}")
)

func TestNewUMBroadcast(t *testing.T) {
	want := &UMBroadcast{}
	got := NewUMBroadcast()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("data.NewUMBroadcast() returned: %v, wanted: %v", got, want)
	}
}

func TestUMBroadcastToJSON(t *testing.T) {}

func TestUMBroadcastParse(t *testing.T) {
	// setup test variables
	var log []byte
	var got = NewUMBroadcast()
	var want = NewUMBroadcast()
	var tt time.Time
	var err error

	// sim bootstrap peers
	setPeers(16, 16)

	// test parse broadcastedVote
	log = UMBroadcastedEx
	tt, _ = parseTime(log)
	want = &UMBroadcast{
		Block:           "0x6d0ae6972cd670a8f7dfd628ef516051d0fd699906c55f80cff12540bd3786a8",
		Height:          11759001,
		Addr:            "0x8d25fa2e7d",
		Signature:       "E1A06D0AE697786A8",
		Timestamp:       1630155386,
		NumPeers:        16,
		SufficientPeers: 16,
		CreatedAt:       tt,
	}

	if err = got.Parse(log); err != nil {
		t.Fatalf("some error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("data.UMBroadcastParse() returned: %v, wanted: %v", got, want)
	}

	// test parse newBlock
	log = UMNewBlockEx
	tt, _ = parseTime(log)
	want = &UMBroadcast{
		Block:           "0xfdca353dd0dcb8d193b1d731e065db547dcdfc6b0af20efdabb1eeff0f430cf2",
		Height:          11759201,
		Addr:            "0x8d25fa2e7d",
		Signature:       "E1A06D0AE697786A8",
		Timestamp:       1630156631,
		NumPeers:        16,
		SufficientPeers: 16,
		CreatedAt:       tt,
	}

	if err = got.Parse(log); err != nil {
		t.Fatalf("some error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("data.UMBroadcastParse() returned: %v, wanted: %v", got, want)
	}

	// test no peers bootstrapped error
	setPeers(0, 0)

	if err = got.Parse(log); err == nil {
		t.Fatalf("data.UMBroadcastParse() returned: %v, wanted error: %v", got, err)
	}
}

func TestFilterUMLogType(t *testing.T) {
	// setup test variables
	var log []byte
	var got int
	var want int

	// test receivedBlock
	log = UMReceivedEx
	want = umReceivedBlockFilter
	got = filterUMLogType(log)

	if got != want {
		t.Fatalf("data.filterUMLogType() returned: %v, wanted: %v", got, want)
	}

	// test NewBlock
	log = UMNewBlockEx
	want = umBroadcastedVoteFilter // umNewBlockFilter
	got = filterUMLogType(log)

	if got != want {
		t.Fatalf("data.filterUMLogType() returned: %v, wanted: %v", got, want)
	}

	// test newRound
	log = UMNewRoundEx
	want = umNewRoundFilter
	got = filterUMLogType(log)

	if got != want {
		t.Fatalf("data.filterUMLogType() returned: %v, wanted: %v", got, want)
	}

	// test BroadcastedVote
	log = UMBroadcastedEx
	want = umBroadcastedVoteFilter
	got = filterUMLogType(log)

	if got != want {
		t.Fatalf("data.filterUMLogType() returned: %v, wanted: %v", got, want)
	}

	// test no match filter
	log = []byte("no match log")
	want = umErrFilter
	got = filterUMLogType(log)

	if got != want {
		t.Fatalf("data.filterUMLogType() returned: %v, wanted: %v", got, want)
	}
}

func TestGetUMLogURL(t *testing.T) {}

func TestSetAddr(t *testing.T) {
	// setup test variables
	var got string
	var want = "Ox8aef2c43"

	// test null address
	if err := setAddr(""); err == nil {
		t.Fatalf("data.setAddr() returned: %v, wanted error: %v", got, err)
	}

	// test set address
	if err := setAddr(want); err != nil {
		t.Fatalf("data.setAddr() returned error: %v", err)
	}

	got = nodeAddr
	if got != want {
		t.Fatalf("data.setAddr() returned: %v, wanted: %v", got, want)
	}
}
