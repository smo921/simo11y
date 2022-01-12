package logs

import (
	"fmt"

	"simo11y/internal/generator/rand"
	"simo11y/internal/types"
)

type account struct {
	id            int
	name          string
	logsGenerated int
}

// AccountLogger that decorates logs with account details
type AccountLogger struct {
	accounts []account
}

// newAccountLogger initializes random account data to be Mutateed to log messages
func newAccountLogger(num int) *AccountLogger {
	al := &AccountLogger{}
	al.accounts = make([]account, num)
	for x := range al.accounts {
		al.accounts[x] = account{
			id:            rand.SeededRand.Int() % 5000,
			name:          rand.String(25, rand.Charset),
			logsGenerated: 0,
		}
	}
	return al
}

// Dump the account logger object as a string
func (al AccountLogger) String() string {
	ret := fmt.Sprintf("Account Logger: %d accounts\n", len(al.accounts))
	for x := range al.accounts {
		ret += fmt.Sprintf("%d: %s\n", al.accounts[x].id, al.accounts[x].name)
	}
	return ret
}

func (al AccountLogger) randomAccount() account {
	return al.accounts[rand.SeededRand.Int()%len(al.accounts)]
}

// Decorator Mutates account details to a structured message
func (al AccountLogger) Decorator(msg types.StructuredMessage) types.StructuredMessage {
	account := al.randomAccount()
	msg["account"] = make(types.StructuredMessage)
	accountLog := msg["account"].(types.StructuredMessage)
	accountLog["id"] = account.id
	accountLog["name"] = account.name
	account.logsGenerated++
	return msg
}

// RandomLog message with account information embedded
func (al AccountLogger) RandomLog() types.StructuredMessage {
	return al.Decorator(newLog())
}
