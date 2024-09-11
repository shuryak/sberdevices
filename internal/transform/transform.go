package transform

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/shuryak/sberhack/pkg/sbertypes"
	"github.com/shuryak/sberhack/pkg/yandex"
)

func SberToYandexDevices(sberDevices []sbertypes.DeviceItem) []yandex.Device {
	var yandexDevices []yandex.Device

	for _, device := range sberDevices {
		yandexDevices = append(yandexDevices, *SberToYandexDeviceInfo(&device))
	}

	return yandexDevices
}

func SberToYandexDevicesState(sberDevices []sbertypes.DeviceItem) []yandex.Device {
	var yandexDevices []yandex.Device

	for _, device := range sberDevices {
		yandexDevices = append(yandexDevices, *SberToYandexDeviceStates(&device))
	}

	return yandexDevices
}

func SberToYandexDeviceStates(sberDevice *sbertypes.DeviceItem) *yandex.Device {
	yandexDevice := &yandex.Device{}

	if sberDevice.Name != nil {
		yandexDevice.ID = sberDevice.ID
	}

	sberCommands, _, _ := getCommandsInfo(sberDevice)

	reportedStates := make(map[sbertypes.StateCommand]*sbertypes.DeviceState)
	for i := 0; i < len(sberDevice.ReportedState); i++ {
		reportedStates[sberDevice.ReportedState[i].Key] = &sberDevice.ReportedState[i]
	}

	sberCommandsKeys := make([]sbertypes.StateCommand, 0, len(sberCommands))
	for k := range sberCommands {
		sberCommandsKeys = append(sberCommandsKeys, k)
	}
	slices.Sort(sberCommandsKeys)

	for _, k := range sberCommandsKeys {
		yandexCapabilityType, ok := stateCommandToCapabilityMap[k]
		if !ok {
			continue
		}

		state := sberToYandexDeviceCapabilityState(k, reportedStates)

		if state != nil {
			yandexDevice.Capabilities = append(yandexDevice.Capabilities, yandex.DeviceCapability{
				Type:  yandexCapabilityType,
				State: state,
			})
		}
	}

	return yandexDevice
}

func sberToYandexDeviceCapabilityState(
	sberCommand sbertypes.StateCommand,
	sberReportedStates map[sbertypes.StateCommand]*sbertypes.DeviceState,
) *yandex.DeviceCapabilityState {
	yandexState := &yandex.DeviceCapabilityState{
		Instance: stateCommandToInstanceMap[sberCommand],
	}

	reportedState := sberReportedStates[sberCommand]

	switch yandexState.Instance {
	case yandex.DeviceInstanceOn:
		yandexState.Value = reportedState.BoolValue
	case yandex.DeviceInstanceBrightness:
		value, _ := strconv.Atoi(reportedState.IntegerValue) // TODO: handle errors for atoi everywhere
		value /= 10
		yandexState.Value = value
	case yandex.DeviceInstanceTemperatureK:
		value, _ := strconv.Atoi(reportedState.IntegerValue)
		value = 7*value + 2000 // normalize [0, 1000] to [2000, 9000]
		yandexState.Value = value
	case yandex.DeviceInstanceScene:
		var ok bool
		yandexState.Value, ok = sberColorSceneIDToYandexMap[sbertypes.ColorSceneID(reportedState.EnumValue)] // TODO: handle ""
		if !ok {
			return nil
		}
	case yandex.DeviceInstanceHSV:
		yandexState.Value = yandex.DeviceHSVColor{
			Hue:        reportedState.ColorValue.Hue,
			Saturation: reportedState.ColorValue.Saturation / 10, // TODO: ?
			Value:      reportedState.ColorValue.Value / 10,      // TODO: ?
		}
	default:
		return nil
	}

	return yandexState
}

