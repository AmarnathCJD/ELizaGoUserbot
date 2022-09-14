package modules

import (
	"fmt"
	"runtime"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("+ping", Ping, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("+r", React, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("^\\+alive", Alive, &telegram.Filters{Outgoing: true})
}

func Ping(m *telegram.NewMessage) error {
	_, err := m.Edit(fmt.Sprintf("<b>Pong! %v\nUptime:</b> %s", m.Client.Ping(), time.Since(time.Unix(startTime, 0))))
	return err
}

func React(m *telegram.NewMessage) error {
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

const (
	AlivePicURL  = "https://te.legra.ph/file/cb37180e3aaa92dac6f40.jpg"
	AliveMessage = `<b>Eliza Userbot is alive!</b>

<b>• Master:</b> <a href="tg://user?id=%d">%s</a>
<b>• Mode:</b> <code>userbot</code>
<b>• Uptime:</b> %s
<b>• Golang:</b> %s
<b>• GoGram:</b> %s (<b><i>Layer %d</i></b>)
<b>• Repo:</b> ...`
)

func Alive(m *telegram.NewMessage) error {
	m.Delete()
	_, err := m.RespondMedia(AlivePicURL, telegram.MediaOptions{Caption: fmt.Sprintf(AliveMessage, m.SenderID(), m.Sender.FirstName, time.Since(time.Unix(startTime, 0)).Truncate(time.Second), runtime.Version(), telegram.Version, telegram.ApiVersion)})
	return err
}
