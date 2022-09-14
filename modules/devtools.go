package modules

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("\\+json", Jsonify, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+reserved", ReservedChannels, &telegram.Filters{Outgoing: true})
}

func Jsonify(m *telegram.NewMessage) error {
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			return err
		}
		_, err = m.Respond(fmt.Sprintf("<code>%s</code>", r.Marshal()))
		return err
	}
	_, err := m.Respond(fmt.Sprintf("<code>%s</code>", m.Marshal()))
	return err
}

func ReservedChannels(m *telegram.NewMessage) error {
	m.Edit("<code>Getting reserved channels...</code>")
	c, err := m.Client.ChannelsGetAdminedPublicChannels(false, false)
	if err != nil {
		_, err = m.Respond(fmt.Sprintf("Error: %v", err))
		return err
	}
	var s = "<b>Reserved Channels:</b>\n"
	switch c := c.(type) {
	case *telegram.MessagesChatsObj:
		for _, v := range c.Chats {
			switch v := v.(type) {
			case *telegram.Channel:
				s += fmt.Sprintf("%s (%d) - @%s\n", v.Title, v.ID, v.Username)
			case *telegram.ChatObj:
				s += fmt.Sprintf("%s (%d)\n", v.Title, v.ID)
			case *telegram.ChatForbidden:
				s += fmt.Sprintf("%s (%d)\n", v.Title, v.ID)
			}
		}
	case *telegram.MessagesChatsSlice:
		for _, v := range c.Chats {
			switch v := v.(type) {
			case *telegram.Channel:
				s += fmt.Sprintf("%s (%d) - %s\n", v.Title, v.ID, v.Username)
			case *telegram.ChatObj:
				s += fmt.Sprintf("%s (%d)\n", v.Title, v.ID)
			case *telegram.ChatForbidden:
				s += fmt.Sprintf("%s (%d)\n", v.Title, v.ID)
			}
		}
	}
	_, err = m.Edit(s)
	return err
}
