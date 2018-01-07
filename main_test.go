package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var HandlerGen map[string]func(http.ResponseWriter, *http.Request)

func doATest(t *testing.T, method, endpoint string, statuscode int,
			 req_body string, res_body string) *http.Response {
	r := httptest.NewRequest(method, endpoint, strings.NewReader(req_body))
	w := httptest.NewRecorder()
	
	HandlerGen[method](w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != statuscode || string(body) != res_body {
		t.Errorf("Returning %q instead of %q", string(body), res_body)
	}
	return resp
}

func assert(t *testing.T, what, expected, given string) {
	if expected != given {
		t.Errorf("Wrong "+what+": expected "+expected+" but given "+given)
	}
}

func TestGetHandler(t *testing.T) {
	store = make(map[string]entry)
	
	_ = doATest(t, "GET", "/",404, "", "")
	_ = doATest(t, "GET", "/an_endpoint", 404, "", "")
	_ = doATest(t, "GET", "/an_endpoint/", 404, "", "")
}

func TestPostHandler(t *testing.T) {
	body := `{"this":"is an example"}`
	res_body := `{"this":"is an example"}`
	url := "/my/endpoint"
	r := doATest(t, "POST",url, 201, body, res_body)
	assert(t, "location", "/"+r.Header.Get("location"), url+"/0")
	res_body = `{"this":"is an example2"}`
	body = res_body
	r = doATest(t, "POST",url, 201, body, res_body)
	assert(t, "location", "/"+r.Header.Get("location"), url+"/1")
}

func TestMain(m *testing.M){
	HandlerGen = make(map[string]func(http.ResponseWriter, *http.Request))
	HandlerGen["GET"] = getHandler
	HandlerGen["POST"] = postHandler
	HandlerGen["PUT"] = putHandler
	HandlerGen["DELETE"] = deleteHandler
	os.Exit(m.Run())
}
