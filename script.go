package main

import (
	"github.com/lxbot/lxlib"
	"github.com/mohemohe/lxbot-mastoadmin/aws"
	"github.com/mohemohe/lxbot-mastoadmin/tootctl"
	"github.com/mohemohe/lxbot-mastoadmin/util"
	"log"
	"plugin"
)

var store *lxlib.Store
var ch *chan util.M

func Boot(s *plugin.Plugin, c *chan util.M) {
	var err error
	store, err = lxlib.NewStore(s)
	if err != nil {
		log.Fatalln(err)
	}
	ch = c
}

func Help() string {
	t := `ping: pong
check: アクセス権限チェック
` + aws.Help() + tootctl.Help()

	return t
}

func OnMessage() []func(util.M) util.M {
	return []func(util.M) util.M{
		ping,
		check,
		aws.EcsServiceLs,
		aws.EcsServiceStatus,
		aws.EcsServiceScale,
		func(m util.M) util.M {
			return tootctl.Tootctl(m, ch)
		},
	}
}

func ping(msg util.M) util.M {
	m, err := lxlib.NewLXMessage(msg)
	if util.IsErr(err) || !util.IsReply(msg) {
		return nil
	}

	if m.Message.Text != "ping" {
		return nil
	}

	r, err := m.SetText("pong").Reply().ToMap()
	if util.IsErr(err) {
		return nil
	}
	return r
}

func check(msg util.M) util.M {
	m, err := lxlib.NewLXMessage(msg)
	if util.IsErr(err) || !util.IsReply(msg) {
		return nil
	}

	if m.Message.Text != "check" {
		return nil
	}

	if util.IsRoot(m) {
		m.SetText("キミの権限は root だよ")
	} else {
		m.SetText("キミに操作権限は無いよ")
	}

	r, err := m.Reply().ToMap()
	if util.IsErr(err) {
		return nil
	}
	return r
}
