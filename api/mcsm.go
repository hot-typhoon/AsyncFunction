package handler

import (
	"net/http"

	"asynclab.club/AsyncFunction/pkg/lib/mcsm"
	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

func HandlerMCSM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		vercelkit.HttpResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	params, err := vercelkit.ReadParamsFromQuery[mcsm.QueryParams](r.URL.Query())
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	remotes, err := mcsm.GetRemotes(params.BaseUrl, params.ApiKey)
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	vercelkit.HttpResponse(w, http.StatusOK, remotes)
}
