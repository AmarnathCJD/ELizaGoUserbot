package modules

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/amarnathcjd/gogram/telegram"
)

func init() {
	UB.AddMessageHandler("\\+json", Jsonify, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+reserved", ReservedChannels, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+sh", Shell, &telegram.Filters{Outgoing: true})
	UB.AddMessageHandler("\\+sysinfo", SystemInfo, &telegram.Filters{Outgoing: true})
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

func SystemInfo(m *telegram.NewMessage) error {
	info, err := GatherPsutilInfo()
	if err != nil {
		return err
	}
	_, err = m.Edit(info)
	return err
}

func GatherPsutilInfo() (string, error) {
	var s string = "<b>System Info:</b>\n"
	kernel, err := host.KernelVersion()
	if err != nil {
		return s, err
	}
	h, _ := host.Info()
	ip, err := externalIP()
	if err != nil {
		return s, err
	}
	OS := fmt.Sprintf("%s %s %s", runtime.GOOS, runtime.GOARCH, runtime.Version())
	s += fmt.Sprintf("<b>Kernel:</b> <code>%s</code>\n<b>Hostname:</b> <code>%s</code>\n<b>IP:</b> <code>%s</code>\n<b>OS:</b> <code>%s</code>", kernel, h.Hostname, ip, OS)
	v, err := mem.VirtualMemory()
	if err != nil {
		return s, err
	}
	s += fmt.Sprintf("\n<b>RAM:</b> <code>%v/%v (%f%%)</code>\n", HumanBytes(v.Used), HumanBytes(v.Total), v.UsedPercent)
	c, err := cpu.Info()
	if err != nil {
		return s, err
	}
	s += fmt.Sprintf("<b>CPU:</b> <code>%s (x%d)</code>\n", c[0].ModelName, len(c))
	p, err := cpu.Percent(0, false)
	if err != nil {
		return s, err
	}
	s += fmt.Sprintf("<b>CPU Usage:</b> <code>%f%%</code>\n", p[0])
	d, err := disk.Usage("/")
	if err != nil {
		return s, err
	}
	s += fmt.Sprintf("<b>Disk:</b> <code>%s/%s (%f%%)</code>\n", HumanBytes(d.Used), HumanBytes(d.Total), d.UsedPercent)
	n, err := net.IOCounters(true)
	if err != nil {
		return s, err
	}
	totalOut := uint64(0)
	totalIn := uint64(0)
	for _, v := range n {
		totalOut += v.BytesSent
		totalIn += v.BytesRecv
	}
	s += fmt.Sprintf("<b>Network:</b> <code>%s (in) / %s (out)</code>", HumanBytes(totalIn), HumanBytes(totalOut))
	return s, nil
}

func externalIP() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HumanBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
