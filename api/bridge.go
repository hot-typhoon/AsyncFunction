package handler

import (
	"io"
	"net/http"

	"asynclab.club/AsyncFunction/pkg/lib/bridge"
	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

func HandlerBridge(w http.ResponseWriter, r *http.Request) {
	params, err := vercelkit.ReadParamsFromQuery[bridge.QueryParams](r.URL.Query())
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	req, err := http.NewRequest(r.Method, params.Url, r.Body)
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Header = r.Header
	resp, err := vercelkit.HttpClient.Do(req)
	if err != nil {
		vercelkit.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error while streaming response", http.StatusInternalServerError)
		return
	}
}
