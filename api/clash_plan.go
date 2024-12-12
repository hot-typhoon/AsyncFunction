package handler

import (
	"net/http"

	"asynclab.club/AsyncFunction/pkg/lib/clash_plan"
	"asynclab.club/AsyncFunction/pkg/util"
)

func HandlerClashPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.HttpResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	params, err := util.ReadParamsFromQuery[clash_plan.QueryParams](r.URL.Query())
	if err != nil {
		util.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := clash_plan.GetClashPlan(params.Url)
	if err != nil {
		util.HttpResponse(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	util.HttpResponse(w, http.StatusOK, plan)
}
