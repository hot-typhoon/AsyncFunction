package clash_plan

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

type QueryParams struct {
	Url string
}

type ClashPlan struct {
	Upload   string `json:"upload"`
	Download string `json:"download"`
	Total    string `json:"total"`
	Expire   string `json:"expire"`
}

func GetClashPlan(url string) (*ClashPlan, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "clash-verge/v2.0.2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %v", resp.Status)
	}

	info := resp.Header.Get("Subscription-Userinfo")
	if info == "" {
		return nil, fmt.Errorf("error: %v", "Subscription-Userinfo not found")
	}

	var upload, download, total, expire int
	fmt.Sscanf(info, "upload=%d; download=%d; total=%d; expire=%d", &upload, &download, &total, &expire)

	return &ClashPlan{
		Upload:   vercelkit.ConvertBytesToHuman(upload),
		Download: vercelkit.ConvertBytesToHuman(download),
		Total:    vercelkit.ConvertBytesToHuman(total),
		Expire:   time.Unix(int64(expire), 0).Local().Format("2006-01-02 15:04:05"),
	}, nil
}
