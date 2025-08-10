package sbertypes

type SberDataType string

const (
	SberDataTypeBool    SberDataType = "BOOL"
	SberDataTypeEnum    SberDataType = "ENUM"
	SberDataTypeInteger SberDataType = "INTEGER"
	SberDataTypeColor   SberDataType = "COLOR"
)

type StateCommand string

const (
	StateCommandSwitchLED       = "switch_led"
	StateCommandOnline          = "online"
	StateCommandOnOff           = "on_off"
	StateCommandLightBrightness = "light_brightness"
	StateCommandLightColourTemp = "light_colour_temp"
	StateCommandLightScene      = "light_scene"
	StateCommandLightMode       = "light_mode"
	StateCommandLightColour     = "light_colour"
)

type LightMode string

const (
	LightModeWhite  = "white"
	LightModeColour = "colour"
	LightModeScene  = "scene"
)
