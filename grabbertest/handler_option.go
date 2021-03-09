package grabbertest

import (
	"errors"
	"net/http"
	"time"
)

type HandlerOption func(*handler) error

func StatusCodeStatic(code int) HandlerOption {
	return func(h *handler) error {
		return StatusCode(func(req *http.Request) int {
			return code
		})(h)
	}
}

func StatusCode(f StatusCodeFunc) HandlerOption {
	return func(h *handler) error {
		if f == nil {
			return errors.New("status code function cannot be nil")
		}

		h.statusCodeFunc = f
		return nil
	}
}

func MethodWhiteList(methods ...string) HandlerOption {
	return func(h *handler) error {
		h.methodWhiteList = methods
		return nil
	}
}

func HeaderBlackList(headers ...string) HandlerOption {
	return func(h *handler) error {
		h.headerBlackList = headers
		return nil
	}
}

func ContentLength(n int) HandlerOption {
	return func(h *handler) error {
		if n < 0 {
			return errors.New("content length must be zero or greater")
		}

		h.contentLength = n
		return nil
	}
}

func AcceptRange(enabled bool) HandlerOption {
	return func(h *handler) error {
		h.acceptRanges = enabled
		return nil
	}
}

func LastModified(t time.Time) HandlerOption {
	return func(h *handler) error {
		h.lastModified = t.UTC()
		return nil
	}
}

func TimeForFirstByte(d time.Duration) HandlerOption {
	return func(h *handler) error {
		if d < 1 {
			return errors.New("time to first byte must be greater than zero")
		}

		h.ttfb = d
		return nil
	}
}

func RateLimiter(bps int) HandlerOption {
	return func(h *handler) error {
		if bps < 1 {
			return errors.New("bytes per second must be greater than zero")
		}

		h.rateLimiter = time.NewTicker(time.Second / time.Duration(bps))
		return nil
	}
}

func AttachmentFilename(filename string) HandlerOption {
	return func(h *handler) error {
		h.attachmentFilename = filename
		return nil
	}
}
