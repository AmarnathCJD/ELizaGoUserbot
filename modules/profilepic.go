package modules

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("\\+autopfp", AutoPfp, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+stopautopfp", StopAutoPfp, &telegram.Filters{Outgoing: true})
        UB.AddMessageHandler("\\+pfp", we, &telegram.Filters{Outgoing: true})
}

func we(m *telegram.NewMessage) error {
 a := m.Args()
 messages, err := m.Client.GetMessages(a, &telegram.SearchOption{Limit: 2,
		Filter: &telegram.InputMessagesFilterChatPhotos{}})
	if err != nil {
		return err
	}
 var p []telegram.InputMedia
 for _, x := range messages {
     v := x.Action.(*telegram.MessageActionChatEditPhoto).Photo.(*telegram.PhotoObj)
     p = append(p, &telegram.InputMediaPhoto{ID: &telegram.InputPhotoObj{ID: v.ID, AccessHash: v.AccessHash, FileReference: v.FileReference}, TtlSeconds: 0})
 }
 m.Client.SendAlbum(m.ChatID(), p)
 return nil
}


func ParseWallpaperURLS(query string) []string {
	API := "https://getwallpapers.com/search?term=" + url.QueryEscape(query)
	resp, err := http.Get(API)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil
	}
	var urls []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if strings.Contains(url, "wallpaper/full") {
			urls = append(urls, "https://getwallpapers.com"+url)
		}
	})
	return urls
}

var (
	ctx = context.Background()
)

func AutoPfp(m *telegram.NewMessage) error {
	Query := m.Args()
	if Query == "" {
		m.Reply("Please provide a search query!")
		return nil
	}
	urls := ParseWallpaperURLS(Query)
	if len(urls) == 0 {
		m.Reply("No results found!")
		return nil
	}
	ctx = context.Background()
	go func() {
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			default:
				pic, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				file, err := os.Create("pfp.jpg")
				if err != nil {
					fmt.Println(err)
					return
				}
				io.Copy(file, pic.Body)
				file.Close()
				p, err := m.Client.UploadFile("pfp.jpg")
				if err != nil {
					fmt.Println(err)
				}
				_, e := m.Client.PhotosUploadProfilePhoto(p, nil, 0)
				if e != nil {
					fmt.Println(e)

				}
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(60 * time.Second)
			}
		}
	}()
	m.Edit("AutoPfp started!")
	return nil
}

func StopAutoPfp(m *telegram.NewMessage) error {
	ctx.Done()
	m.Edit("AutoPfp stopped!")
	return nil
}
