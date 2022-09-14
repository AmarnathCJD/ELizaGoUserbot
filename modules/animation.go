package modules

import (
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("\\+angry", AngryMessage, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+stupid", StupidMessage, &telegram.Filters{Outgoing: true})
}

var (
	angryChars = [11]string{
		"😡😡😡",
		"I am angry with you",
		"Just shut up",
		"And RUN Away NOW",
		"Or else",
		"I would call CEO of Telegram",
		"He is my friend warning you",
		"My friend is also a hacker...",
		"I would call him if you don't shup up",
		"🤬🤬Warning you, Don't repeat it again and shup up now...🤬🤬",
		"🤬🤬🤬🤬🤬",
	}
	StupidChars = [14]string{
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n🧠         <(^_^ <)🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n🧠       <(^_^ <)  🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n🧠     <(^_^ <)    🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n🧠   <(^_^ <)      🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n🧠 <(^_^ <)        🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n🧠<(^_^ <)         🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n(> ^_^)>🧠         🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n  (> ^_^)>🧠       🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n    (> ^_^)>🧠     🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n      (> ^_^)>🧠   🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n        (> ^_^)>🧠 🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n          (> ^_^)>🧠🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n           (> ^_^)>🗑",
		"YOᑌᖇ ᗷᖇᗩIᑎ ➡️ 🧠\n\n           <(^_^ <)🗑",
	}
)

func AngryMessage(m *telegram.NewMessage) error {
	editSleep := 2 * time.Second
	for i := 0; i < len(angryChars); i++ {
		m.Edit(angryChars[i], telegram.SendOptions{ParseMode: "-"})
		time.Sleep(editSleep)
	}
	return nil
}

func StupidMessage(m *telegram.NewMessage) error {
	editSleep := 1 * time.Second
	for i := 0; i < len(StupidChars); i++ {
		m.Edit(StupidChars[i], telegram.SendOptions{ParseMode: "-"})
		time.Sleep(editSleep)
	}
	return nil
}
