package modules

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("\\+stoi", StickerToImage, &telegram.Filters{Outgoing: true})
}

func StickerToImage(m *telegram.NewMessage) error {
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			return err
		}
		if !r.IsMedia() {
			m.Reply("Reply to a sticker!")
			return nil
		}
		if r.MediaType() != "document" {
			m.Reply(fmt.Sprintf("Reply to a sticker! Got %s", r.MediaType()))
			return nil
		}
		f, err := r.Download(&telegram.DownloadOptions{FileName: "sticker.jpg"})
		if err != nil {
			return err
		}
		_, err = m.RespondMedia(f, telegram.MediaOptions{Caption: "Here's your image!"})
		return err
	}
	m.Reply("Reply to a sticker!")
	return nil
}
