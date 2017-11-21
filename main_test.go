package main

import (
	"fmt" //TODO
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func doATest(t *testing.T, operation, endpoint string, req_body io.Reader, statuscode int, res_body string) {

	r := httptest.NewRequest(operation, endpoint, req_body)
	w := httptest.NewRecorder()
	getHandler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(resp.StatusCode, string(body))
	if resp.StatusCode != 404 || string(body) != res_body {
		t.Errorf("Returning %q instead of %q", string(body), res_body)
	}
}

func TestGetHandler(t *testing.T) {
	doATest(t, "GET", "/", nil, 404, "")
	doATest(t, "GET", "/an_endpoint", nil, 404, "")
	doATest(t, "GET", "/an_endpoint/", nil, 404, "")
	// TODO: two problems here:
	// 1. the test seems wrong, the latter actually returns something
	// 2. the handler shouldn't return anything
}
