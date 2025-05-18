package handler_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	handler "asynclab.club/AsyncFunction/api"
	"github.com/dsx137/go-vercel/pkg/vercelkit"
	"github.com/joho/godotenv"
)

func TestMCSM(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Error(err)
		return
	}
	p := url.Values{}
	p.Add("base_url", os.Getenv("MCSM_BASE_URL"))
	p.Add("api_key", os.Getenv("MCSM_API_KEY"))
	vercelkit.HttpTest(t, http.MethodGet, handler.HandlerMCSM, p)
}
