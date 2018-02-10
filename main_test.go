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
	if resp.StatusCode != statuscode {
		t.Errorf("%q - %q\nStatus code is %d instead of %d", method, endpoint,
			resp.StatusCode, statuscode)
	}
	if string(body) != resBody {
		t.Errorf("%q - %q\nReturning %q instead of %q", method, endpoint,
			string(body), resBody)
	}
	return resp
}

func assert(t *testing.T, what, expected, given string) {
	if expected != given {
		t.Errorf("Wrong " + what + ": expected " + expected + " but given " + given)
	}
}

func TestGetHandler(t *testing.T) {

	_ = doATest(t, "GET", "/", 404, "", "")
	_ = doATest(t, "GET", "/an_endpoint", 404, "", "")
	_ = doATest(t, "GET", "/an_endpoint/", 404, "", "")
}

func TestPostHandler(t *testing.T) {
	body := `{"this":"is an example"}`
	resBody := `{"this":"is an example"}`
	url := "/my/endpoint"
	r := doATest(t, "POST", url, 201, body, resBody)
	assert(t, "location", r.Header.Get("location"), url+"/0")
	resBody2 := `{"this":"is an example2"}`
	body = resBody2
	r = doATest(t, "POST", url, 201, body, resBody2)
	assert(t, "location", r.Header.Get("location"), url+"/1")
	// not so beautiful: this obviously depends from TestGetHandler to work
	_ = doATest(t, "GET", "/my/endpoint/0", 200, "", resBody)
	_ = doATest(t, "GET", "/my/endpoint/1", 200, "", resBody2)

	resourcesBody := `{"Resources":[0,1]}`
	_ = doATest(t, "GET", "/my/endpoint", 200, "", resourcesBody)
}

func TestPutHandler(t *testing.T) {
	body := `{"this":"was an example, now it is a test"}`
	res := doATest(t, "PUT", "/my/endpoint/0", 200, body, body)
	assert(t, "PUT location", "/my/endpoint/0", res.Header.Get("location"))
	_ = doATest(t, "PUT", "/my/endpoint/8000", 404, body, "")
}

func TestDeleteHandler(t *testing.T) {
	body := `{"this":"is an example, now it will be created"}`
	baseURL := "/just/an/endpoint"
	res0 := doATest(t, "POST", baseURL, 201, body, body)
	res := doATest(t, "POST", baseURL, 201, body, body)
	resBody := `{"Resources":[0,1]}`
	_ = doATest(t, "GET", baseURL, 200, "", resBody)
	loc := res.Header.Get("location")
	_ = doATest(t, "GET", loc, 200, body, body)
	_ = doATest(t, "DELETE", baseURL, 405, "", "")
	_ = doATest(t, "DELETE", baseURL+"/", 405, "", "")
	_ = doATest(t, "DELETE", loc, 204, "", "")
	resBody = `{"Resources":[0]}`
	_ = doATest(t, "GET", baseURL, 200, "", resBody)
	loc0 := res0.Header.Get("location")
	_ = doATest(t, "DELETE", loc0, 204, "", "")
	resBody = `{"Resources":[]}`
	_ = doATest(t, "GET", "/just/an/endpoint", 200, "", resBody)
}

func TestMain(m *testing.M) {
	HandlerGen = make(map[string]func(http.ResponseWriter, *http.Request))
	HandlerGen["GET"] = getHandler
	HandlerGen["POST"] = postHandler
	HandlerGen["PUT"] = putHandler
	HandlerGen["DELETE"] = deleteHandler
	os.Exit(m.Run())
}
