package utils

import (
	"fmt"
	"log"

	"github.com/cskr/pubsub"
)

func NewStdLogger(prefix string, flags int) *log.Logger {

	l := stdioLogger{}
	logger := log.New(&l, prefix, flags)

	return logger
}

type stdioLogger struct {
}

func (l *stdioLogger) Write(input []byte) (n int, err error) {

	fmt.Printf("%v", string(input))
	return len(input), nil
}

func NewChLogger(ps *pubsub.PubSub, topic string, prefix string) *log.Logger {

	cl := chanLogger{}
	cl.ps = ps
	cl.topic = topic
	logger := log.New(&cl, prefix, log.Ltime)
	return logger
}

type chanLogger struct {
	ps    *pubsub.PubSub
	topic string
}

func (c *chanLogger) Write(input []byte) (n int, err error) {
	c.ps.Pub(string(input), c.topic)
	return len(input), nil
}

type nullLogger struct {
}

func NewNullLogger() *log.Logger {
	l := nullLogger{}
	logger := log.New(&l, "", 0)

	return logger
}

func (l *nullLogger) Write(input []byte) (n int, err error) {

	return len(input), nil
}
