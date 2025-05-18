package handler

import (
	"net/http"

	"asynclab.club/AsyncFunction/pkg/lib/uptime"
	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

func HandlerUptime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		vercelkit.HttpResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	params, err := vercelkit.ReadParamsFromQuery[uptime.QueryParams](r.URL.Query())
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := uptime.GetMetricsFromUptime(params.BaseUrl, params.ApiKey)
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	statuses, err := uptime.ExtractMetrics(data)
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	vercelkit.HttpResponse(w, http.StatusOK, statuses)
}
