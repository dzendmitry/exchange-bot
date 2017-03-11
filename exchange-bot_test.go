package main

import (
	"testing"
	"fmt"
	"net/http"
	"strings"
	"regexp"
	"encoding/json"
	"io/ioutil"
	"io"
)

const TEST_URL = "http://127.0.0.1:8080/event"
const REGEX = `%s[%s:|%s:]+\s%s%s.*`
const REGEX_FLOAT = `[0-9]+\.?[0-9]+`

func TestBehavior(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"биток все, эфир тузэмун", fmt.Sprintf(REGEX, MESSAGE_PREFIX, strings.ToUpper(CURSTR[BTC]), strings.ToUpper(CURSTR[ETH]), REGEX_FLOAT, TO_CUR_SYMBOL)},
		{"usd/euro", fmt.Sprintf(REGEX, MESSAGE_PREFIX, strings.ToUpper(CURSTR[USD]), strings.ToUpper(CURSTR[EUR]), REGEX_FLOAT, TO_CUR_SYMBOL)},
	}
	for _, c := range cases {
		resp, err := http.Post(TEST_URL, "application/json", strings.NewReader(fmt.Sprintf("{\"text\": \"%s\",\"username\": \"tester\",\"display_name\": \"tester\"}", c.in)))
		if err != nil {
			t.Errorf("HTTP POST got %q for %q, want %q", err, c.in, c.want)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("HTTP STATUS CODE got %d want %d", resp.StatusCode, http.StatusCreated)
		}
		if v, ok := resp.Header["Content-Type"]; !ok {
			t.Errorf("HTTP CONTENT TYPE got %d want %d", "nil", "application/json")
		} else if v[0] != "application/json" {
			t.Errorf("HTTP CONTENT TYPE got %d want %d", v[0], "application/json")
		}
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, BODY_BUFFER))
		if err != nil {
			t.Errorf("READING BODY got %q for %q, want %q", err, c.in, c.want)
		}
		if err = resp.Body.Close(); err != nil {
			t.Errorf("CLOSING BODY got %q for %q, want %q", err, c.in, c.want)
		}
		var answer Answer
		if err := json.Unmarshal(body, &answer); err != nil {
			t.Errorf("UNMARSHALLING JSON got %q for %q, want %q", err, c.in, c.want)
		}
		r := regexp.MustCompile(c.want)
		if !r.MatchString(answer.Text) {
			t.Errorf("ANSWER MISMATCH got %q, want %q", answer.Text, c.want)
		}
	}
}

func TestInvalidJson(t *testing.T) {
	resp, err := http.Post(TEST_URL, "application/json", strings.NewReader(fmt.Sprintf("{\"texwt\": \"usd\",\"useame\": \"tester\",\"dislay_name\": \"tester\"}")))
	if err != nil {
		t.Errorf("HTTP POST got %q", err)
	}
	if resp.StatusCode != http.StatusExpectationFailed {
		t.Errorf("WRONG STATUS CODE got %d, want %d", resp.StatusCode, http.StatusExpectationFailed)
	}
}

func TestEmptyMessage(t *testing.T) {
	resp, err := http.Post(TEST_URL, "application/json", strings.NewReader(fmt.Sprintf("{\"text\": \"\",\"username\": \"tester\",\"display_name\": \"tester\"}")))
	if err != nil {
		t.Errorf("HTTP POST got %q", err)
	}
	if resp.StatusCode != http.StatusExpectationFailed {
		t.Errorf("WRONG STATUS CODE got %d, want %d", resp.StatusCode, http.StatusExpectationFailed)
	}
}