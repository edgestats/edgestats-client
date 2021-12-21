package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/edgestats/edgestats-client/data"
	"github.com/fsnotify/fsnotify"
)

func ProcessEvent(watcher *fsnotify.Watcher, event fsnotify.Event, fp string, offset int64) (int64, error) {
	var err error

	if event.Op&fsnotify.Write == fsnotify.Write {
		offset, err = processLog(fp, offset)
		if err != nil {
			return offset, err
		}
	}

	if event.Op&fsnotify.Rename == fsnotify.Rename {
		// allow edge node log rotate process to complete
		time.Sleep(1000 * time.Millisecond)

		// scan writes to new "log.old.log", ie old "log.log"
		fpOld := fp[:len(fp)-3] + "old.log"
		offset, err = processLog(fpOld, offset)
		if err != nil {
			return offset, err
		}

		// stop watching new "log.old.log", ie old "log.log"
		if err := watcher.Remove(fp); err != nil {
			return offset, err
		}

		// start watching new log file, ie new "log.log"
		if err := watcher.Add(fp); err != nil {
			return offset, err
		}

		// scan writes to new log file, ie new "log.log"
		offset, err = processLog(fp, offset) // not neccessary to set offset = 0
		if err != nil {
			return offset, err
		}
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		// should not have anything to do
	}

	if event.Op&fsnotify.Create == fsnotify.Create {
		// should not have anything to do
	}

	if event.Op&fsnotify.Chmod == fsnotify.Chmod {
		// should not have anything to do
	}

	return offset, err
}

func processLog(fp string, offset int64) (int64, error) {
	var lines int
	var matches int
	var misses int
	var size int64

	f, err := os.Open(fp)
	if err != nil {
		return offset, err
	}
	defer f.Close()

	// offset either last file size or file start
	offset, size, err = getOffset(f, offset)
	if err != nil {
		return offset, err
	}

	// seek to offset in file
	ro := 0 // offset relative to file origin
	_, err = f.Seek(offset, ro)
	if err != nil {
		return offset, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// filter log entries
		lines++
		b := scanner.Bytes()
		i := data.Filter(b)

		switch i {
		case data.UMFilter:
			// process UM log type
			matches++
			p := data.NewUMBroadcast()
			if err := data.SendData(p, b); err != nil {
				continue // perhaps log to log file
			}
		case data.P2PFilter:
			// process P2P log type
			matches++
			p := data.NewP2PNumPeers()
			if err := data.SendData(p, b); err != nil {
				continue // perhaps log to log file
			}
		default: //filter.ErrFilter
			misses++
		}
	}

	if err := scanner.Err(); err != nil {
		return offset, err // log.Println(err)
	}

	// fmt.Printf("Process stats - lines: %v, matches: %v, misses: %v, skipped: %v\n", lines, matches, misses, lines-matches-misses)

	// adjust offset for next file read
	offset = size

	return offset, nil
}

func getOffset(f *os.File, offset int64) (int64, int64, error) {
	var size int64

	// get file path stats
	info, err := f.Stat()
	if err != nil {
		return offset, size, err
	}
	size = info.Size()

	// reset offset if file rotated
	if size < offset {
		offset = 0
	}

	return offset, size, nil
}

func GetFilePath(rt string) (string, error) {
	fp := os.Getenv("LOG_FILEPATH")

	if fp == "" {
		u, err := user.Current()
		if err != nil {
			return fp, err
		}
		hd := u.HomeDir

		switch rt { // runtime.GOOS {
		case "darwin":
			fp = fmt.Sprintf("%s/Library/Logs/Theta Edge Node/log.log", hd) // "~/Library/Logs/Theta Edge Node/log.log" // darwin default filepath
		case "windows":
			fp = fmt.Sprintf("%s\\AppData\\Roaming\\Theta Edge Node\\log.log", hd) // "C:\\Users\\<user>\\AppData\\Roaming\\Theta Edge Node\\log.log" // windows default filepath
		// case "linux":
		// 	fp = fmt.Sprintf("%s/path/to/logfile/log.log", hd) // linux default filepath
		default:
			err := errors.New("os not supported")
			return fp, err
		}
	}

	return fp, nil
}

func PokeFilePath(fp string) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
