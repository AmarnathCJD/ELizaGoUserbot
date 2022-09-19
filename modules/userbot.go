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
)

func UserbotClient() (*telegram.Client, *telegram.Client) {
	godotenv.Load()
	STRING_SESSION := os.Getenv("STRING_SESSION")
	fmt.Println("Starting Userbot...")
	startTime = time.Now().Unix()
	API_ID := os.Getenv("API_ID")
	API_HASH := os.Getenv("API_HASH")
	DC_ID := os.Getenv("DC_ID")
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	b, _ := telegram.TelegramClient(telegram.ClientConfig{
		AppID:         parseINT(API_ID),
		AppHash:       API_HASH,
		DataCenter:    parseINT(DC_ID),
		StringSession: STRING_SESSION,
		ParseMode:     "HTML",
	})
	bot, _ := telegram.TelegramClient(telegram.ClientConfig{
		AppID:       parseINT(API_ID),
		AppHash:     API_HASH,
		DataCenter:  parseINT(DC_ID),
		ParseMode:   "HTML",
		SessionFile: filepath.Join(WorkDir, "bot.session"),
	})
	if err := bot.LoginBot(BOT_TOKEN); err != nil {
		panic(err)
	}
	me, err := b.GetMe()
	bt, _ := bot.GetMe()
	BotUName = bt.Username
	b.Cache.UpdateUser(me)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Logged in", fmt.Sprintf("as @%s", me.Username))
	fmt.Println("Bot logged in", fmt.Sprintf("as @%s", bt.Username))
	fmt.Println("Userbot started!")
	return b, bot
}

func parseINT(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