func SberToYandexDeviceInfo(sberDevice *sbertypes.DeviceItem) *yandex.Device {
	yandexDevice := &yandex.Device{}

	if sberDevice.Name != nil {
		yandexDevice.ID = sberDevice.ID
		yandexDevice.Name = sberDevice.Name.Name
	}

	yandexDevice.Description = sberDevice.DeviceTypeName

	if sberDevice.DeviceInfo != nil {
		yandexDevice.DeviceInfo = &yandex.DeviceInfo{
			Manufacturer: sberDevice.DeviceInfo.Manufacturer,
			Model:        sberDevice.DeviceInfo.Model,
			HWVersion:    sberDevice.DeviceInfo.HWVersion,
			SWVersion: fmt.Sprintf(
				"%s (%s)",
				sberDevice.SWVersion,
				sberDevice.DeviceInfo.SWVersionInt,
			),
		}
	}

	var sberCommands, allStateFields map[sbertypes.StateCommand]struct{}
	sberCommands, allStateFields, yandexDevice.Type = getCommandsInfo(sberDevice)

	capabilitiesMap := make(map[yandex.DeviceCapabilityType]*yandex.DeviceCapability)
	parametersMap := make(map[yandex.DeviceCapabilityType]*yandex.DeviceCapabilitiesParameters)

	for _, attribute := range sberDevice.Attributes {
		if _, ok := sberCommands[attribute.Key]; !ok {
			continue
		}

		if attribute.Key == sbertypes.StateCommandLightMode { // TODO: temp
			continue
		}

		yandexCapabilityType, ok := stateCommandToCapabilityMap[attribute.Key]
		if !ok {
			continue
		}

		_, retrievable := allStateFields[attribute.Key]

		capabilitiesMap[yandexCapabilityType] = &yandex.DeviceCapability{
			Type:        yandexCapabilityType,
			Retrievable: &retrievable,
			Reportable:  false,
		}

		parameters := makeYandexCapabilitiesParameters(&attribute)

		if v, ok := parametersMap[yandexCapabilityType]; !ok {
			parametersMap[yandexCapabilityType] = parameters
		} else if v != nil && parameters != nil {
			if parameters.Split != nil {
				parametersMap[yandexCapabilityType].Split = parameters.Split
			}
			if len(parameters.Instance) != 0 {
				parametersMap[yandexCapabilityType].Instance = parameters.Instance
			}
			if len(parameters.Unit) != 0 {
				parametersMap[yandexCapabilityType].Unit = parameters.Unit
			}
			if parameters.RandomAccess != nil {
				parametersMap[yandexCapabilityType].RandomAccess = parameters.RandomAccess
			}
			if parameters.Range != nil {
				parametersMap[yandexCapabilityType].Range = parameters.Range
			}
			if len(parameters.ColorModel) != 0 {
				parametersMap[yandexCapabilityType].ColorModel = parameters.ColorModel
			}
			if parameters.TemperatureK != nil {
				parametersMap[yandexCapabilityType].TemperatureK = parameters.TemperatureK
			}
			if parameters.ColorScene != nil {
				parametersMap[yandexCapabilityType].ColorScene = parameters.ColorScene
			}
		}
	}

	capabilitiesKeys := make([]yandex.DeviceCapabilityType, 0, len(capabilitiesMap))
	for k := range capabilitiesMap {
		capabilitiesKeys = append(capabilitiesKeys, k)
	}
	slices.Sort(capabilitiesKeys)

	for _, k := range capabilitiesKeys {
		if capabilitiesMap[k] == nil {
			continue
		}

		capabilitiesMap[k].Parameters = parametersMap[k]
		yandexDevice.Capabilities = append(yandexDevice.Capabilities, *capabilitiesMap[k])
	}

	return yandexDevice
}

