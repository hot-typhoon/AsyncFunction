package mcsm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

type QueryParams struct {
	BaseUrl string
	ApiKey  string
}

type JCpuMemChart struct {
	Cpu float64 `json:"cpu"`
	Mem float64 `json:"mem"`
}

type JRemote struct {
	Uuid         string         `json:"uuid"`
	CpuMemCharts []JCpuMemChart `json:"cpuMemChart"`
	Remarks      string         `json:"remarks"`
}

type JOverviewData struct {
	Remotes []JRemote `json:"remote"`
}

type JOverview struct {
	Data JOverviewData `json:"data"`
}

type JInstanceConfig struct {
	Nickname string `json:"nickname"`
}

type JInstanceInfo struct {
	CurrentPlayers int `json:"currentPlayers"`
}

type JInstance struct {
	Status int             `json:"status"`
	Config JInstanceConfig `json:"config"`
	Info   JInstanceInfo   `json:"info"`
}

type JInstancesData struct {
	Data []JInstance `json:"data"`
}

type JInstances struct {
	Data JInstancesData `json:"data"`
}

var StatusMap = map[int]string{
	-1: "BUSY",
	0:  "OFFLINE",
	1:  "STOPPING",
	2:  "STARTING",
	3:  "ONLINE",
}

type Instance struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Players int    `json:"players"`
}

type Remote struct {
	Uuid      string     `json:"-"`
	Name      string     `json:"name"`
	Cpu       string     `json:"cpu"`
	Memory    string     `json:"memory"`
	Instances []Instance `json:"instances"`
}

func SetRequestHeader(req *http.Request) {
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
}

func GetRemotes(baseUrl, apiKey string) ([]Remote, error) {
	remotes := make([]Remote, 0)

	p := url.Values{}
	p.Add("apikey", apiKey)

	req, err := http.NewRequest(http.MethodGet, baseUrl+"/api/overview?"+p.Encode(), nil)
	if err != nil {
		return nil, err
	}
	SetRequestHeader(req)

	data, err := vercelkit.HttpRequest(req)
	if err != nil {
		return nil, err
	}

	var overview JOverview
	err = json.Unmarshal(data, &overview)
	if err != nil {
		return nil, err
	}

	//-------------------------------------------------------------------

	wg := sync.WaitGroup{}
	for i, remoteData := range overview.Data.Remotes {
		remotes = append(remotes, Remote{Uuid: remoteData.Uuid, Name: remoteData.Remarks, Cpu: fmt.Sprintf("%.2f", remoteData.CpuMemCharts[0].Cpu) + "%", Memory: fmt.Sprintf("%.2f", remoteData.CpuMemCharts[0].Mem) + "%"})
		for j := -1; j <= 3; j++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				p := url.Values{}
				p.Add("apikey", apiKey)
				p.Add("daemonId", remoteData.Uuid)
				p.Add("page", "1")
				p.Add("page_size", "10000")
				p.Add("instance_name", "")
				p.Add("status", fmt.Sprintf("%d", j))

				req, err := http.NewRequest(http.MethodGet, baseUrl+"/api/service/remote_service_instances?"+p.Encode(), nil)
				if err != nil {
					return
				}
				SetRequestHeader(req)

				data, err := vercelkit.HttpRequest(req)
				if err != nil {
					return
				}

				var instances JInstances
				err = json.Unmarshal(data, &instances)
				if err != nil {
					return
				}

				for _, instance := range instances.Data.Data {
					remotes[i].Instances = append(remotes[i].Instances, Instance{Name: instance.Config.Nickname, Status: StatusMap[instance.Status], Players: instance.Info.CurrentPlayers})
				}
			}(i)
		}
	}

	wg.Wait()

	return remotes, nil
}
