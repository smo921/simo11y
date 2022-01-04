package logs

import (
	"ar/internal/generator/rand"
	"ar/internal/types"
	"fmt"
)

type service struct {
	name, product, team string
}

// ServiceLogger that decorates logs with service details
type ServiceLogger struct {
	services []service
}

func newServiceLogger(num int) *ServiceLogger {
	sl := &ServiceLogger{}
	sl.services = make([]service, num)
	for i := range sl.services {
		sl.services[i] = service{
			name:    rand.String(32, rand.Charset),
			product: rand.String(32, rand.Charset),
			team:    rand.String(32, rand.Charset),
		}
	}
	return sl
}

// Dump returns the service logger runtime details as a string
func (sl ServiceLogger) Dump() string {
	ret := fmt.Sprintf("Service Logger: %d services\n", len(sl.services))
	for i := range sl.services {
		s := sl.services[i]
		ret += fmt.Sprintf("%d: %s, %s, %s\n", i, s.name, s.product, s.team)
	}
	return ret
}

// Decorator adds service details to a structured message
func (sl ServiceLogger) Decorator(msg types.StructuredMessage) types.StructuredMessage {
	service := sl.randomService()
	msg["service"] = service.name
	return msg
}

func (sl ServiceLogger) randomService() service {
	return sl.services[rand.SeededRand.Int()%len(sl.services)]
}
