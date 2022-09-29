package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var (
	PMPERMIT_PIC       = getENV("PMPERMIT_PIC", "https://telegra.ph/file/db92ed3d77377856ef911.mp4")
	PM_WARNS           = make(map[int64]int)
	CUSTOM_MIDDLE_TEXT = getENV("CUSTOM_MIDDLE_TEXT", "<b>YOU HAVE TRESPASSED TO MY MASTERS INBOX</b>\nTHIS IS ILLEGAL AND REGARDED AS A CRIME")
	USER_BOT_WARN_ZERO = "<code>You were spamming my sweet master's inbox, henceforth your retarded lame ass has been blocked by my master's userbot⭕️.</code>\n<b>Now GTFO, i'm busy</b>"
	USER_BOT_NO_WARN   = "<code>Hello, This is ELIZA⚠️.You have found your way here to my master,</code>%s's inbox. do your work..\n\n<b>%s</b>\n\n<code>Leave your Name,Reason and 10k$ and hopefully you'll get a reply within 100 light years.</code>⭕️\n\n⭕️<b>Now You Are In Trouble So Send</b> 🔥 <code>/start</code> 🔥 <b>To Start A Valid Conversation!!</b>⭕️"
	APPROVED_IDS       = []int64{}
	PREV_MSG           = make(map[int64]*telegram.NewMessage)

	PM  = "<code>Hello. You are accessing the availabe menu of my Master,</code>%s.\n<i>Let's make this smooth and let me know why you are here.</i>\n</b>Choose one of the following reasons why you are here:<b>\n\n<code>1</code>. To chat with my master\n<code>2</code>. To Give Your Details.\n<code>3</code>. To enquire something\n<code>4</code>. To request something\n"
	ONE = "<i>Okay. Your request has been registered. Do not spam my master's inbox.You can expect a reply within 24 light years. He is a busy guy, unlike you probably.</i>\n\n<b>⚠️ You will be blocked and reported if you spam  ⚠️</b>\n\n<i>Use</i> <code>/start</code> <i>to go back to the main menu.</i>"

	TWO   = "<b>So uncool, this is not your home. Go bother someone else. You have been blocked and reported until further notice.</b>"
	FOUR  = "<i>Okay. My master has not seen your message yet.He usually responds to people,though idk about retarted ones.</i>\n <i>He'll respond when he comes back, if he wants to.There's already a lot of pending messages😶</i>\n <b>Please do not spam unless you wish to be blocked and reported.</b>"
	FIVE  = "<code>Okay. please have the basic manners as to not bother my master too much. If he wishes to help you, he will respond to you soon.</code>\n<b>Do not ask repeatdly else you will be blocked and reported.</b>"
	LWARN = "<b>This is your last warning. DO NOT send another message else you will be blocked and reported. Keep patience. My Master will respond Your Request.</b>\n<i>Use</i> <code>/start</code> <i>to go back to the main menu.<i>"
)

func init() {
	UB.AddMessageHandler(telegram.OnNewMessage, NewPrivateMsg, &telegram.Filters{Incoming: true})
	UB.AddMessageHandler("^\\+approve", ApproveUser, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("^\\+disapprove", DisapproveUser, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("^\\+block", BlockUser, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("^\\+unblock", UnblockUser, &telegram.Filters{Outgoing: true})
}

func NewPrivateMsg(m *telegram.NewMessage) error {
	if !m.IsPrivate() || m.Sender.Bot || m.IsCommand() || m.SenderID() == MasterID {
		return nil
	}
	for _, id := range APPROVED_IDS {
		if id == m.Sender.ID {
			return nil
		}
	}
	if strings.HasPrefix(m.Text(), "/start") {
		m.Reply("started..., now GTFO")
		return nil
	}
	if _, ok := PM_WARNS[m.SenderID()]; !ok {
		PM_WARNS[m.SenderID()] = 1
	} else {
		PM_WARNS[m.SenderID()] += 1
		if PM_WARNS[m.SenderID()] >= 4 {
			r, _ := m.Reply(USER_BOT_WARN_ZERO)
			_, err := m.Client.ContactsBlock(m.Peer)
			if err != nil {
				return err
			}
			PREV_MSG[m.SenderID()] = r
			PM_WARNS[m.SenderID()] = 0
			return nil
		}
	}
	msg, err := m.RespondMedia(PMPERMIT_PIC, telegram.MediaOptions{Caption: fmt.Sprintf(USER_BOT_NO_WARN, aliveName, CUSTOM_MIDDLE_TEXT)})
	if PREV, ok := PREV_MSG[m.SenderID()]; ok {
		PREV.Delete()
	}
	PREV_MSG[m.SenderID()] = msg
	return err
}

func ApproveUser(m *telegram.NewMessage) error {
	SenderID := getUserIDFromEvent(m)
	delete(PM_WARNS, SenderID)
	if _, ok := PREV_MSG[SenderID]; ok {
		PREV_MSG[m.SenderID()].Delete()
		delete(PREV_MSG, SenderID)
	}
	APPROVED_IDS = append(APPROVED_IDS, SenderID)
	msg, err := m.Edit("You have been approved to PM my master.")
	go func() {
		_, _ = m.Client.ContactsUnblock(m.Peer)
		time.Sleep(5 * time.Second)
		msg.Delete()
	}()
	return err
}

func DisapproveUser(m *telegram.NewMessage) error {
	SenderID := getUserIDFromEvent(m)
	delete(PM_WARNS, SenderID)
	if _, ok := PREV_MSG[SenderID]; ok {
		PREV_MSG[m.SenderID()].Delete()
		delete(PREV_MSG, SenderID)
	}
	for i, id := range APPROVED_IDS {
		if id == SenderID {
			APPROVED_IDS = append(APPROVED_IDS[:i], APPROVED_IDS[i+1:]...)
			break
		}
	}
	msg, err := m.Edit("You have been disapproved to PM my master.")
	go func() {
		time.Sleep(5 * time.Second)
		msg.Delete()
	}()
	return err
}

func BlockUser(m *telegram.NewMessage) error {
	SenderID := getUserIDFromEvent(m)
	delete(PM_WARNS, SenderID)
	if _, ok := PREV_MSG[SenderID]; ok {
		PREV_MSG[m.SenderID()].Delete()
		delete(PREV_MSG, SenderID)
	}
	for i, id := range APPROVED_IDS {
		if id == SenderID {
			APPROVED_IDS = append(APPROVED_IDS[:i], APPROVED_IDS[i+1:]...)
			break
		}
	}
	msg, err := m.Edit("Blocked User!")
	go func() {
		_, _ = m.Client.ContactsBlock(m.Peer)
		time.Sleep(5 * time.Second)
		msg.Delete()
	}()
	return err
}

func UnblockUser(m *telegram.NewMessage) error {
	SenderID := getUserIDFromEvent(m)
	delete(PM_WARNS, SenderID)
	if _, ok := PREV_MSG[SenderID]; ok {
		PREV_MSG[m.SenderID()].Delete()
		delete(PREV_MSG, SenderID)
	}
	for i, id := range APPROVED_IDS {
		if id == SenderID {
			APPROVED_IDS = append(APPROVED_IDS[:i], APPROVED_IDS[i+1:]...)
			break
		}
	}
	msg, err := m.Edit("Unblocked User!")
	go func() {
		_, _ = m.Client.ContactsUnblock(m.Peer)
		time.Sleep(5 * time.Second)
		msg.Delete()
	}()
	return err
}
