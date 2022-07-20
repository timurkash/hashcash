package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type HashCashData struct {
	Version    int
	ZerosCount int
	Date       int64
	Resource   string
	Rand       string
	Counter    int
	hash       []byte
}

const (
	complexity   = 3
	maxIteration = 100_000_000
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (h *HashCashData) IsHashCorrect() bool {
	if h.ZerosCount > len(h.hash) {
		return false
	}
	for _, ch := range h.hash[:h.ZerosCount] {
		if ch != 0 {
			return false
		}
	}
	return true
}

func (h *HashCashData) Stringify() []byte {
	return []byte(fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter))
}

func (h *HashCashData) PrintHash() {
	fmt.Printf("%x\n", h.hash)
}

func (h *HashCashData) CalcHash() {
	hash := sha256.New()
	hash.Write(h.Stringify())
	h.hash = hash.Sum(nil)
}

func (h *HashCashData) ComputeCounter(maxIterations int) error {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		h.CalcHash()
		if h.IsHashCorrect() {
			return nil
		}
		h.Counter++
	}
	return fmt.Errorf("max iterations exceeded")
}

func main() {
	defer func() func() {
		start := time.Now()
		return func() {
			fmt.Println(time.Since(start))
		}
	}()()
	date := time.Now()
	clientInfo := "hashCash"
	rand.Seed(time.Now().UnixNano())
	hashCash := &HashCashData{
		Version:    1,
		ZerosCount: complexity,
		Date:       date.Unix(),
		Resource:   clientInfo,
		Rand:       RandStringRunes(10),
	}
	if err := hashCash.ComputeCounter(maxIteration); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", hashCash)
	fmt.Println("=======")
	fmt.Println("after receiving hash and has to be verified")
	hashCash.CalcHash()
	hashCash.PrintHash()
	fmt.Println(hashCash.IsHashCorrect())
}
