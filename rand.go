package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomLogEntry() interface{} {
	var ret interface{}
	x := rand.Int() % 100
	if x < 10 {
		ret = rand.Int()
	} else if x < 20 {
		ret = rand.Float32()
	} else if x < 40 {
		ret = randomURL()
	} else if x > 80 {
		ret = randomString(rand.Int()%100, charset+" ")
	} else {
		ret = 0
	}
	return ret
}

// Borrowed from https://www.calhoun.io/creating-random-strings-in-go/
const charsetLower = "abcdefghijklmnopqrstuvwxyz"
const charsetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var charset = charsetLower + charsetUpper

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int, characterSet string) string {
	return stringWithCharset(length, characterSet)
}

func randomURL() string {
	host := fmt.Sprintf("%s.%s.%s",
		randomString(seededRand.Int()%30+5, charsetLower),
		randomString(seededRand.Int()%30+5, charsetLower),
		randomString(seededRand.Int()%5+3, charsetLower))
	return fmt.Sprintf("https://%s:%d/%s/%s",
		host,
		seededRand.Int()%5000,
		randomString(seededRand.Int()%30, charset),
		randomString(seededRand.Int()%30, charset))
}
