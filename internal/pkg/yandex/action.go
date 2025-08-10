package yandex

type DevicesActionRequest struct {
	Payload struct {
		Devices []Device `json:"devices"`
	} `json:"payload"`
}
