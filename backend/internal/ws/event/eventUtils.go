package event

import (
	"encoding/json"
	"fmt"
)

func CreateEvent(eventType string, payload interface{}) (Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Event{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return Event{Type: eventType, Payload: data}, nil
}