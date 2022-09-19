package modules

import (
	"fmt"
	"runtime"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("+ping", Ping, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("+r(?: |$)(.*)", React, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("^\\+alive", Alive, &telegram.Filters{Outgoing: true})
	Bot.AddInlineHandler("alive", InlineAliveMessage)
}

func Ping(m *telegram.NewMessage) error {
	_, err := m.Edit(fmt.Sprintf("<b>Pong! %v\nUptime:</b> %s", m.Client.Ping(), time.Since(time.Unix(startTime, 0))))
	return err
}

func React(m *telegram.NewMessage) error {
	defer m.Delete()
	emo := "👍"
	if m.Args() != "" {
		emo = m.Args()
	}
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			return err
		}
		e := r.React(emo)
		if e != nil {
			m.Edit(fmt.Sprintf("Error: %v", e))
		}
		return e
	}
	e := m.React(emo)
	if e != nil {
		m.Edit(fmt.Sprintf("Error: %v", e))
	}
	return e
}

var (
	AlivePicURL = getENV("ALIVE_PIC", "https://te.legra.ph/file/cb37180e3aaa92dac6f40.jpg")
)

const (
	AliveMessage = `<b>Eliza Userbot is alive!</b>

<b>• Master:</b> <a href="tg://user?id=%d">%s</a>
<b>• Mode:</b> <b>dual</b>
<b>• Uptime:</b> %s
<b>• Golang:</b> %s
<b>• GoGram:</b> %s (<b><i>%d</i></b>)`
)

func Alive(m *telegram.NewMessage) error {
	m.Delete()
	q, e := m.Client.InlineQuery(BotUName, &telegram.InlineOptions{Dialog: m.ChatID(), Query: "alive"})
	peer, _ := m.Client.GetSendablePeer(m.ChatID())
	if len(q.Results) > 0 {
		_, e = m.Client.MessagesSendInlineBotResult(&telegram.MessagesSendInlineBotResultParams{
			Peer:     peer,
			RandomID: telegram.GenRandInt(),
			QueryID:  q.QueryID,
			ID:       q.Results[0].(*telegram.BotInlineMediaResult).ID,
		})
	}
	return e
}

var (
	PHOTO *telegram.InputMediaPhoto
)

func InlineAliveMessage(m *telegram.InlineQuery) error {
	messageText := fmt.Sprintf(AliveMessage, m.SenderID, m.Sender.FirstName, time.Since(time.Unix(startTime, 0)).Truncate(time.Second), runtime.Version(), telegram.Version, telegram.ApiVersion)
	var b = &telegram.Button{}
	var r = m.Builder()

	r.Photo(PHOTO, &telegram.ArticleOptions{
		Caption:     messageText,
		ReplyMarkup: b.Keyboard(b.Row(b.URL("Repo", "https://github.com/Amarnathcjd/elizagouserbot"))),
	})
	_, err := m.Answer(r.Results())
	return err
}

func init() {
	pic := telegram.InputMediaPhotoExternal{
		URL: AlivePicURL,
	}
	PIC, _ := Bot.MessagesUploadMedia(&telegram.InputPeerSelf{}, &pic)
	IMAGE := *PIC.(*telegram.MessageMediaPhoto).Photo.(*telegram.PhotoObj)
	PHOTO = &telegram.InputMediaPhoto{
		ID: &telegram.InputPhotoObj{
			ID:            IMAGE.ID,
			AccessHash:    IMAGE.AccessHash,
			FileReference: IMAGE.FileReference,
		},
	}
}
