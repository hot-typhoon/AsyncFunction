package handler

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"asynclab.club/AsyncFunction/pkg/test"
	"github.com/joho/godotenv"
)

func TestUptime(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Error(err)
		return
	}
	p := url.Values{}
	p.Add("base_url", os.Getenv("UPTIME_BASE_URL"))
	p.Add("api_key", os.Getenv("UPTIME_API_KEY"))
	test.HttpTest(t, http.MethodGet, HandlerUptime, p)
}
