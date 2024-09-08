package yandex

type DeviceCapabilityType string

const (
	DeviceCapabilityTypeOnOff        DeviceCapabilityType = "devices.capabilities.on_off"
	DeviceCapabilityTypeColorSetting DeviceCapabilityType = "devices.capabilities.color_setting"
	DeviceCapabilityTypeVideoStream  DeviceCapabilityType = "devices.capabilities.video_stream"
	DeviceCapabilityTypeMode         DeviceCapabilityType = "devices.capabilities.mode"
	DeviceCapabilityTypeRange        DeviceCapabilityType = "devices.capabilities.range"
	DeviceCapabilityTypeToggle       DeviceCapabilityType = "devices.capabilities.toggle"
)
