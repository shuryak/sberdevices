package yandex

type DeviceInstance string

const (
	DeviceInstanceBrightness   DeviceInstance = "brightness"
	DeviceInstanceChannel      DeviceInstance = "channel"
	DeviceInstanceHumidity     DeviceInstance = "humidity"
	DeviceInstanceOpen         DeviceInstance = "open"
	DeviceInstanceTemperature  DeviceInstance = "temperature"
	DeviceInstanceVolume       DeviceInstance = "volume"
	DeviceInstanceRGB          DeviceInstance = "rgb"
	DeviceInstanceTemperatureK DeviceInstance = "temperature_k"
	DeviceInstanceHSV          DeviceInstance = "hsv"
	DeviceInstanceScene        DeviceInstance = "scene"
	DeviceInstanceOn           DeviceInstance = "on"
)
