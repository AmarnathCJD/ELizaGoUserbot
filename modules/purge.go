package modules

import (
	"fmt"
	"strconv"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func Purge(m *telegram.NewMessage) error {
	if !GetRight(m).CanDeleteMessages() {
		m.Reply("You don't have the rights to delete messages.")
		return nil
	}
	MessageCount := m.Args()
	IDS := []int32{}
	if MessageCount == "" && !m.IsReply() {
		m.Reply("Please specify a number of messages to delete.")
		return nil
	}
	if MessageCount == "" {
		ReplyID := m.ReplyToMsgID()
		CurrentID := m.ID - 1
		for CurrentID >= ReplyID {
			IDS = append(IDS, CurrentID)
			CurrentID--
		}
	} else {
		MsgCount, err := strconv.Atoi(MessageCount)
		if err != nil || MsgCount < 1 {
			m.Reply("Please specify a valid number of messages to delete.")
			return nil
		}
		CurrentID := m.ID - 1
		for i := 0; i < MsgCount; i++ {
			IDS = append(IDS, CurrentID)
			CurrentID--
		}
	}
	SubIDs := [][]int32{}
	for i := 0; i < len(IDS); i += 200 {
		end := i + 200
		if end > len(IDS) {
			end = len(IDS)
		}
		SubIDs = append(SubIDs, IDS[i:end])
	}
	for _, sub := range SubIDs {
		err := m.Client.DeleteMessage(m.ChatID(), sub...)
		if err != nil {
			m.Reply(fmt.Sprintf("Error: %s", err))
			return nil
		}
	}
	msg, _ := m.Respond("Purge complete.")
	time.Sleep(5 * time.Second)
	m.Client.DeleteMessage(m.ChatID(), msg.ID)
	return nil
}

func Del(m *telegram.NewMessage) error {
	if !GetRight(m).CanDeleteMessages() {
		m.Reply("You don't have the rights to delete messages.")
		return nil
	}
	if !m.IsReply() {
		m.Reply("Please reply to a message to delete.")
		return nil
	}
	Delerr := m.Client.DeleteMessage(m.ChatID(), m.ReplyToMsgID())
	if Delerr != nil {
		m.Respond(fmt.Sprintf("Error: %s", Delerr))
		return nil
	}
	return nil
}
