package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("^\\+dl", downloadFile, &telegram.Filters{Outgoing: true})
}

func downloadFile(m *telegram.NewMessage) error {
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			return err
		}
		if !r.IsMedia() {
			_, err = m.Edit("<code>Reply to a media message to download it.</code>")
			return err
		}
		_, err = m.Edit("<code>Downloading...</code>")
		if err != nil {
			return err
		}
		p := PROGRESS_GEN{StartTime: time.Now(), LastTime: time.Now()}
		prog := telegram.Progress{}
		go func() {
			for {
				if prog.Total != 0 {
					p.TotalSize = int64(prog.Total)
                                        p.CurrentSize = int64(0)
                                        p.LastSize = int64(0)
					break
				}
			}
			for {
				time.Sleep(2 * time.Second)
				p.UpdateProgress(int(prog.Current))
				_, err = m.Edit(p.GenProgressString())
				if err != nil {
					return
				}
			}
		}()

		//prog.OnStart = func() {
		//	p.StartProgress()
		// }
		f, _ := m.Client.DownloadMedia(r, &telegram.DownloadOptions{Progress: &prog, Threaded: true, Threads: 20})
		p.EndProgress()
		_, err = m.Edit(fmt.Sprintf("<code>Downloaded to %s</code> in %s", f, p.GetElapsedTime()))
		if err != nil {
			return err
		}
		return err
	}
	_, err := m.Edit("<code>Reply to a media message to download it.</code>")
	return err
}

type PROGRESS_GEN struct {
	TotalSize   int64
	CurrentSize int64
	StartTime   time.Time
	LastTime    time.Time
	LastSize    int64
}

func (p *PROGRESS_GEN) GetProgress() string {
	//elapsed := time.Since(p.StartTime)
	elapsedLast := time.Since(p.LastTime)
	speed := float64(p.CurrentSize-p.LastSize) / elapsedLast.Seconds()
	if speed == 0 {
		speed = 1
	}
	remaining := time.Duration((p.TotalSize-p.CurrentSize)/int64(speed)) * time.Second
	return fmt.Sprintf("%s / %s (%.2f%%) - %.2f KB/s - %s remaining", FormatBytes(p.CurrentSize), FormatBytes(p.TotalSize), float64(p.CurrentSize)/float64(p.TotalSize)*100, speed/1024, remaining)
}

func FormatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/1024/1024/1024)
	}
}

func (p *PROGRESS_GEN) UpdateProgress(size int) {
	p.CurrentSize += int64(size)
        fmt.Println(p.CurrentSize, p.LastSize)
	p.LastSize = p.CurrentSize
	p.LastTime = time.Now()
}

func (p *PROGRESS_GEN) StartProgress() {
	p.StartTime = time.Now()
	p.LastTime = time.Now()
}

func (p *PROGRESS_GEN) EndProgress() {
	p.CurrentSize = p.TotalSize
	p.LastSize = p.CurrentSize
	p.LastTime = time.Now()
}

func (p *PROGRESS_GEN) GetProgressString() string {
	return fmt.Sprintf("%s / %s (%.2f%%)", FormatBytes(p.CurrentSize), FormatBytes(p.TotalSize), float64(p.CurrentSize)/float64(p.TotalSize)*100)
}

func (p *PROGRESS_GEN) GetSpeedString() string {
	elapsed := time.Since(p.StartTime)
	speed := float64(p.CurrentSize) / elapsed.Seconds()
	if speed == 0 {
		speed = 1
	}
	return fmt.Sprintf("%.2f KB/s", speed/1024)
}

func (p *PROGRESS_GEN) GetRemainingString() string {
	elapsed := time.Since(p.StartTime)
	speed := float64(p.CurrentSize) / elapsed.Seconds()
	if speed == 0 {
		speed = 1
	}
	remaining := time.Duration((p.TotalSize-p.CurrentSize)/int64(speed)) * time.Second
	return fmt.Sprintf("%s remaining", remaining)
}

func (p *PROGRESS_GEN) GetETAString() string {
	elapsed := time.Since(p.StartTime)
	speed := float64(p.CurrentSize) / elapsed.Seconds()
	if speed == 0 {
		speed = 1
	}
	remaining := time.Duration((p.TotalSize-p.CurrentSize)/int64(speed)) * time.Second
	return time.Now().Add(remaining).Format("15:04:05")
}

func (p *PROGRESS_GEN) GetElapsedTime() string {
	elapsed := time.Since(p.StartTime)
	return elapsed.String()
}

func (p *PROGRESS_GEN) GetElapsedTimeSeconds() int64 {
	elapsed := time.Since(p.StartTime)
	return int64(elapsed.Seconds())
}

func (p *PROGRESS_GEN) GenProgressString() string {
	return fmt.Sprintf("<b>Downloading... %s</b>\n<b>Progress:</b> %s\n<b>Speed:</b> %s\n<b>Remaining:</b> %s\n<b>ETA:</b> %s\n<b>Elapsed:</b> %s", "-", p.GetProgress(), p.GetSpeedString(), p.GetRemainingString(), p.GetETAString(), p.GetElapsedTime())
}
