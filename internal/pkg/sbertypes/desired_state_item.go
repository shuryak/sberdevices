package sbertypes

import "time"

type StateRequest struct {
	Timestamp    time.Time      `json:"timestamp"`
	DeviceID     string         `json:"device_id"`
	DesiredState []*DeviceState `json:"desired_state"`
}

type StateResponse struct {
	DeviceID     string        `json:"device_id"`
	DesiredState []DeviceState `json:"desired_state"`
}
