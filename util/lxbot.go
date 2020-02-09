package util

import (
	"github.com/lxbot/lxlib"
	"log"
	"os"
)

type M = map[string]interface{}

func IsErr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func IsReply(msg M) bool {
	return msg["is_reply"] != nil && msg["is_reply"].(bool)
}

func IsRoot(m *lxlib.LXMessage) bool {
	envRoot := os.Getenv("LXBOT_MASTOADMIN_ROOT")
	if envRoot != "" && envRoot == m.User.ID {
		return true
	}

	// TODO: store

	return false
}