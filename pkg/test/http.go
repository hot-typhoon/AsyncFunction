package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func HttpTest(t *testing.T, method string, handler func(http.ResponseWriter, *http.Request), params url.Values) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "https://test.com?"+params.Encode(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Error("Status code: " + fmt.Sprintf("%d", w.Code) + ", Body: " + w.Body.String())
		return
	}
}
