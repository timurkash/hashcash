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
func (h *HashCashData) Stringify() []byte {
	return []byte(fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter))
}

func sha256Hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	bs := h.Sum(nil)
	return []byte(fmt.Sprintf("%x", bs))
}

func (h *HashCashData) IsHashCorrect(hash []byte) bool {
	if h.ZerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:h.ZerosCount] {
		if ch != '0' {
			return false
		}
	}
	return true
}

func (h *HashCashData) GetHash() []byte {
	return sha256Hash(h.Stringify())
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
		ZerosCount: 5,
		Date:       date.Unix(),
		Resource:   clientInfo,
		Rand:       RandStringRunes(10),
		Counter:    0,
	}
	if err := hashCash.ComputeCounter(10_000_000); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v\n", hashCash)
	fmt.Println("=======")
	fmt.Println("after receiving hash and has to be verified")
	hash := hashCash.GetHash()
	fmt.Println(string(hash), hashCash.IsHashCorrect(hash))
}
