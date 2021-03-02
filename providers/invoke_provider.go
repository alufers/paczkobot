package providers

import (
	"context"
	"fmt"
	"github.com/alufers/paczkobot/commondata"
	"sort"
)

func InvokeProvider(ctx context.Context, provider Provider, trackingNumber string) (result *commondata.TrackingData, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	result, err = provider.Track(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(result.TrackingSteps, func(i, j int) bool {
		return result.TrackingSteps[i].Datetime.UnixNano() < result.TrackingSteps[j].Datetime.UnixNano()
	})
	return result, nil
}
