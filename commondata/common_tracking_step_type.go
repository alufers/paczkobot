package commondata

type CommonTrackingStepType int64

var (
	CommonTrackingStepType_UNKNOWN              = CommonTrackingStepType(0)
	CommonTrackingStepType_INFORMATION_PREPARED = CommonTrackingStepType(1)
	CommonTrackingStepType_SENT                 = CommonTrackingStepType(2)
	CommonTrackingStepType_IN_TRANSIT           = CommonTrackingStepType(3)
	CommonTrackingStepType_OUT_FOR_DELIVERY     = CommonTrackingStepType(4)
	CommonTrackingStepType_READY_FOR_PICKUP     = CommonTrackingStepType(5)
	CommonTrackingStepType_DELIVERED            = CommonTrackingStepType(6)
	CommonTrackingStepType_FAILURE              = CommonTrackingStepType(7)
)

var CommonTrackingStepTypeName = map[CommonTrackingStepType]string{
	0: "INFORMATION_PREPARED",
	1: "SENT",
	2: "IN_TRANSIT",
	3: "OUT_FOR_DELIVERY",
	4: "READY_FOR_PICKUP",
	5: "DELIVERED",
	6: "FAILURE",
}

var CommonTrackingStepTypeEmoji = map[CommonTrackingStepType]string{
	1: "📝",
	2: "✉️",
	3: "🚚",
	4: "🚶",
	5: "📦",
	6: "✅",
	7: "❗",
}
