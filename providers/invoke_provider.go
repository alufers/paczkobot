package providers

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sort"

	"github.com/alufers/paczkobot/commondata"
)

func InvokeProvider(ctx context.Context, provider Provider, trackingNumber string) (result *commondata.TrackingData, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic at %v provider: %v\nstacktrace: %v", provider.GetName(), r, string(debug.Stack()))
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	result, err = provider.Track(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("provider returned nil result")
	}
	sort.SliceStable(result.TrackingSteps, func(i, j int) bool {
		return result.TrackingSteps[i].Datetime.UnixNano() < result.TrackingSteps[j].Datetime.UnixNano()
	})
	return result, nil
}
