package grabbertest

import (
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

