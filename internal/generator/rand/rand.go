package rand

import (
	"fmt"
	"math/rand"
	"time"
)

// LogEntry value generator
func LogEntry() interface{} {
	var ret interface{}
	x := rand.Int() % 100
	if x < 10 {
		ret = SeededRand.Int()
	} else if x < 20 {
		ret = SeededRand.Float32()
	} else if x < 40 {
		ret = randomURL()
	} else if x > 80 {
		ret = String(SeededRand.Int()%100, Charset+" ")
	} else {
		ret = 0
	}
	return ret
}

// Borrowed from https://www.calhoun.io/creating-random-strings-in-go/

// CharsetLower case
const CharsetLower = "abcdefghijklmnopqrstuvwxyz"

// CharsetUpper case
const CharsetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Charset upper and lower case
const Charset = CharsetLower + CharsetUpper

// SeededRand random number generator
var SeededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[SeededRand.Intn(len(charset))]
	}
	return string(b)
}

// MetricName generator
func MetricName() string {
	return fmt.Sprintf("%s.%s.%s",
		String(SeededRand.Int()%15+5, CharsetLower),
		String(SeededRand.Int()%15+5, CharsetLower),
		String(SeededRand.Int()%15+5, CharsetLower))
}

// String of random characters from characterSet
func String(length int, characterSet string) string {
	return stringWithCharset(length, characterSet)
}

// generate a string resembling a randomized valid URL
func randomURL() string {
	host := fmt.Sprintf("%s.%s.%s",
		String(SeededRand.Int()%30+5, CharsetLower),
		String(SeededRand.Int()%30+5, CharsetLower),
		String(SeededRand.Int()%5+3, CharsetLower))
	return fmt.Sprintf("https://%s:%d/%s/%s",
		host,
		SeededRand.Int()%5000,
		String(SeededRand.Int()%30, Charset),
		String(SeededRand.Int()%30, Charset))
}
