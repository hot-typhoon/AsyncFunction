package handler_test

import (
	"net/http"
	"net/url"
	"testing"

	handler "asynclab.club/AsyncFunction/api"
	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

func TestBridge(t *testing.T) {
	p := url.Values{}
	p.Add("url", "https://baidu.com")
	vercelkit.HttpTest(t, http.MethodGet, handler.HandlerBridge, p)
}
