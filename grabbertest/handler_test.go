package grabbertest

import (
	"testing"
)

func TestHandlerDefaults(t *testing.T) {
	StartTestServer(t, func(url string) {
		resp := MustHTTPDo(t, MustHTTPNewRequest(t, "GET", url))
		AssertHTTPResponseStatusCode(t, resp, 200)
	})
}