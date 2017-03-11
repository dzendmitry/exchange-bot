package main

import "strings"

const CURRENCIES_COUNT = 4

const USD = 0
const EUR = 1
const BTC = 2
const ETH = 3

var WORDS = map[string]int {
	"usd":      USD,
	"бакс":     USD,
	"доллар":   USD,
	"евро":     EUR,
	"еуро":     EUR,
	"eur":      EUR,
	"btc":      BTC,
	"биткоин":  BTC,
	"биток":    BTC,
	"ethereum": ETH,
	"эфир":     ETH,
}

type Parser struct {
	CurIndexes []bool
}

func (p *Parser) Parse(message string) {
	words := strings.Fields(message)
	p.CurIndexes = make([]bool, CURRENCIES_COUNT)
	for _, word := range words {
		word = strings.ToLower(word)
		for ws, wv := range WORDS {
			if strings.Contains(word, ws) {
				p.CurIndexes[wv] = true
			}
		}
		if p.IsFull() {
			break
		}
	}
}

func (p *Parser) IsCur() bool {
	for _, v := range p.CurIndexes {
		if v {
			return true
		}
	}
	return false
}

func (p *Parser) IsFull() bool {
	for _, v := range p.CurIndexes {
		if !v {
			return false
		}
	}
	return true
}