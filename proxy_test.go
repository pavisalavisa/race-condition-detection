package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fleetName      = "testfleet"
	goodTenantSlug = "goodone"
	inputBaseURL   = "https://api.com"
	systemABaseUrl = "https://system-a.com"
	systemBBaseURL = "https://system-b.com"
)

func TestProxy_V1(t *testing.T) {
	testCases := []struct {
		desc             string
		originalPath     string
		originalMethod   string
		expectedProxyURL string
	}{
		{
			desc:             "System A POST",
			originalPath:     "/a/foo/bar",
			originalMethod:   "POST",
			expectedProxyURL: fmt.Sprintf("%s/v1/foo/bar", systemABaseUrl),
		},
		{
			desc:             "System B POST",
			originalPath:     "/b/baz/14",
			originalMethod:   "POST",
			expectedProxyURL: fmt.Sprintf("%s/baz/14", systemBBaseURL),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var proxiedRequest *http.Request
			p := fixtureProxy(t, func(r *http.Request) {
				proxiedRequest = r
			})

			writer := fixtureWriter()
			req := fixtureRequest(t, tC.originalPath, tC.originalMethod)
			p.ServeHTTP(writer, req)
			require.Equal(t, tC.expectedProxyURL, proxiedRequest.URL.String())
			require.Equal(t, tC.originalMethod, proxiedRequest.Method, "HTTP method should not be modified on proxy")
		})
	}
}

func TestProxy_newProxy_errorProxying(t *testing.T) {
	p := fixtureProxy(t, func(r *http.Request) {})
	p.Transport = noopRoundTripper{onRoundTrip: errorRoundTrip}

	writer := fixtureWriter()
	req := fixtureRequest(t, "/any/v1/whtvr", "POST")
	p.ServeHTTP(writer, req)

	require.Equal(t, http.StatusInternalServerError, writer.StatusCode, "expected internal server error to be returned on proxy error")
	require.Equal(t, "{\"Message\":\"Something went wrong while handling your request. Our team has been notified. Please try again later.\"}", string(writer.Content))
}

func TestProxy_QueryParamsAndFragments_ShouldPropagate(t *testing.T) {
	var proxiedRequest *http.Request
	p := fixtureProxy(t, func(r *http.Request) {
		proxiedRequest = r
	})

	writer := fixtureWriter()
	req := fixtureRequest(t, "/a/foo/bar?id=NTEyMA%3D%3D&step=20#abcd", "GET")
	expectedProxiedUrl := fmt.Sprintf("%s/v1/foo/bar?id=NTEyMA%%3D%%3D&step=20#abcd", systemABaseUrl)

	p.ServeHTTP(writer, req)
	require.Equal(t, expectedProxiedUrl, proxiedRequest.URL.String())
}

func TestProxy_ConcurrentRequests(t *testing.T) {
	p := fixtureProxy(t, func(r *http.Request) {})
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			writer := fixtureWriter()
			req := fixtureRequest(t, "/a/foo/bar", "GET")

			p.ServeHTTP(writer, req)
		}()
	}

	wg.Wait()
}

func fixtureProxy(t *testing.T, onOutgoing func(r *http.Request)) *httputil.ReverseProxy {
	p, err := NewProxy(systemABaseUrl, systemBBaseURL)
	require.NoError(t, err)

	originalDirector := p.Director
	p.Director = func(outgoing *http.Request) {
		onOutgoing(outgoing)
		originalDirector(outgoing)
	}
	p.Transport = noopRoundTripper{onRoundTrip: successRoundTrip}
	return p
}

func fixtureRequest(t *testing.T, path, method string) *http.Request {
	return httptest.NewRequest(method, fmt.Sprintf("%s%s", inputBaseURL, path), nil)
}

func fixtureWriter() *testWriter {
	return &testWriter{headers: make(http.Header)}
}

type testWriter struct {
	StatusCode int
	Content    []byte
	headers    http.Header
}

func (w *testWriter) Header() http.Header {
	return w.headers
}

func (w *testWriter) Write(bytes []byte) (int, error) {
	w.Content = bytes
	return len(bytes), nil
}

func (w *testWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

type noopRoundTripper struct {
	onRoundTrip func(r *http.Request) (*http.Response, error)
}

func (n noopRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return n.onRoundTrip(request)
}

func successRoundTrip(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "OK",
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func errorRoundTrip(request *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("Test round trip error")
}
