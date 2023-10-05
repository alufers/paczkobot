package inpostextra_test

import (
	"net/http"
	"testing"

	"github.com/alufers/paczkobot/httphelpers"
	"github.com/alufers/paczkobot/inpostextra"
	"github.com/stretchr/testify/assert"
)

func TestInpostServiceSendSMSCode(t *testing.T) {
	//nolint:bodyclose
	mockTransport := httphelpers.NewMockHTTPTransport(httphelpers.MockJSONEndpoint(func(c *httphelpers.MockJSONCtx) (*http.Response, error) {
		assert.Equal(t, "POST", c.Method)
		assert.Equal(t, "/v1/sendSMSCode", c.Path)
		assert.Equal(t, true, c.HasBody)
		assert.Equal(t, map[string]interface{}{
			"phoneNumber": "123456789",
		}, c.Body)
		return c.MakeResp(200, map[string]interface{}{
			"status": "OK",
		})
	}))
	serv := inpostextra.NewInpostService(&http.Client{
		Transport: mockTransport,
	})
	err := serv.SendSMSCode("123456789")
	assert.NoError(t, err)
	assert.Equal(t, 1, mockTransport.RequestCount)
}
