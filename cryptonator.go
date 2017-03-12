package main

import (
	"time"
	"net/http"
	"errors"
	"fmt"
	"io"
	"encoding/json"
	"io/ioutil"
	"strings"
)

const ALL_REQUESTS_TIMEOUT = 5

const CRYPTONATOR_URL = "https://api.cryptonator.com/api/ticker/%s-%s"
const TO_CUR = "rur"
const TO_CUR_SYMBOL = "₽"
const MESSAGE_PREFIX = "Exchange courses are "

var CURSTR = []string{"usd", "eur", "btc", "eth"}

type CryptonatorTicker struct {
	Base string                `json:"base"`
	Target string              `json:"target"`
	Price string               `json:"price"`
	Volume string              `json:"volume"`
	Change string              `json:"change"`
}

type CryptonatorAnswer struct {
	Ticker CryptonatorTicker   `json:"ticker"`
	Timestamp int64            `json:"timestamp"`
	Success bool               `json:"success"`
	Error string               `json:"error"`
}

type CryptonatorAnswerStatus struct {
	C *CryptonatorAnswer
	Err error
}

type Crypronator struct {
	parser *Parser
}

func (c *Crypronator) Download() (string, error) {
	if !c.parser.IsCur() {
		return "", errors.New("Nothing to download")
	}
	message := MESSAGE_PREFIX
	resChan := make(chan *CryptonatorAnswerStatus)
	var ctr int
	for i, v := range c.parser.CurIndexes {
		if !v {
			continue
		}
		go Fetch(CURSTR[i], resChan)
		ctr += 1
	}
	var err error
	timer := time.NewTimer(ALL_REQUESTS_TIMEOUT * time.Second)
	L: for i := 0; i < ctr; i++ {
		select {
		case r := <- resChan:
			if r.Err != nil || !r.C.Success || len(r.C.Error) > 0 {
				continue
			}
			ts := time.Unix(r.C.Timestamp, 0)
			if ts.Before(time.Now().AddDate(0, 0, -1)) {
				continue
			}
			message += strings.ToUpper(r.C.Ticker.Base) + ": " + r.C.Ticker.Price + TO_CUR_SYMBOL + " "
		case <- timer.C:
			break L
		}
	}
	if message == MESSAGE_PREFIX {
		err = errors.New("No data")
	}
	return message, err
}

func Fetch(fromCur string, resChan chan *CryptonatorAnswerStatus) {
	res, err := http.Get(fmt.Sprintf(CRYPTONATOR_URL, fromCur, TO_CUR))
	if err != nil {
		resChan <- &CryptonatorAnswerStatus{nil, err}
		return
	}
	if res.StatusCode != http.StatusOK {
		resChan <- &CryptonatorAnswerStatus{
			C: nil,
			Err: errors.New(fmt.Sprintf("Downloading result for %s is %s", fromCur, res.Status)),
		}
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, BODY_BUFFER))
	if err != nil {
		resChan <- &CryptonatorAnswerStatus{
			C: nil,
			Err: errors.New(fmt.Sprintf("Reading body %s: %s", fromCur, res.Status)),
		}
		return
	}
	if err = res.Body.Close(); err != nil {
		resChan <- &CryptonatorAnswerStatus{
			C: nil,
			Err: errors.New(fmt.Sprintf("Closing body %s: %s", fromCur, res.Status)),
		}
		return
	}

	var ca CryptonatorAnswer
	if err := json.Unmarshal(body, &ca); err != nil {
		resChan <- &CryptonatorAnswerStatus{
			C: nil,
			Err: errors.New(fmt.Sprintf("Unmarshalling cryptonator answer %s: %s", fromCur, res.Status)),
		}
		return
	}

	resChan <- &CryptonatorAnswerStatus{&ca, nil}
}