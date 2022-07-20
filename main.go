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
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (h *HashCashData) IsHashCorrect(hash []byte) bool {
	if h.ZerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:h.ZerosCount] {
		if ch != 0 {
			return false
		}
	}
	return true
}

func (h *HashCashData) Stringify() []byte {
	return []byte(fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter))
}

func (h *HashCashData) GetHash() []byte {
	hash := sha256.New()
	hash.Write(h.Stringify())
	return hash.Sum(nil)
}

func (h *HashCashData) ComputeCounter(maxIterations int) error {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		if h.IsHashCorrect(h.GetHash()) {
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
		ZerosCount: 2,
		Date:       date.Unix(),
		Resource:   clientInfo,
		Rand:       RandStringRunes(10),
	}
	if err := hashCash.ComputeCounter(1_000_000); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v\n", hashCash)
	fmt.Println("=======")
	fmt.Println("after receiving hash and has to be verified")
	hash := hashCash.GetHash()
	fmt.Printf("%x %v\n", string(hash), hashCash.IsHashCorrect(hash))
}
