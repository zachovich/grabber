package grabbertest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func AssertHTTPResponseStatusCode(t *testing.T, resp *http.Response, expected int) bool {
	actual := resp.StatusCode
	if actual != expected {
		t.Errorf("wrong response status code. extected: %d, got: %d", expected, actual)
		return false
	}

	return true
}

func AssertHTTPResponseHeader(t *testing.T, resp *http.Response, key, format string, a ...interface{}) bool {
	expected := fmt.Sprintf(format, a...)
	actual := resp.Header.Get(key)

	if expected != actual {
		t.Errorf("wrong response header. expected: %s, got: %s", expected, actual)
		return false
	}

	return true
}

func AssertHTTPContentBodyLength(t *testing.T, resp *http.Response, expected int64) bool {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if int64(len(body)) != expected {
		t.Errorf("wrong body length. exptected: %d, got: %d", expected, len(body))
		return false
	}

	return true
}

func AssertHTTPContentLength(t *testing.T, resp *http.Response, expected int64) bool {
	actual := resp.ContentLength
	if actual != expected {
		t.Errorf("wrong content-length. expected: %d, got: %d", expected, actual)
		return false
	}

	if ! AssertHTTPContentBodyLength(t, resp, expected) {
		return false
	}

	return true
}
