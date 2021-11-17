package main

import (
	"fmt"
)

type account struct {
	id   int
	name string
}

type accountLogger struct {
	accounts []account
}

func newAccountLogger(num int) *accountLogger {
	al := &accountLogger{}
	al.accounts = make([]account, num)
	for x := range al.accounts {
		al.accounts[x] = account{seededRand.Int() % 5000, randomString(25, charset)}
	}
	return al
}

func (al accountLogger) dump() string {
	ret := fmt.Sprintf("Account Logger: %d accounts\n", len(al.accounts))
	for x := range al.accounts {
		ret += fmt.Sprintf("%d: %s\n", al.accounts[x].id, al.accounts[x].name)
	}
	return ret
}

func (al accountLogger) randomAccount() account {
	return al.accounts[seededRand.Int()%len(al.accounts)]
}

func (al accountLogger) randomLog() map[string]interface{} {
	log := newMessage()
	account := al.randomAccount()
	log["account"] = make(map[string]interface{})
	accountLog := log["account"].(map[string]interface{})
	accountLog["id"] = account.id
	accountLog["name"] = account.name
	return log
}
