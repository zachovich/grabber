package grabbertest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MustHTTPDo(t *testing.T, req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func MustHTTPNewRequestCtx(t *testing.T, ctx context.Context, method, url string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		panic(err)
	}

	return req
}

func MustHTTPNewRequest(t *testing.T, method, url string) *http.Request {
	return MustHTTPNewRequestCtx(t, context.Background(), method, url)
}

func StartTestServer(t *testing.T, f func(url string), options ...HandlerOption) {
	h, err := NewHandler(options...)
	if err != nil {
		t.Fatalf("unable to create test server handler: %v", err)
		return
	}

	s := httptest.NewServer(h)
	defer func() {
		h.(*handler).close()
		s.Close()
	}()

	f(s.URL)
}

