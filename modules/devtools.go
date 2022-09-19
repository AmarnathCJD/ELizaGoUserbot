package modules

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("\\+json", Jsonify, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+reserved", ReservedChannels, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+sh", Shell, &telegram.Filters{Outgoing: true})
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

func Shell(m *telegram.NewMessage) error {
	var err error
	msg, _ := m.Edit("<code>Processing...</code>")
	if m.Args() == "" {
		_, err := msg.Edit("No code provided!")
		return err
	} else {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		proc := exec.Command("bash", "-c", m.Args())
		proc.Stdout = &stdout
		proc.Stderr = &stderr
		err = proc.Run()
		var result string
		if stdout.String() != string("") {
			result = stdout.String()
		} else if stderr.String() != string("") {
			result = stderr.String()
		} else if err != nil {
			result = err.Error()
		} else {
			result = "No output"
		}
		if len(result) > 4096 {
			defer os.Remove("sh.txt")
			defer m.Delete()
			fs, _ := os.OpenFile("sh.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			defer fs.Close()
			fs.WriteString(result)
			f, _ := m.Client.UploadFile("sh.txt")
			_, err := m.Client.SendMedia(m.ChatID(), f)
			return err
		}
		_, err = msg.Edit(fmt.Sprintf("<b>BASH:</b><code>%s</code>", result))
	}
	return err
}

func GolangCodeEvaluate(code string, variables []interface{}) (string, error) {
	// TODO
	return "", nil
}
