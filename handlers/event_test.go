package handlers

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
)

var (
	logs = []byte(`[2021-09-29T20:54:51Z] [info] [ThetaEdgeLauncher] [2021-09-29T20:54:51Z]  INFO [p2p] Already has sufficient number of peers, numPeers: 16, sufficientNumPeers: 16
	[2021-09-29T20:54:57Z] [info] [ThetaEdgeLauncher] [2021-09-29T20:54:57Z]  INFO [p2p] Already has sufficient number of peers, numPeers: 16, sufficientNumPeers: 16
	[2021-08-28 09:00:26.951] [info] [ThetaEdgeLauncher] [2021-08-28 09:00:26]  INFO [uptime miner] Broadcasted vote: EENVote{Block: 0x6d0ae6972cd670a8f7dfd628ef516051d0fd699906c55f80cff12540bd3786a8, Height: 11759001, Address: 0x8d25fa2e7d, Signature: E1A06D0AE697786A8, CreationTimestamp: 1630155386}
	`)
)

func TestProcessEvent(t *testing.T) {
	// setup test variables
	var gotOffset int64
	var wantOffset int64
	var prevOffset int64
	var event fsnotify.Event
	var buf []byte
	var err error

	// create tmp file with data
	buf = logs
	tmp := t.TempDir()
	fp := filepath.Join(tmp, "log.log")
	_ = os.WriteFile(fp, buf, 0664)

	f, _ := os.Open(fp)
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()

	watcher, _ := fsnotify.NewWatcher()
	_ = watcher.Add(fp)
	defer watcher.Close()

	// test process write event
	event = fsnotify.Event{Op: fsnotify.Write}
	prevOffset = 0
	wantOffset = size
	gotOffset, err = ProcessEvent(watcher, event, fp, prevOffset)
	if err != nil {
		t.Fatalf("handlers.ProcessEvent() returned error: %v", err)
	}

	if gotOffset != wantOffset {
		t.Fatalf("handlers.ProcessEvent() returned: %v, wanted: %v", gotOffset, wantOffset)
	}

	// test process rename event
	fpo := filepath.Join(tmp, "log.old.log")
	_ = os.Rename(fp, fpo)          // rename log.log to log.old.log
	_ = os.WriteFile(fp, buf, 0664) // create new log.log

	event = fsnotify.Event{Op: fsnotify.Rename}
	prevOffset = size * 2
	wantOffset = size
	gotOffset, err = ProcessEvent(watcher, event, fp, prevOffset)
	if err != nil {
		t.Fatalf("handlers.ProcessEvent() returned error: %v", err)
	}

	if gotOffset != wantOffset {
		t.Fatalf("handlers.ProcessEvent() returned: %v, wanted: %v", gotOffset, wantOffset)
	}
}

func TestProcessLog(t *testing.T) {
	// setup test variables
	var gotOffset int64
	var wantOffset int64
	var prevOffset int64
	var buf []byte
	var err error

	// create tmp file with data
	buf = logs
	tmp := t.TempDir()
	fp := filepath.Join(tmp, "log.log")
	_ = os.WriteFile(fp, buf, 0664)

	f, _ := os.Open(fp)
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()

	// test process log
	prevOffset = 0
	wantOffset = size
	gotOffset, err = processLog(fp, prevOffset)
	if err != nil {
		t.Fatalf("handlers.processLog() returned error: %v", err)
	}

	if gotOffset != wantOffset {
		t.Fatalf("handlers.ProcessLog() returned: %v, wanted: %v", gotOffset, wantOffset)
	}

	// test process error
	//
}

