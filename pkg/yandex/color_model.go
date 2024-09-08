package yandex

type ColorModel string

const (
	DeviceColorModelHSV ColorModel = "hsv"
	DeviceColorModelRGB ColorModel = "rgb"
)

type DeviceHSVColor struct {
	Hue        int `json:"h"`
	Saturation int `json:"s"`
	Value      int `json:"v"`
}
