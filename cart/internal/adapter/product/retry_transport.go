package product

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type RetryTransport struct {
	MaxAttempts int
}

var ErrTooManyRequests = errors.New("too many requests")

func NewRetryTransport(maxAttempts int) *RetryTransport {
	return &RetryTransport{MaxAttempts: maxAttempts}
}

func (rt *RetryTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	bytesBody, errRead := io.ReadAll(r.Body)
	if errRead != nil {
		panic(errRead)
	}

	originReq := *r

	var lastRes *http.Response
	var err error
	totalAttempts := 0

	for i := 0; i < rt.MaxAttempts; i++ {
		currReq := originReq
		currReq.Body = io.NopCloser(bytes.NewReader(bytesBody))

		lastRes, err = http.DefaultTransport.RoundTrip(&currReq)
		if err != nil {
			return lastRes, err
		}

		if lastRes.StatusCode == 420 || lastRes.StatusCode == http.StatusTooManyRequests {
			totalAttempts++
			continue
		} else {
			break
		}
	}

	if totalAttempts == rt.MaxAttempts {
		return lastRes, ErrTooManyRequests
	}

	return lastRes, err
}
