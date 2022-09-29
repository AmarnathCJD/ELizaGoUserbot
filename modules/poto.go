package modules

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("^\\+poto", Poto, &telegram.Filters{Outgoing: true})
}

func downloadPFP(user telegram.InputPeerUser, client *telegram.Client) ([]string, error) {
	p, err := client.PhotosGetUserPhotos(&telegram.InputUserObj{UserID: user.UserID, AccessHash: user.AccessHash}, 0, 0, 4)
	if err != nil {
		return nil, err
	}
	FileNames := []string{}
	switch p := p.(type) {
	case *telegram.PhotosPhotosObj:
		for _, photo := range p.Photos {
			file, err := client.DownloadProfilePhoto(user, photo)
			if err != nil {
				return nil, err
			}
			FileNames = append(FileNames, file)
		}
	}
	return FileNames, nil
}

func Poto(m *telegram.NewMessage) error {
	user, err := GetInputPeerUserFromEvent(m)
	if err != nil {
		return err
	}
	FileNames, err := downloadPFP(*user, m.Client)
	if err != nil {
		return err
	}
	fmt.Println(FileNames, "ok")
	return nil
}
