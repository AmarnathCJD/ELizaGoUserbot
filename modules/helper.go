package modules

import "github.com/amarnathcjd/gogram/telegram"

func init() {
	Bot.AddInlineHandler("help", HelpBot)
	UB.AddMessageHandler("^\\+help", Help, &telegram.Filters{Outgoing: true})
}

var (
	Modules = []string{
		"Start",
		"Animation",
		"Admin", "Devtools",
		"StickerTools", "PMPermit", "ProfilePic", "Misc", "FileTools", "Afk",
	}

	b = telegram.Button{}
)

func genHelpKeyboard() *telegram.ReplyInlineMarkup {
	var rows []*telegram.KeyboardButtonRow
	var buttons []telegram.KeyboardButton
	for _, module := range Modules {
		buttons = append(buttons, (b.Data("📚 "+module, "help_"+module)))
		if len(buttons) == 2 {
			rows = append(rows, b.Row(buttons...))
			buttons = nil
		}
	}
	if len(buttons) > 0 {
		rows = append(rows, b.Row(buttons...))
	}
	return b.Keyboard(rows...)
}

const (
	helpCaption = `📚 <b>Help</b>
Here are the available modules. Click on the buttons below to get help for that module.`
)

func HelpBot(m *telegram.InlineQuery) error {
	var b = m.Builder()
	b.Article("help", helpCaption, helpCaption, &telegram.ArticleOptions{ReplyMarkup: genHelpKeyboard()})
	_, err := m.Answer(b.Results())
	return err
}

func Help(m *telegram.NewMessage) error {
	m.Delete()
	q, e := m.Client.InlineQuery(BotUName, &telegram.InlineOptions{Dialog: m.ChatID(), Query: "help"})
	peer, _ := m.Client.GetSendablePeer(m.ChatID())
	if len(q.Results) > 0 {
		_, e = m.Client.MessagesSendInlineBotResult(&telegram.MessagesSendInlineBotResultParams{
			Peer:     peer,
			RandomID: telegram.GenRandInt(),
			QueryID:  q.QueryID,
			ID:       q.Results[0].(*telegram.BotInlineResultObj).ID,
		})
	}
	return e
}
