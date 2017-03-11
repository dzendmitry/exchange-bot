package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"log"
)

const BOT_NAME = "exchange-bot"

const BODY_BUFFER = 1 * 1024 * 1024

type Event struct {
	Text string         `json:"text"`
	Username string     `json:"username"`
	DisplayName string  `json:"display_name"`
}

type Answer struct {
	Text string         `json:"text"`
	Bot  string         `json:"bot"`
}

func eventHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, BODY_BUFFER))
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	if err = req.Body.Close(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	var parser Parser
	parser.Parse(event.Text)
	cryptonator := Crypronator{&parser}
	message, err := cryptonator.Download()
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	answer := Answer{
		Text: message,
		Bot: BOT_NAME,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
}

func infoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	info := map[string]string{
		"author": "Dmitry Zenin aka dzendmitry",
		"info":   "Bot reacts to words from list and shows money exchange rates",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
	}
}

func main() {
	http.HandleFunc("/event", eventHandler)
	http.HandleFunc("/info", infoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
