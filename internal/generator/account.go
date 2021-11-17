package generator

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

// NewAccountLogger initializes random account data to be added to log messages
func NewAccountLogger(num int) *AccountLogger {
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

// RandomLog message with account information embedded
func (al AccountLogger) RandomLog() map[string]interface{} {
	log := newLog()
	account := al.randomAccount()
	log["account"] = make(map[string]interface{})
	accountLog := log["account"].(map[string]interface{})
	accountLog["id"] = account.id
	accountLog["name"] = account.name
	return log
}
