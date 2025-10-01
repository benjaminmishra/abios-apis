package abios

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	abiosAuthHeaderKey  = "Abios-Secret"
	retryAfterHeaderKey = "Retry-After"
)

type authTransport struct {
	token     string
	transport http.RoundTripper
}

type rateLimitTransport struct {
	limiter   *rate.Limiter
	transport http.RoundTripper
}

type retryTransport struct {
	transport  http.RoundTripper
	maxRetries int
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(abiosAuthHeaderKey, a.token)

	return a.transport.RoundTrip(req)
}

func (t *rateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return t.transport.RoundTrip(req)
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastResp *http.Response
	var err error

	for attempt := 0; attempt < t.maxRetries; attempt++ {
		// cloning request so body/headers are not consumed
		newReq := req.Clone(req.Context())

		lastResp, err = t.transport.RoundTrip(newReq)
		if err != nil {
			return nil, err
		}

		if lastResp.StatusCode != http.StatusTooManyRequests {
			return lastResp, nil
		}

		retryAfter := parseRetryAfter(lastResp.Header.Get(retryAfterHeaderKey))
		time.Sleep(retryAfter)
	}

	return lastResp, fmt.Errorf("abios: too many retries after %d attempts", t.maxRetries)
}