func YandexToSberDeviceState(
	currentState map[sbertypes.StateCommand]sbertypes.DeviceState,
	yandexCapability *yandex.DeviceCapability,
) []*sbertypes.DeviceState {
	var states []*sbertypes.DeviceState

	key := instanceToStateCommandMap[yandexCapability.State.Instance]
	now := time.Now()

	switch yandexCapability.State.Instance {
	case yandex.DeviceInstanceOn:
		states = append(states, &sbertypes.DeviceState{
			Type:      sbertypes.SberDataTypeBool,
			BoolValue: yandexCapability.State.Value.(bool),
		})
	case yandex.DeviceInstanceBrightness:
		value := int(yandexCapability.State.Value.(float64)) * 10

		states = append(states,
			&sbertypes.DeviceState{
				Key:          sbertypes.StateCommandLightBrightness,
				Type:         sbertypes.SberDataTypeInteger,
				IntegerValue: strconv.Itoa(value),
			},
		)
	case yandex.DeviceInstanceTemperatureK:
		value := int(yandexCapability.State.Value.(float64))
		value = ((value - 2000) * 1000) / 9000

		states = append(states,
			&sbertypes.DeviceState{
				Key:       sbertypes.StateCommandLightMode,
				Type:      sbertypes.SberDataTypeEnum,
				EnumValue: sbertypes.LightModeWhite,
			},
			&sbertypes.DeviceState{
				Type:         sbertypes.SberDataTypeInteger,
				IntegerValue: strconv.Itoa(value),
			},
		)
	case yandex.DeviceInstanceScene:
		states = append(states, &sbertypes.DeviceState{
			Type: sbertypes.SberDataTypeEnum,
			IntegerValue: string(
				yandexMapToSberColorSceneID[yandex.ColorSceneID(yandexCapability.State.Value.(string))],
			),
		})
	case yandex.DeviceInstanceHSV:
		value := yandexCapability.State.Value.(map[string]interface{})

		states = append(states,
			&sbertypes.DeviceState{
				Key:       sbertypes.StateCommandLightMode,
				Type:      sbertypes.SberDataTypeEnum,
				EnumValue: sbertypes.LightModeColour,
			},
			&sbertypes.DeviceState{
				Type: sbertypes.SberDataTypeColor,
				ColorValue: &sbertypes.DeviceStateColorValue{
					Hue:        int(value["h"].(float64)),
					Saturation: int(value["s"].(float64)) * 10,
					Value:      int(value["v"].(float64)) * 10,
				},
			},
		)
	}

	for i := range states {
		if len(states[i].Key) == 0 {
			states[i].Key = key
		}
		states[i].LastSync = now
	}

	return states
}

func makeYandexCapabilitiesParameters(sberAttribute *sbertypes.DeviceAttribute) *yandex.DeviceCapabilitiesParameters {
	switch sberAttribute.Key {
	case sbertypes.StateCommandOnOff, sbertypes.StateCommandSwitchLED:
		return &yandex.DeviceCapabilitiesParameters{
			Split: nilableFalse,
		}
	case sbertypes.StateCommandLightBrightness:
		if sberAttribute.IntValues == nil {
			return nil
		}

		return &yandex.DeviceCapabilitiesParameters{
			Instance:     yandex.DeviceInstanceBrightness,
			Unit:         yandex.UnitPercent,
			RandomAccess: nilableTrue,
			Range: &yandex.DeviceCapabilitiesParametersRange{
				// TODO: range min_max * 0.1 everywhere, problem with min = 5%. For yandex.UnitPercent
				Min:       0,
				Max:       float64(sberAttribute.IntValues.Range.Max / 10),
				Precision: float64(sberAttribute.IntValues.Range.Step),
			},
		}
	case sbertypes.StateCommandLightColourTemp:
		if sberAttribute.IntValues == nil {
			return nil
		}

		return &yandex.DeviceCapabilitiesParameters{
			Instance: yandex.DeviceInstanceTemperatureK,
			TemperatureK: &yandex.DeviceCapabilitiesParametersRange{
				// TODO: normalize int_values.range.min and int_values.range.max to [2000, 9000]
				Min:       2000,
				Max:       9000,
				Precision: float64(sberAttribute.IntValues.Range.Step),
			},
		}
	case sbertypes.StateCommandLightScene:
		if sberAttribute.EnumValues == nil {
			return nil
		}

		var scenes []yandex.DeviceColorSceneItem
		for _, scene := range sberAttribute.EnumValues.Values {
			scenes = append(scenes, yandex.DeviceColorSceneItem{
				ID: sberColorSceneIDToYandexMap[sbertypes.ColorSceneID(scene)],
			})
		}

		return &yandex.DeviceCapabilitiesParameters{
			Instance: yandex.DeviceInstanceScene,
			ColorScene: &yandex.DeviceColorScene{
				Scenes: scenes,
			},
		}
	case sbertypes.StateCommandLightMode:
		// TODO: light_mode
		return nil
	case sbertypes.StateCommandLightColour:
		return &yandex.DeviceCapabilitiesParameters{
			ColorModel: yandex.DeviceColorModelHSV,
		}
	}

	return nil
}

