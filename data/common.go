package data

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	isoDatetimeRE = `(?m)\d+-\d+-\d+\s?\w?\d+:\d+:\d+\.?\d*?\w+`
	timeStrFormat = `2006-01-02T15:04:05.999`
)

var (
	apiAddr = "http://127.0.0.1:8000"
	apiKey  = "devkey"
)

type Parser interface {
	ToJSON() ([]byte, error)
	Parse([]byte) error
}

func SendData(p Parser, b []byte) error {
	// parse log
	if err := p.Parse(b); err != nil {
		return err
	}

	// create json
	d, err := p.ToJSON()
	if err != nil {
		return err
	}

	// fuzzing request
	fuzzRequest()

	// network request
	url := getServiceURI(p)
	j := strings.NewReader(string(d))
	req, err := http.NewRequest(http.MethodPost, url, j)
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// get request status code
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("~%s %s %v\n", resp.Request.URL.Path, resp.Request.Method, resp.StatusCode)

	return nil
}

func parseTime(b []byte) (time.Time, error) {
	re := regexp.MustCompile(isoDatetimeRE)
	re.Longest()

	m := re.Find(b)
	s := strings.Replace(string(m), " ", "T", -1)

	// get relevant locations
	loc := time.Now().Location()       // loc, err := time.LoadLocation("Local")
	utc := time.Now().UTC().Location() // utc, err := time.LoadLocation("UTC")

	// parse time at location
	t, err := time.ParseInLocation(timeStrFormat, s, loc)
	if err != nil {
		return t, err
	}

	// return time at utc
	return t.In(utc), nil
}

func getServiceURI(p Parser) string {
	var url string

	switch p.(type) {
	case *UMBroadcast:
		url = umBroadcastedServiceURL
	case *P2PNumPeers:
		url = p2pNumPeersServiceURL
	}

	return url
}

func fuzzRequest() {
	// rethink how seed is randomized
	rand.Seed(time.Now().UnixNano())
	d := time.Duration((rand.Intn(6000)))
	time.Sleep(d * time.Millisecond)
}
