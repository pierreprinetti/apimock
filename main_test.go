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
			 reqBody string, resBody string) *http.Response {
	r := httptest.NewRequest(method, endpoint, strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	
	HandlerGen[method](w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != statuscode || string(body) != resBody {
		t.Errorf("Returning %q instead of %q", string(body), resBody)
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
	resBody := `{"this":"is an example"}`
	url := "/my/endpoint"
	r := doATest(t, "POST",url, 201, body, resBody)
	assert(t, "location", r.Header.Get("location"), url+"/0")
	resBody2 := `{"this":"is an example2"}`
	body = resBody2
	r = doATest(t, "POST",url, 201, body, resBody2)
	assert(t, "location", r.Header.Get("location"), url+"/1")
	_ = doATest(t, "GET", "/my/endpoint/1", 200, "", resBody2)
	_ = doATest(t, "GET", "/my/endpoint/0", 200, "", resBody)

  resourcesBody := "{\"Resources\":[0,1]}"
	_ = doATest(t, "GET", "/my/endpoint", 200, "", resourcesBody)
}

func TestMain(m *testing.M){
	HandlerGen = make(map[string]func(http.ResponseWriter, *http.Request))
	HandlerGen["GET"] = getHandler
	HandlerGen["POST"] = postHandler
	HandlerGen["PUT"] = putHandler
	HandlerGen["DELETE"] = deleteHandler
	os.Exit(m.Run())
}
