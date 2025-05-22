package message

type BMTPublicOutboxesMsg struct {
	After AfterPayload `json:"after"`
}

type AfterPayload struct {
	ID             string `json:"id"`
	AggregatedType string `json:"aggregated_type"`
	AggregatedID   int32  `json:"aggregated_id"`
	EventType      string `json:"event_type"`
	Payload        string `json:"payload"`
	CreatedAt      int64  `json:"created_at"`
}
