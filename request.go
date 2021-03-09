package grabber

import (
	"context"
	"hash"
	"net/http"
	"net/url"
)

type Hook func(*Response) error

type Request struct {
	Label                string
	Tag                  interface{}
	HTTPRequest          *http.Request
	Filename             string
	SkipExisting         bool
	NoResume             bool
	NoStore              bool
	IgnoreBadStatusCodes bool
	IgnoreRemoteTime     bool
	Size                 int64
	BufferSize           int
	RateLimiter          RateLimiter
	BeforeCopy           Hook
	AfterCopy            Hook
	hash                 hash.Hash
	checksum             []byte
	deleteOnError        bool
	ctx                  context.Context
}

func (r *Request) NewRequest(dst, urlStr string) (*Request, error) {
	if dst == "" {
		dst = "."
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	return &Request{
		HTTPRequest: req,
		Filename:    dst,
	}, nil
}

func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}

	return context.Background()
}

func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}

	r2 := new(Request)
	*r2 = *r
	r2.ctx = ctx
	r2.HTTPRequest = r2.HTTPRequest.WithContext(ctx)
	return r2
}

func (r *Request) URL() *url.URL {
	return r.HTTPRequest.URL
}

func (r *Request) SetChecksum(h hash.Hash, sum []byte, deleteOnError bool) {
	r.hash = h
	r.checksum = sum
	r.deleteOnError = deleteOnError
}
