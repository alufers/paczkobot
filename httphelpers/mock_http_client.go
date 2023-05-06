package httphelpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type MockHTTPClient struct {
	DoFunc       func(req *http.Request) (*http.Response, error)
	RequestCount int
}

func NewMockHTTPClient(doFunc func(req *http.Request) (*http.Response, error)) *MockHTTPClient {
	return &MockHTTPClient{
		DoFunc: doFunc,
	}
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.RequestCount++
	return c.DoFunc(req)
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
