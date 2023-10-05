package providerutil

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/httphelpers"
)

func FetchGoqueryDocument(ctx context.Context, providerName string, url string, checkStatusCode bool) (*goquery.Document, error) {
	client := httphelpers.NewClientWithLogger()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	commondata.SetCommonHTTPHeaders(&req.Header)

	httpResponse, err := client.Do(req)
	if err != nil {
		// Assuming NewNetworkError is available in your package
		return nil, commonerrors.NewNetworkError(providerName, req)
	}
	defer httpResponse.Body.Close()

	if checkStatusCode && httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code %v", httpResponse.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML response: %w", err)
	}

	return doc, nil
}
