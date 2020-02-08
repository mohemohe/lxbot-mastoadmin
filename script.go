package main

import (
	"github.com/lxbot/lxlib"
	"log"
	"plugin"
)

type M = map[string]interface{}

var store *lxlib.Store
var ch *chan M

func Boot(s *plugin.Plugin, c *chan M) {
	var err error
	store, err = lxlib.NewStore(s)
	if err != nil {
		log.Fatalln(err)
	}
	ch = c
}

func OnMessage() []func(M) M {
	return []func(M) M{
		ping,
	}
}

func ping(msg M) M {
	m, err := lxlib.NewLXMessage(msg)
	if isErr(err) || !isReply(msg) {
		return nil
	}

	if m.Message.Text != "ping" {
		return nil
	}

	r, err := m.SetText("pong").Reply().ToMap()
	if isErr(err) {
		return nil
	}
	return r
}

func isErr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func isReply(msg M) bool {
	return msg["is_reply"] != nil && msg["is_reply"].(bool)
}