package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSolveBasic(t *testing.T) {
	handler := Solve()

	body, _ := url.ParseQuery("rows=8&cols=8&e=&e=32&e=&e=&e=&e=16&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=8&e=&e=4&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=&e=2&e=&e=&e=&e=&e=&e=&e=&e=1&e=1")
	req, err := http.NewRequest("POST", "/solve", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = body

	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("response status code %d", rr.Code)
	}

	t.Log(rr.Body.String())
}
