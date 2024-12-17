package handler

import (
	"net/http"

	"asynclab.club/AsyncFunction/pkg/lib/mcsm"
	"asynclab.club/AsyncFunction/pkg/util"
)

func HandlerMCSM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.HttpResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	params, err := util.ReadParamsFromQuery[mcsm.QueryParams](r.URL.Query())
	if err != nil {
		util.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	remotes, err := mcsm.GetRemotes(params.BaseUrl, params.ApiKey)
	if err != nil {
		util.HttpResponse(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	util.HttpResponse(w, http.StatusOK, remotes)
}
