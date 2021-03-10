package grabbertest

import (
	"net/http"
	"testing"
)

func TestHandlerDefaults(t *testing.T) {
	StartTestServer(t, func(url string) {
		resp := MustHTTPDo(t, MustHTTPNewRequest(t, "GET", url))
		AssertHTTPResponseStatusCode(t, resp, http.StatusOK)
		AssertHTTPContentLength(t, resp, 1<<20)
		AssertHTTPResponseHeader(t, resp, "Accept-Ranges", "%s", "bytes")
	})
}

func TestHandlerMethodWhitelist(t *testing.T) {
	tests := []struct {
		whiteList          []string
		method             string
		expectedStatusCode int
	}{
		{[]string{"GET", "HEAD"}, "GET", http.StatusOK},
		{[]string{"GET", "HEAD"}, "HEAD", http.StatusOK},
		{[]string{"GET"}, "HEAD", http.StatusMethodNotAllowed},
		{[]string{"HEAD"}, "GET", http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		StartTestServer(t, func(url string){
			resp := MustHTTPDo(t, MustHTTPNewRequest(t, test.method, url))
			AssertHTTPResponseStatusCode(t, resp, test.expectedStatusCode)
		}, MethodWhiteList(test.whiteList...))
	}
}

func TestHandlerHeaderBlacklist(t *testing.T) {

}