func getCommandsInfo(sberDevice *sbertypes.DeviceItem) (
	commands map[sbertypes.StateCommand]struct{},
	allStateFields map[sbertypes.StateCommand]struct{},
	yandexDeviceType yandex.DeviceType,
) {
	commands = make(map[sbertypes.StateCommand]struct{})
	allStateFields = make(map[sbertypes.StateCommand]struct{})

	for _, command := range sberDevice.Commands {
		commands[command.Key] = struct{}{}
		for _, stateField := range command.StateFields {
			allStateFields[stateField] = struct{}{}
		}
	}

	for _, category := range sberDevice.FullCategories {
		if category.Slug == "light" || category.Slug == "led_strip" { // TODO: constants for slugs
			yandexDeviceType = yandex.DeviceTypeLightStrip

			// switch_led and on_off have the same effect
			if _, ok := commands[sbertypes.StateCommandSwitchLED]; ok {
				delete(commands, sbertypes.StateCommandOnOff)
			}
		}
	}

	return
}

var stateCommandToCapabilityMap = map[sbertypes.StateCommand]yandex.DeviceCapabilityType{
	sbertypes.StateCommandOnOff:           yandex.DeviceCapabilityTypeOnOff,
	sbertypes.StateCommandSwitchLED:       yandex.DeviceCapabilityTypeOnOff,
	sbertypes.StateCommandLightBrightness: yandex.DeviceCapabilityTypeRange,
	sbertypes.StateCommandLightColourTemp: yandex.DeviceCapabilityTypeColorSetting,
	sbertypes.StateCommandLightScene:      yandex.DeviceCapabilityTypeColorSetting,
	sbertypes.StateCommandLightMode:       yandex.DeviceCapabilityTypeMode, // TODO: ?
	sbertypes.StateCommandLightColour:     yandex.DeviceCapabilityTypeColorSetting,
}

var stateCommandToInstanceMap = map[sbertypes.StateCommand]yandex.DeviceInstance{
	sbertypes.StateCommandOnOff:           yandex.DeviceInstanceOn,
	sbertypes.StateCommandSwitchLED:       yandex.DeviceInstanceOn,
	sbertypes.StateCommandLightBrightness: yandex.DeviceInstanceBrightness,
	sbertypes.StateCommandLightColourTemp: yandex.DeviceInstanceTemperatureK,
	sbertypes.StateCommandLightScene:      yandex.DeviceInstanceScene,
	sbertypes.StateCommandLightMode:       "", // TODO: ?
	sbertypes.StateCommandLightColour:     yandex.DeviceInstanceHSV,
}

var instanceToStateCommandMap = map[yandex.DeviceInstance]sbertypes.StateCommand{
	yandex.DeviceInstanceOn:           sbertypes.StateCommandOnOff,
	yandex.DeviceInstanceBrightness:   sbertypes.StateCommandLightBrightness,
	yandex.DeviceInstanceTemperatureK: sbertypes.StateCommandLightColourTemp,
	yandex.DeviceInstanceScene:        sbertypes.StateCommandLightScene,
	yandex.DeviceInstanceHSV:          sbertypes.StateCommandLightColour,
}

var sberColorSceneIDToYandexMap = map[sbertypes.ColorSceneID]yandex.ColorSceneID{
	sbertypes.ColorSceneIDCandle:    yandex.ColorSceneIDCandle,
	sbertypes.ColorSceneIDArctic:    yandex.ColorSceneIDOcean,
	sbertypes.ColorSceneIDRomantic:  yandex.ColorSceneIDRomance,
	sbertypes.ColorSceneIDSunset:    yandex.ColorSceneIDSunset,
	sbertypes.ColorSceneIDDawn:      yandex.ColorSceneIDSunrise,
	sbertypes.ColorSceneIDChristmas: yandex.ColorSceneIDGarland,
	sbertypes.ColorSceneIDFito:      yandex.ColorSceneIDRest,
}

var yandexMapToSberColorSceneID = map[yandex.ColorSceneID]sbertypes.ColorSceneID{
	yandex.ColorSceneIDCandle:  sbertypes.ColorSceneIDCandle,
	yandex.ColorSceneIDOcean:   sbertypes.ColorSceneIDArctic,
	yandex.ColorSceneIDRomance: sbertypes.ColorSceneIDRomantic,
	yandex.ColorSceneIDSunset:  sbertypes.ColorSceneIDSunset,
	yandex.ColorSceneIDSunrise: sbertypes.ColorSceneIDDawn,
	yandex.ColorSceneIDGarland: sbertypes.ColorSceneIDChristmas,
	yandex.ColorSceneIDRest:    sbertypes.ColorSceneIDFito,
}

func nilableBool(v bool) *bool {
	return &v
}

var (
	nilableFalse = nilableBool(false)
	nilableTrue  = nilableBool(true)
)
