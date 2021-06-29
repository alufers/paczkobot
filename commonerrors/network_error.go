package commonerrors

import (
	"fmt"
	"net/http"
)

var NetworkError = fmt.Errorf("not the real network error do not use")

type networkErrorWithText struct {
	req          *http.Request
	providerName string
}

func (n *networkErrorWithText) Error() string {
	return fmt.Sprintf("Network error when sending a request to %v (%v %v)", n.providerName, n.req.Method, n.req.URL.String())
}

func (n *networkErrorWithText) Is(target error) bool {
	return target == NetworkError
}

func NewNetworkError(providerName string, req *http.Request) error {
	return &networkErrorWithText{
		providerName: providerName,
		req:          req,
	}
}
