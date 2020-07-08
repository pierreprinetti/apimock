package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Integration test
func TestMain(t *testing.T) {
	t.Run("complete PUT-GET-DELETE-GET cycle", func(t *testing.T) {

		// Run the application
		srvAddr := "localhost:29109"
		os.Setenv("HOST", srvAddr)
		defer os.Unsetenv("HOST")

		go func() {
			main()
		}()

		// Make sure that the http listener is in place
		time.Sleep(time.Millisecond)

		// Define the data that will be sent and the expected
		targetEndpoint := "http://" + srvAddr + "/endpoint1"
		expectedBody := "yay"
		expectedContentType := "WUARGH"

		// Perform the PUT call
		var client http.Client
		req, _ := http.NewRequest("PUT", targetEndpoint, strings.NewReader(expectedBody))
		req.Header.Set("Content-Type", expectedContentType)
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("calling PUT: %v", err)
		}

		// Test the PUT response code
		if want, have := 200, res.StatusCode; want != have {
			t.Errorf("expected PUT response status code %d, found %d", want, have)
		}

		// Perform the GET call
		res, err = http.Get(targetEndpoint)
		if err != nil {
			t.Fatalf("calling GET: %v", err)
		}

		// Test the GET response code
		if want, have := 200, res.StatusCode; want != have {
			t.Errorf("expected GET response status code %d, found %d", want, have)
		}

		// Test the GET response body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("reading the GET response body: %v", err)
		}

		if want, have := expectedBody, string(body); want != have {
			t.Errorf("expected GET response body %q, found %q", want, have)
		}

		// Test the GET response content-type
		if want, have := expectedContentType, res.Header.Get("Content-Type"); want != have {
			t.Errorf("expected GET response content-type %q, found %q", want, have)
		}

		// Perform the DELETE call
		req, _ = http.NewRequest("DELETE", targetEndpoint, nil)
		res, err = client.Do(req)
		if err != nil {
			t.Fatalf("calling DELETE: %v", err)
		}

		// Test the DELETE response code
		if want, have := 204, res.StatusCode; want != have {
			t.Errorf("expected DELETE response status code %d, found %d", want, have)
		}

		// Perform a GET call to check if the resource has been deleted
		res, err = http.Get(targetEndpoint)
		if err != nil {
			t.Fatalf("calling GET: %v", err)
		}

		// Test the GET response code
		if want, have := 404, res.StatusCode; want != have {
			t.Errorf("expected GET response status code %d, found %d", want, have)
		}
	})

	t.Run("OPTIONS call", func(t *testing.T) {

		// Run the application
		srvAddr := "localhost:29108"
		os.Setenv("HOST", srvAddr)
		defer os.Unsetenv("HOST")

		go func() {
			main()
		}()

		// Make sure that the http listener is in place
		time.Sleep(time.Millisecond)

		// Define the data that will be sent and the expected
		targetEndpoint := "http://" + srvAddr + "/endpoint2"

		// Perform the OPTIONS call
		var client http.Client
		req, _ := http.NewRequest("OPTIONS", targetEndpoint, nil)
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("calling OPTIONS: %v", err)
		}

		// Test the response code
		if want, have := 204, res.StatusCode; want != have {
			t.Errorf("expected response status code %d, found %d", want, have)
		}
	})

	t.Run("POST call not implemented", func(t *testing.T) {

		// Run the application
		srvAddr := "localhost:29110"
		os.Setenv("HOST", srvAddr)
		defer os.Unsetenv("HOST")

		go func() {
			main()
		}()

		// Make sure that the http listener is in place
		time.Sleep(time.Millisecond)

		// Define the data that will be sent and the expected
		targetEndpoint := "http://" + srvAddr + "/endpoint3"

		// Perform the OPTIONS call
		var client http.Client
		req, _ := http.NewRequest("POST", targetEndpoint, nil)
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("calling POST: %v", err)
		}

		// Test the response code
		if want, have := 501, res.StatusCode; want != have {
			t.Errorf("expected response status code %d, found %d", want, have)
		}
	})
}
