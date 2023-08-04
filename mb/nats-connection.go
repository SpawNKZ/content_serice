package mb

import (
	"github.com/go-kit/log"
	"github.com/nats-io/nats.go"
)

func NewNatsConnection(connString string, logger log.Logger) (*nats.Conn, func(), error) {
	nc, err := nats.Connect(connString)
	if err != nil {
		logger.Log("cannot connect to NATS", err)
		return nil, nil, err
	}

	return nc, nc.Close, nil
}
