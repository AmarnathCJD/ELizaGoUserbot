package modules

import (
	"errors"
	"os"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func getENV(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getUserIDFromEvent(ev *telegram.NewMessage) int64 {
	if ev.IsReply() {
		r, err := ev.GetReplyMessage()
		if err != nil {
			return 0
		}
		return r.SenderID()
	}
	args := ev.Args()
	userID := strings.Split(args, " ")[0]
	if userID == "" {
		return 0
	}
	peerObj, err := ev.Client.GetSendablePeer(userID)
	if err != nil {
		return 0
	}
	peerID := ev.Client.GetPeerID(peerObj)
	return peerID
}

func GetInputPeerUserFromEvent(ev *telegram.NewMessage) (*telegram.InputPeerUser, error) {
	userID := getUserIDFromEvent(ev)
	if userID == 0 {
		return &telegram.InputPeerUser{}, errors.New("Please reply to a user to get their profile picture.")
	}
	peerObj, err := ev.Client.GetSendablePeer(userID)
	if err != nil {
		return &telegram.InputPeerUser{}, err
	}
	user, ok := peerObj.(*telegram.InputPeerUser)
	if !ok {
		return &telegram.InputPeerUser{}, errors.New("Please reply to a user to get their profile picture.")
	}
	return user, nil
}

func GetRight(ev *telegram.NewMessage) *telegram.ChatAdminRights {
	if ev.Chat != nil {
		return ev.Chat.AdminRights
	}
	if ev.Channel != nil {
		return ev.Channel.AdminRights
	}
	return &telegram.ChatAdminRights{}
}
