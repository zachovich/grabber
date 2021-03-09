package grabbertest

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	DefaultHandlerContentLength = 1 << 20 // 1024 KB
)

type StatusCodeFunc func(req *http.Request) int

type handler struct {
	statusCodeFunc     StatusCodeFunc
	methodWhiteList    []string
	headerBlackList    []string
	contentLength      int
	acceptRanges       bool
	attachmentFilename string
	lastModified       time.Time
	ttfb               time.Duration // time to first byte
	rateLimiter        *time.Ticker
}

func NewHandler(options ...HandlerOption) (http.Handler, error) {
	h := &handler{
		statusCodeFunc:  func(req *http.Request) int { return http.StatusOK },
		methodWhiteList: []string{"GET", "HEAD"},
		contentLength:   DefaultHandlerContentLength,
		acceptRanges:    true,
	}

	for _, option := range options {
		if err := option(h); err != nil {
			return nil, err
		}
	}

	return h, nil
}

func (h *handler) close() {
	if h.rateLimiter != nil {
		h.rateLimiter.Stop()
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// delay response
	if h.ttfb > 0 {
		time.Sleep(h.ttfb)
	}

	allowed := false
	for _, method := range h.methodWhiteList {
		if r.Method == method {
			allowed = true
			break
		}
	}

	if !allowed {
		httpErr(w, http.StatusMethodNotAllowed)
		return
	}

	// set server options
	if h.acceptRanges {
		w.Header().Set("Accept-Ranges", "bytes")
	}

	// set attachment filename
	if h.attachmentFilename != "" {
		w.Header().Set(
			"Content-Disposition",
			fmt.Sprintf("attachment;filename=%s", h.attachmentFilename),
		)
	}

	// set last modified timestamp
	lastMod := time.Now()
	if !h.lastModified.IsZero() {
		lastMod = h.lastModified
	}

	w.Header().Set("Last-Modified", lastMod.UTC().Format(http.TimeFormat))

	// set content-length
	offset := 0
	if h.acceptRanges {
		if reqRange := r.Header.Get("Range"); reqRange != "" {
			if _, err := fmt.Sscan(reqRange, "bytes=%d-", &offset); err != nil {
				httpErr(w, http.StatusBadRequest)
				return
			}

			if offset >= h.contentLength || offset < 0 {
				httpErr(w, http.StatusRequestedRangeNotSatisfiable)
				return
			}
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(h.contentLength-offset))

	// apply header blacklist
	if h.headerBlackList != nil {
		for _, key := range h.headerBlackList {
			w.Header().Del(key)
		}
	}

	// send header and status code
	w.WriteHeader(h.statusCodeFunc(r))

	// send body
	if r.Method == "GET" {
		// use buffered io to reduce overhead on the reader
		bw := bufio.NewWriterSize(w, 4096)
		for i := offset; !isRequestClosed(r) && i < h.contentLength; i++ {
			bw.WriteByte(byte(i))

			if h.rateLimiter != nil {
				bw.Flush()
				w.(http.Flusher).Flush() // force the server to send the data to the client
				select {
				case <-h.rateLimiter.C:
				case <-r.Context().Done():
				}
			}
		}

		if !isRequestClosed(r) {
			bw.Flush()
		}
	}
}

func isRequestClosed(r *http.Request) bool {
	return r.Context().Err() != nil
}

func httpErr(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
