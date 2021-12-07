package logs

import (
	"fmt"

	"ar/internal/generator/rand"
)

type account struct {
	id   int
	name string
}

// AccountLogger that decorates logs with account details
type AccountLogger struct {
	accounts []account
}

// newAccountLogger initializes random account data to be added to log messages
func newAccountLogger(num int) *AccountLogger {
	al := &AccountLogger{}
	al.accounts = make([]account, num)
	for x := range al.accounts {
		al.accounts[x] = account{rand.SeededRand.Int() % 5000, rand.String(25, rand.Charset)}
	}
	return al
}

// Dump the account logger object as a string
func (al AccountLogger) Dump() string {
	ret := fmt.Sprintf("Account Logger: %d accounts\n", len(al.accounts))
	for x := range al.accounts {
		ret += fmt.Sprintf("%d: %s\n", al.accounts[x].id, al.accounts[x].name)
	}
	return ret
}

func (al AccountLogger) randomAccount() account {
	return al.accounts[rand.SeededRand.Int()%len(al.accounts)]
}

func (al AccountLogger) Decorator(msg structuredMessage) structuredMessage {
	account := al.randomAccount()
	msg["account"] = make(structuredMessage)
	accountLog := msg["account"].(structuredMessage)
	accountLog["id"] = account.id
	accountLog["name"] = account.name
	return msg
}

// RandomLog message with account information embedded
func (al AccountLogger) RandomLog() structuredMessage {
	return al.Decorator(newLog())
}