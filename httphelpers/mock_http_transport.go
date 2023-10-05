package httphelpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RoundTripFuncType func(req *http.Request) (*http.Response, error)

type MockHTTPTransport struct {
	RoundTripFunc RoundTripFuncType
	RequestCount  int
}

func NewMockHTTPTransport(roundTripFunc RoundTripFuncType) *MockHTTPTransport {
	return &MockHTTPTransport{
		RoundTripFunc: roundTripFunc,
	}
}

func (c *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.RequestCount++
	return c.RoundTripFunc(req)
}

type MockJSONCtx struct {
	Method  string
	URL     string
	Path    string
	HasBody bool
	Body    any
}

func (m *MockJSONCtx) MakeResp(status int, data any) (*http.Response, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBuffer(dataBytes)),
	}, nil
}

func MockJSONEndpoint(theFunc func(c *MockJSONCtx) (*http.Response, error)) func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		ctx := &MockJSONCtx{
			Method:  req.Method,
			URL:     req.URL.String(),
			Path:    req.URL.Path,
			HasBody: req.Body != nil,
		}
		if ctx.HasBody {
			defer req.Body.Close()
			err := json.NewDecoder(req.Body).Decode(&ctx.Body)
			if err != nil {
				return nil, err
			}
		}
		return theFunc(ctx)
	}
}
