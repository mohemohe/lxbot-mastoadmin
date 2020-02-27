package rails

import (
	"github.com/lxbot/lxlib"
	"github.com/mohemohe/lxbot-mastoadmin/tootctl"
	"github.com/mohemohe/lxbot-mastoadmin/util"
)

func Help() string {
	return `rails_db_migrate: できるとでも思ったのか？
`
}

func RailsDbMigrate(msg util.M, ch *chan util.M) util.M {
	m, err := lxlib.NewLXMessage(msg)
	if util.IsErr(err) || !util.IsReply(msg) || !util.IsRoot(m) {
		return nil
	}
	if m.Message.Text != "rails_db_migrate" {
		return nil
	}

	script := "rails db:migrate"
	go tootctl.Run(msg, script, ch)
	return nil
}
