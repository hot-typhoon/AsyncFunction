package mcsm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"asynclab.club/AsyncFunction/pkg/util"
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
	Nickname string `json:"nickname"`
	Status   string `json:"status"`
	Players  int    `json:"players"`
}

type Remote struct {
	Uuid                  string
	CpuUsagePercentage    float64    `json:"cpu_usage_percentage"`
	MemoryUsagePercentage float64    `json:"memory_usage_percentage"`
	Instances             []Instance `json:"instances"`
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

	data, err := util.HttpRequest(req)
	if err != nil {
		return nil, err
	}

	var overview JOverview
	err = json.Unmarshal(data, &overview)
	if err != nil {
		return nil, err
	}

	for _, remote := range overview.Data.Remotes {
		remotes = append(remotes, Remote{Uuid: remote.Uuid, CpuUsagePercentage: remote.CpuMemCharts[len(remote.CpuMemCharts)-1].Cpu, MemoryUsagePercentage: remote.CpuMemCharts[len(remote.CpuMemCharts)-1].Mem})
	}
	//-------------------------------------------------------------------

	wg := sync.WaitGroup{}
	for i := range remotes {
		for j := -1; j <= 3; j++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				p := url.Values{}
				p.Add("apikey", apiKey)
				p.Add("daemonId", remotes[i].Uuid)
				p.Add("page", "1")
				p.Add("page_size", "10000")
				p.Add("instance_name", "")
				p.Add("status", fmt.Sprintf("%d", j))

				req, err := http.NewRequest(http.MethodGet, baseUrl+"/api/service/remote_service_instances?"+p.Encode(), nil)
				if err != nil {
					return
				}
				SetRequestHeader(req)

				data, err := util.HttpRequest(req)
				if err != nil {
					return
				}

				var instances JInstances
				err = json.Unmarshal(data, &instances)
				if err != nil {
					return
				}

				for _, instance := range instances.Data.Data {
					remotes[i].Instances = append(remotes[i].Instances, Instance{Nickname: instance.Config.Nickname, Status: StatusMap[instance.Status], Players: instance.Info.CurrentPlayers})
				}
			}(i)
		}
	}

	wg.Wait()

	return remotes, nil
}
