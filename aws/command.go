package aws

import (
	"github.com/lxbot/lxlib"
	"github.com/mohemohe/lxbot-mastoadmin/util"
	"log"
	"strconv"
	"strings"
)

func Help() string {
	return `aws ecs service ls: ECSで動かしてるサービスのリストを表示
aws ecs service status [service_name]: ECSで動かしてるサービスの稼働状況を表示
aws ecs service scale [service_name] [count]: ECSで動かしてるサービスのタスク数を設定
`
}

func EcsServiceLs(msg util.M) util.M {
	m, err := lxlib.NewLXMessage(msg)
	if util.IsErr(err) || !util.IsReply(msg) || !util.IsRoot(m) {
		return nil
	}

	if !util.Equals(m.Message.Text, "aws", "ecs", "service", "ls") {
		return nil
	}

	services := ListServices()
	r, err := m.SetText("登録してるサービスだよ\n" + strings.Join(services, "\n")).Reply().ToMap()
	if util.IsErr(err) {
		return nil
	}
	return r
}

func EcsServiceStatus(msg util.M) util.M {
	m, err := lxlib.NewLXMessage(msg)
	if util.IsErr(err) || !util.IsReply(msg) || !util.IsRoot(m) {
		return nil
	}
	if !util.Prefix(m.Message.Text, "aws", "ecs", "service", "status") {
		return nil
	}
	args := strings.Fields(m.Message.Text)
	if len(args) != 5 {
		return nil
	}

	name := args[4]
	status := GetServiceStatus(name)
	if status == nil {
		return nil
	}
	text := status.Name + "\n" + "STATUS: " + status.Status + "\n" + "DESIRED: " + strconv.Itoa(status.DesiredCount) + "\n" + "RUNNING: " + strconv.Itoa(status.RunningCount) + "\n" + "PENDING: " + strconv.Itoa(status.PendingCount)
	r, err := m.SetText(text).Reply().ToMap()
	if util.IsErr(err) {
		return nil
	}
	return r
}

func EcsServiceScale(msg util.M) util.M {
	m, err := lxlib.NewLXMessage(msg)
	if util.IsErr(err) || !util.IsReply(msg) || !util.IsRoot(m) {
		return nil
	}

	if !util.Prefix(m.Message.Text, "aws", "ecs", "service", "scale") {
		return nil
	}

	args := strings.Fields(m.Message.Text)
	if len(args) != 6 {
		return nil
	}

	text := ""

	name := args[4]
	countStr := args[5]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Println(err)
		text = "タスク数がおかしいよ"
	}
	if count == 0 {
		text = name + " のタスク数を 0 にするなんてとんでもない"
	} else if count < 0 {
		text = "どーやって " + name + " のタスク数をマイナスにするの"
	} else if ok := ScaleService(name, count); !ok {
		text = name + " のタスク数の変更に失敗しちゃった"
	} else {
		text = name + " のタスク数を "+ countStr + " にしたよ"
	}
	r, err := m.SetText(text).Reply().ToMap()
	if util.IsErr(err) {
		return nil
	}
	return r
}
