package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

func getStatusText(code string) string {
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

func extract(data string) (string, error) {
	re := regexp.MustCompile(`^monitor_status\{.*?monitor_name="(.*?)".*?\}\s*(\d+)$`)
	lines := strings.Split(data, "\n")
	var statuses []MonitorStatus

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			name := matches[1]
			statusCode := matches[2]
			status := getStatusText(statusCode)
			statuses = append(statuses, MonitorStatus{Name: name, Status: status})
		}
	}

	jsonData, err := json.MarshalIndent(statuses, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return "", err
	}

	return string(jsonData), nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.HttpResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	params, err := util.ReadParamsFromQuery[QueryParams](r.URL.Query())
	if err != nil {
		util.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	req, err := http.NewRequest("GET", params.BaseUrl+"/metrics", nil)
	if err != nil {
		util.HttpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+params.ApiKey)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.HttpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		util.HttpResponse(w, resp.StatusCode, "Failed to fetch metrics")
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		util.HttpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonData, err := extract(string(data))
	if err != nil {
		util.HttpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.HttpResponse(w, http.StatusOK, jsonData)
}
