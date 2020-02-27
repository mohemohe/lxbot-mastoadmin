package tootctl

import (
	"github.com/lxbot/lxlib"
	"github.com/mohemohe/lxbot-mastoadmin/util"
	"strings"
)

func Help() string {
	return `tootctl: tootctl
`
}

func Tootctl(msg util.M, ch *chan util.M) util.M {
	m, err := lxlib.NewLXMessage(msg)
	if util.IsErr(err) || !util.IsReply(msg) || !util.IsRoot(m) {
		return nil
	}
	if !util.Prefix(m.Message.Text, "tootctl") {
		return nil
	}

	script := strings.Join(strings.Split(m.Message.Text, "\n"), " ")
	go Run(msg, script, ch)
	return nil
}
