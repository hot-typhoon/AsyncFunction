package handler

import (
	"net/http"

	"asynclab.club/AsyncFunction/pkg/lib/clash_plan"
	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

func HandlerClashPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		vercelkit.HttpResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	params, err := vercelkit.ReadParamsFromQuery[clash_plan.QueryParams](r.URL.Query())
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := clash_plan.GetClashPlan(params.Url)
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	vercelkit.HttpResponse(w, http.StatusOK, plan)
}