func TestGetOffset(t *testing.T) {
	// setup test variables
	var gotOffset int64
	var wantOffset int64
	var gotSize int64
	var wantSize int64
	var prevOffset int64
	var buf []byte
	var err error

	// create tmp file with data
	buf = logs
	tmp := t.TempDir()
	fp := filepath.Join(tmp, "log.log")
	_ = os.WriteFile(fp, buf, 0664)

	// test offset < file size
	prevOffset = 0
	f, _ := os.Open(fp)
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()

	wantOffset, wantSize = prevOffset, size
	gotOffset, gotSize, err = getOffset(f, prevOffset)
	if err != nil {
		t.Fatalf("handlers.GetOffset() returned error: %v", err)
	}

	if gotOffset != wantOffset || gotSize != wantSize {
		t.Fatalf("handlers.GetOffset() returned: %v, %v, wanted: %v, %v", gotOffset, gotSize, wantOffset, wantSize)
	}

	// test offset > file size
	prevOffset = size * 2
	f, _ = os.Open(fp)
	defer f.Close()
	info, _ = f.Stat()
	size = info.Size()

	wantOffset, wantSize = 0, size
	gotOffset, gotSize, err = getOffset(f, prevOffset)
	if err != nil {
		t.Fatalf("handlers.GetOffset() returned error: %v", err)
	}

	if gotOffset != wantOffset || gotSize != wantSize {
		t.Fatalf("handlers.GetOffset() returned: %v, %v, wanted: %v, %v", gotOffset, gotSize, wantOffset, wantSize)
	}

	// test none file path
	fp = filepath.Join(tmp, "error.log")
	f, _ = os.Open(fp)
	defer f.Close()

	gotOffset, gotSize, err = getOffset(f, prevOffset)
	if err == nil {
		t.Fatalf("handlers.GetOffset() returned: %v, %v, wanted error: %v", gotOffset, gotSize, err)
	}
}

func TestGetFilePath(t *testing.T) {
	// setup test variables
	var got string
	var want string
	var rt string
	var hd string
	var err error

	u, _ := user.Current()
	hd = u.HomeDir

	// test darwin file path
	rt = "darwin"
	want = fmt.Sprintf("%s/Library/Logs/Theta Edge Node/log.log", hd)
	got, err = GetFilePath(rt)
	if err != nil {
		t.Fatalf("handlers.GetFilePath() returned error: %v", err)
	}

	if got != want {
		t.Fatalf("handlers.GetFilePath() returned: %v, wanted: %v", got, want)
	}

	// test windows file path
	rt = "windows"
	want = fmt.Sprintf("%s\\AppData\\Roaming\\Theta Edge Node\\log.log", hd)
	got, err = GetFilePath(rt)
	if err != nil {
		t.Fatalf("handlers.GetFilePath() returned error: %v", err)
	}

	if got != want {
		t.Fatalf("handlers.GetFilePath() returned: %v, wanted: %v", got, want)
	}

	// // test linux file path
	// rt = "linux"
	// want = fmt.Sprintf("%s/path/to/logfile/log.log", hd)
	// got, err = GetFilePath(rt)
	// if err != nil {
	// 	t.Fatalf("handlers.GetFilePath() returned error: %v", err)
	// }

	// if got != want {
	// 	t.Fatalf("handlers.GetFilePath() returned: %v, wanted: %v", got, want)
	// }

	// test none OS file path
	rt = "breebsd"
	want = ""
	got, err = GetFilePath(rt)
	if err == nil {
		t.Fatalf("handlers.GetFilePath() returned: %v, wanted error: %v", got, err)
	}
}

func TestPokeFilePath(t *testing.T) {
	var fp string
	var tmp string
	var err error

	// test existing file path
	tmp = t.TempDir()
	fp = filepath.Join(tmp, "log.log")
	os.Create(fp)
	if err = PokeFilePath(fp); err != nil {
		t.Fatalf("handlers.PokeFilePath() returned error: %v", err)
	}

	// test not extst file path
	fp = ""
	if err = PokeFilePath(fp); err == nil {
		t.Fatalf("handlers.PokeFilePath() returned: %v, wanted error: %v", nil, err)
	}
}
