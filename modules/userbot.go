package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/joho/godotenv"
)

var (
	UB, Bot    = UserbotClient()
	startTime  int64
	WorkDir, _ = os.Getwd()
	BotUName   string
	aliveName  string
	MasterID   int64
)

func UserbotClient() (*telegram.Client, *telegram.Client) {
	godotenv.Load()
	STRING_SESSION := os.Getenv("STRING_SESSION")
	fmt.Println("Starting Userbot...")
	startTime = time.Now().Unix()
	API_ID := os.Getenv("API_ID")
	API_HASH := os.Getenv("API_HASH")
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	b, err := telegram.TelegramClient(telegram.ClientConfig{
		AppID:         parseINT(API_ID),
		AppHash:       API_HASH,
		StringSession: STRING_SESSION,
		ParseMode:     "HTML",
	})
	if err != nil {
		panic(err)
	}
	bot, _ := telegram.TelegramClient(telegram.ClientConfig{
		AppID:       parseINT(API_ID),
		AppHash:     API_HASH,
		ParseMode:   "HTML",
		SessionFile: filepath.Join(WorkDir, "bot.session"),
	})
	if err := bot.LoginBot(BOT_TOKEN); err != nil {
		panic(err)
	}
	me, err := b.GetMe()
	MasterID = me.ID
	bt, _ := bot.GetMe()
	BotUName = bt.Username
	aliveName = me.FirstName + " " + me.LastName
	b.Cache.UpdateUser(me)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Logged in", fmt.Sprintf("as @%s", me.Username))
	fmt.Println("Bot logged in", fmt.Sprintf("as @%s", bt.Username))
	fmt.Println("STARTUP TIME:", time.Now().Unix()-startTime, "seconds")
	fmt.Println("ElizaUserbot - @ElizaUserbot")
	return b, bot
}

func parseINT(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
