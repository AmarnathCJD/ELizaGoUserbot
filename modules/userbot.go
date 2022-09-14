package modules

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/joho/godotenv"
)

var (
	UB        = UserbotClient()
	startTime int64
)

func UserbotClient() *telegram.Client {
	godotenv.Load()
	STRING_SESSION := os.Getenv("STRING_SESSION")
	fmt.Println("Starting Userbot...")
	startTime = time.Now().Unix()
	API_ID := os.Getenv("API_ID")
	API_HASH := os.Getenv("API_HASH")
	DC_ID := os.Getenv("DC_ID")
	b, _ := telegram.TelegramClient(telegram.ClientConfig{
		AppID:         parseINT(API_ID),
		AppHash:       API_HASH,
		DataCenter:    parseINT(DC_ID),
		StringSession: STRING_SESSION,
		ParseMode:     "HTML",
	})
	me, err := b.GetMe()
	b.Cache.UpdateUser(me)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Logged in", fmt.Sprintf("as @%s", me.Username))
	return b
}

func parseINT(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
