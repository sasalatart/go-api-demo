package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type PongTestCase struct {
	desc   string
	method string
	query  string
	body   string
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err.Error())
	}
}

func assertResponseStatus(t *testing.T, expected int, got int) {
	t.Helper()

	if expected != got {
		t.Errorf("Expected response status to be %d, but got %d", expected, got)
	}
}

func assertResponseContentType(t *testing.T, expected string, got string) {
	t.Helper()

	if expected != got {
		t.Errorf("Expected response headers to have Content-Type %s, but got %s", expected, got)
	}
}

func assertResponseMethod(t *testing.T, expected string, response string) {
	t.Helper()

	if expected != response {
		t.Errorf("Expected %s, but got %s", expected, response)
	}
}

func assertKeyValues(t *testing.T, key string, expectedValue string, gotValue string) {
	t.Helper()

	if expectedValue != gotValue {
		t.Errorf("Expected key %s to have value %v, but got %v", key, expectedValue, gotValue)
	}
}

func deepCompare(t *testing.T, sent string, response interface{}) {
	t.Helper()

	expected, err := url.ParseQuery(sent)
	assertNoError(t, err)

	got := response.(map[string]interface{})

	for key, expectedValue := range expected {
		assertKeyValues(t, key, fmt.Sprintf("%v", expectedValue), fmt.Sprintf("%v", got[key]))
	}

	for key, gotValue := range got {
		assertKeyValues(t, key, fmt.Sprintf("%v", expected[key]), fmt.Sprintf("%v", gotValue))
	}
}

func doRequest(t *testing.T, testCase PongTestCase) *http.Request {
	t.Helper()

	requestURL := "/ping"
	if testCase.query != "" {
		requestURL = requestURL + "?" + testCase.query
	}

	var requestBody io.Reader
	if testCase.body != "" {
		requestBody = strings.NewReader(testCase.body)
	}

	request, err := http.NewRequest(testCase.method, requestURL, requestBody)
	assertNoError(t, err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return request
}

func TestPongHandler(t *testing.T) {
	testCases := []PongTestCase{
		{
			desc:   "GET with no query string nor body",
			method: http.MethodGet,
		},
		{
			desc:   "GET with a query string",
			method: http.MethodGet,
			query:  "key1=value1&key2=value2",
		},
		{
			desc:   "POST with no query string nor body",
			method: http.MethodPost,
		},
		{
			desc:   "POST with a query string but no body",
			method: http.MethodPost,
			query:  "key1=value1&key2=value2",
		},
		{
			desc:   "POST with a query string and a body",
			method: http.MethodPost,
			query:  "key1=value1&key2=value2",
			body:   "key3=value3&key4=value4",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			request := doRequest(t, testCase)
			response := httptest.NewRecorder()

			if testCase.body != "" {
				defer request.Body.Close()
			}

			PongHandler(response, request)

			responseBody, err := ioutil.ReadAll(response.Body)
			assertNoError(t, err)
			var pongResponse PongResponse
			if err := json.Unmarshal(responseBody, &pongResponse); err != nil {
				t.Fatalf(err.Error())
			}

			assertResponseStatus(t, http.StatusOK, response.Code)
			assertResponseContentType(t, "application/json", response.HeaderMap["Content-Type"][0])
			assertResponseMethod(t, testCase.method, pongResponse.Method)

			if testCase.query != "" {
				deepCompare(t, testCase.query, pongResponse.QueryParams)
			}

			if testCase.body != "" {
				deepCompare(t, testCase.body, pongResponse.Body)
			}
		})
	}
}
