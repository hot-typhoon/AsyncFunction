package uptime

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"

	"asynclab.club/AsyncFunction/pkg/util"
)

type QueryParams struct {
	BaseUrl string
	ApiKey  string
}

type MonitorStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func GetStatusText(code string) string {
	switch code {
	case "1":
		return "UP"
	case "0":
		return "DOWN"
	case "2":
		return "PENDING"
	case "3":
		return "MAINTENANCE"
	default:
		return "UNKNOWN"
	}
}

func GetMetricsFromUptime(baseUrl string, apiKey string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, baseUrl+"/metrics", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+apiKey)))

	data, err := util.HttpRequest(req)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ExtractMetrics(data string) ([]MonitorStatus, error) {
	re := regexp.MustCompile(`^monitor_status\{.*?monitor_name="(.*?)".*?\}\s*(\d+)$`)
	lines := strings.Split(data, "\n")
	var statuses []MonitorStatus

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			name := matches[1]
			statusCode := matches[2]
			status := GetStatusText(statusCode)
			statuses = append(statuses, MonitorStatus{Name: name, Status: status})
		}
	}

	return statuses, nil
}
