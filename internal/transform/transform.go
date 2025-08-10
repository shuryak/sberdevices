package transform

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	sbertypes2 "github.com/shuryak/sberdevices/internal/pkg/sbertypes"
	yandex2 "github.com/shuryak/sberdevices/internal/pkg/yandex"
)

func SberToYandexDevices(sberDevices []sbertypes2.DeviceItem) []yandex2.Device {
	var yandexDevices []yandex2.Device

	for _, device := range sberDevices {
		yandexDevices = append(yandexDevices, *SberToYandexDeviceInfo(&device))
	}

	return yandexDevices
}

func SberToYandexDevicesState(sberDevices []sbertypes2.DeviceItem) []yandex2.Device {
	var yandexDevices []yandex2.Device

	for _, device := range sberDevices {
		yandexDevices = append(yandexDevices, *SberToYandexDeviceStates(&device))
	}

	return yandexDevices
}

func SberToYandexDeviceStates(sberDevice *sbertypes2.DeviceItem) *yandex2.Device {
	yandexDevice := &yandex2.Device{}

	if sberDevice.Name != nil {
		yandexDevice.ID = sberDevice.ID
	}

	sberCommands, _, _ := getCommandsInfo(sberDevice)

	reportedStates := make(map[sbertypes2.StateCommand]*sbertypes2.DeviceState)
	for i := 0; i < len(sberDevice.ReportedState); i++ {
		reportedStates[sberDevice.ReportedState[i].Key] = &sberDevice.ReportedState[i]
	}

	sberCommandsKeys := make([]sbertypes2.StateCommand, 0, len(sberCommands))
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
			yandexDevice.Capabilities = append(yandexDevice.Capabilities, yandex2.DeviceCapability{
				Type:  yandexCapabilityType,
				State: state,
			})
		}
	}

	return yandexDevice
}

func sberToYandexDeviceCapabilityState(
	sberCommand sbertypes2.StateCommand,
	sberReportedStates map[sbertypes2.StateCommand]*sbertypes2.DeviceState,
) *yandex2.DeviceCapabilityState {
	yandexState := &yandex2.DeviceCapabilityState{
		Instance: stateCommandToInstanceMap[sberCommand],
	}

	reportedState := sberReportedStates[sberCommand]

	switch yandexState.Instance {
	case yandex2.DeviceInstanceOn:
		yandexState.Value = reportedState.BoolValue
	case yandex2.DeviceInstanceBrightness:
		value, _ := strconv.Atoi(reportedState.IntegerValue) // TODO: handle errors for atoi everywhere
		value /= 10
		yandexState.Value = value
	case yandex2.DeviceInstanceTemperatureK:
		value, _ := strconv.Atoi(reportedState.IntegerValue)
		value = 7*value + 2000 // normalize [0, 1000] to [2000, 9000]
		yandexState.Value = value
	case yandex2.DeviceInstanceScene:
		var ok bool
		yandexState.Value, ok = sberColorSceneIDToYandexMap[sbertypes2.ColorSceneID(reportedState.EnumValue)] // TODO: handle ""
		if !ok {
			return nil
		}
	case yandex2.DeviceInstanceHSV:
		yandexState.Value = yandex2.DeviceHSVColor{
			Hue:        reportedState.ColorValue.Hue,
			Saturation: reportedState.ColorValue.Saturation / 10, // TODO: ?
			Value:      reportedState.ColorValue.Value / 10,      // TODO: ?
		}
	default:
		return nil
	}

	return yandexState
}

func SberToYandexDeviceInfo(sberDevice *sbertypes2.DeviceItem) *yandex2.Device {
	yandexDevice := &yandex2.Device{}

	if sberDevice.Name != nil {
		yandexDevice.ID = sberDevice.ID
		yandexDevice.Name = sberDevice.Name.Name
	}

	yandexDevice.Description = sberDevice.DeviceTypeName

	if sberDevice.DeviceInfo != nil {
		yandexDevice.DeviceInfo = &yandex2.DeviceInfo{
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

	var sberCommands, allStateFields map[sbertypes2.StateCommand]struct{}
	sberCommands, allStateFields, yandexDevice.Type = getCommandsInfo(sberDevice)

	capabilitiesMap := make(map[yandex2.DeviceCapabilityType]*yandex2.DeviceCapability)
	parametersMap := make(map[yandex2.DeviceCapabilityType]*yandex2.DeviceCapabilitiesParameters)

	for _, attribute := range sberDevice.Attributes {
		if _, ok := sberCommands[attribute.Key]; !ok {
			continue
		}

		if attribute.Key == sbertypes2.StateCommandLightMode { // TODO: temp
			continue
		}

		yandexCapabilityType, ok := stateCommandToCapabilityMap[attribute.Key]
		if !ok {
			continue
		}

		_, retrievable := allStateFields[attribute.Key]

		capabilitiesMap[yandexCapabilityType] = &yandex2.DeviceCapability{
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

	capabilitiesKeys := make([]yandex2.DeviceCapabilityType, 0, len(capabilitiesMap))
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
	currentState map[sbertypes2.StateCommand]sbertypes2.DeviceState,
	yandexCapability *yandex2.DeviceCapability,
) []*sbertypes2.DeviceState {
	var states []*sbertypes2.DeviceState

	key := instanceToStateCommandMap[yandexCapability.State.Instance]
	now := time.Now()

	switch yandexCapability.State.Instance {
	case yandex2.DeviceInstanceOn:
		states = append(states, &sbertypes2.DeviceState{
			Type:      sbertypes2.SberDataTypeBool,
			BoolValue: yandexCapability.State.Value.(bool),
		})
	case yandex2.DeviceInstanceBrightness:
		value := int(yandexCapability.State.Value.(float64)) * 10

		cur := currentState[sbertypes2.StateCommandLightColour]

		states = append(states,
			&sbertypes2.DeviceState{
				Key:  cur.Key,
				Type: sbertypes2.SberDataTypeColor,
				ColorValue: &sbertypes2.DeviceStateColorValue{
					Hue:        cur.ColorValue.Hue,
					Saturation: cur.ColorValue.Saturation,
					Value:      value,
				},
			},
			&sbertypes2.DeviceState{
				Key:          sbertypes2.StateCommandLightBrightness,
				Type:         sbertypes2.SberDataTypeInteger,
				IntegerValue: strconv.Itoa(value),
			},
		)
	case yandex2.DeviceInstanceTemperatureK:
		value := int(yandexCapability.State.Value.(float64))
		value = ((value - 2000) * 1000) / 7000

		states = append(states,
			&sbertypes2.DeviceState{
				Key:       sbertypes2.StateCommandLightMode,
				Type:      sbertypes2.SberDataTypeEnum,
				EnumValue: sbertypes2.LightModeWhite,
			},
			&sbertypes2.DeviceState{
				Type:         sbertypes2.SberDataTypeInteger,
				IntegerValue: strconv.Itoa(value),
			},
		)
	case yandex2.DeviceInstanceScene:
		states = append(states, &sbertypes2.DeviceState{
			Type: sbertypes2.SberDataTypeEnum,
			IntegerValue: string(
				yandexMapToSberColorSceneID[yandex2.ColorSceneID(yandexCapability.State.Value.(string))],
			),
		})
	case yandex2.DeviceInstanceHSV:
		value := make(map[string]interface{})
		if yandexCapability.State.Value == nil { // for Marusia command "turn on the black color"
			value["h"] = float64(0)
			value["s"] = float64(0)
			value["v"] = float64(0)
		} else {
			// TODO: log type of value (check for ok on type assertion)
			value = yandexCapability.State.Value.(map[string]interface{})
		}

		cur := currentState[sbertypes2.StateCommandLightColour]

		states = append(states,
			&sbertypes2.DeviceState{
				Key:       sbertypes2.StateCommandLightMode,
				Type:      sbertypes2.SberDataTypeEnum,
				EnumValue: sbertypes2.LightModeColour,
			},
			&sbertypes2.DeviceState{
				Type: sbertypes2.SberDataTypeColor,
				ColorValue: &sbertypes2.DeviceStateColorValue{
					Hue:        int(value["h"].(float64)),
					Saturation: int(value["s"].(float64)) * 10,
					Value:      cur.ColorValue.Value,
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

func makeYandexCapabilitiesParameters(sberAttribute *sbertypes2.DeviceAttribute) *yandex2.DeviceCapabilitiesParameters {
	switch sberAttribute.Key {
	case sbertypes2.StateCommandOnOff, sbertypes2.StateCommandSwitchLED:
		return &yandex2.DeviceCapabilitiesParameters{
			Split: nilableFalse,
		}
	case sbertypes2.StateCommandLightBrightness:
		if sberAttribute.IntValues == nil {
			return nil
		}

		return &yandex2.DeviceCapabilitiesParameters{
			Instance:     yandex2.DeviceInstanceBrightness,
			Unit:         yandex2.UnitPercent,
			RandomAccess: nilableTrue,
			Range: &yandex2.DeviceCapabilitiesParametersRange{
				// TODO: range min_max * 0.1 everywhere, problem with min = 5%. For yandex.UnitPercent
				Min:       0,
				Max:       float64(sberAttribute.IntValues.Range.Max / 10),
				Precision: float64(sberAttribute.IntValues.Range.Step),
			},
		}
	case sbertypes2.StateCommandLightColourTemp:
		if sberAttribute.IntValues == nil {
			return nil
		}

		return &yandex2.DeviceCapabilitiesParameters{
			Instance: yandex2.DeviceInstanceTemperatureK,
			TemperatureK: &yandex2.DeviceCapabilitiesParametersRange{
				// TODO: normalize int_values.range.min and int_values.range.max to [2000, 9000]
				Min:       2000,
				Max:       9000,
				Precision: float64(sberAttribute.IntValues.Range.Step),
			},
		}
	case sbertypes2.StateCommandLightScene:
		if sberAttribute.EnumValues == nil {
			return nil
		}

		var scenes []yandex2.DeviceColorSceneItem
		for _, scene := range sberAttribute.EnumValues.Values {
			scenes = append(scenes, yandex2.DeviceColorSceneItem{
				ID: sberColorSceneIDToYandexMap[sbertypes2.ColorSceneID(scene)],
			})
		}

		return &yandex2.DeviceCapabilitiesParameters{
			Instance: yandex2.DeviceInstanceScene,
			ColorScene: &yandex2.DeviceColorScene{
				Scenes: scenes,
			},
		}
	case sbertypes2.StateCommandLightMode:
		// TODO: light_mode
		return nil
	case sbertypes2.StateCommandLightColour:
		return &yandex2.DeviceCapabilitiesParameters{
			ColorModel: yandex2.DeviceColorModelHSV,
		}
	}

	return nil
}

func getCommandsInfo(sberDevice *sbertypes2.DeviceItem) (
	commands map[sbertypes2.StateCommand]struct{},
	allStateFields map[sbertypes2.StateCommand]struct{},
	yandexDeviceType yandex2.DeviceType,
) {
	commands = make(map[sbertypes2.StateCommand]struct{})
	allStateFields = make(map[sbertypes2.StateCommand]struct{})

	for _, command := range sberDevice.Commands {
		commands[command.Key] = struct{}{}
		for _, stateField := range command.StateFields {
			allStateFields[stateField] = struct{}{}
		}
	}

	for _, category := range sberDevice.FullCategories {
		if category.Slug == "light" || category.Slug == "led_strip" { // TODO: constants for slugs
			yandexDeviceType = yandex2.DeviceTypeLightStrip

			// switch_led and on_off have the same effect
			if _, ok := commands[sbertypes2.StateCommandSwitchLED]; ok {
				delete(commands, sbertypes2.StateCommandOnOff)
			}
		}
	}

	return
}

var stateCommandToCapabilityMap = map[sbertypes2.StateCommand]yandex2.DeviceCapabilityType{
	sbertypes2.StateCommandOnOff:           yandex2.DeviceCapabilityTypeOnOff,
	sbertypes2.StateCommandSwitchLED:       yandex2.DeviceCapabilityTypeOnOff,
	sbertypes2.StateCommandLightBrightness: yandex2.DeviceCapabilityTypeRange,
	sbertypes2.StateCommandLightColourTemp: yandex2.DeviceCapabilityTypeColorSetting,
	sbertypes2.StateCommandLightScene:      yandex2.DeviceCapabilityTypeColorSetting,
	sbertypes2.StateCommandLightMode:       yandex2.DeviceCapabilityTypeMode, // TODO: ?
	sbertypes2.StateCommandLightColour:     yandex2.DeviceCapabilityTypeColorSetting,
}

var stateCommandToInstanceMap = map[sbertypes2.StateCommand]yandex2.DeviceInstance{
	sbertypes2.StateCommandOnOff:           yandex2.DeviceInstanceOn,
	sbertypes2.StateCommandSwitchLED:       yandex2.DeviceInstanceOn,
	sbertypes2.StateCommandLightBrightness: yandex2.DeviceInstanceBrightness,
	sbertypes2.StateCommandLightColourTemp: yandex2.DeviceInstanceTemperatureK,
	sbertypes2.StateCommandLightScene:      yandex2.DeviceInstanceScene,
	sbertypes2.StateCommandLightMode:       "", // TODO: ?
	sbertypes2.StateCommandLightColour:     yandex2.DeviceInstanceHSV,
}

var instanceToStateCommandMap = map[yandex2.DeviceInstance]sbertypes2.StateCommand{
	yandex2.DeviceInstanceOn:           sbertypes2.StateCommandOnOff,
	yandex2.DeviceInstanceBrightness:   sbertypes2.StateCommandLightBrightness,
	yandex2.DeviceInstanceTemperatureK: sbertypes2.StateCommandLightColourTemp,
	yandex2.DeviceInstanceScene:        sbertypes2.StateCommandLightScene,
	yandex2.DeviceInstanceHSV:          sbertypes2.StateCommandLightColour,
}

var sberColorSceneIDToYandexMap = map[sbertypes2.ColorSceneID]yandex2.ColorSceneID{
	sbertypes2.ColorSceneIDCandle:    yandex2.ColorSceneIDCandle,
	sbertypes2.ColorSceneIDArctic:    yandex2.ColorSceneIDOcean,
	sbertypes2.ColorSceneIDRomantic:  yandex2.ColorSceneIDRomance,
	sbertypes2.ColorSceneIDSunset:    yandex2.ColorSceneIDSunset,
	sbertypes2.ColorSceneIDDawn:      yandex2.ColorSceneIDSunrise,
	sbertypes2.ColorSceneIDChristmas: yandex2.ColorSceneIDGarland,
	sbertypes2.ColorSceneIDFito:      yandex2.ColorSceneIDRest,
}

var yandexMapToSberColorSceneID = map[yandex2.ColorSceneID]sbertypes2.ColorSceneID{
	yandex2.ColorSceneIDCandle:  sbertypes2.ColorSceneIDCandle,
	yandex2.ColorSceneIDOcean:   sbertypes2.ColorSceneIDArctic,
	yandex2.ColorSceneIDRomance: sbertypes2.ColorSceneIDRomantic,
	yandex2.ColorSceneIDSunset:  sbertypes2.ColorSceneIDSunset,
	yandex2.ColorSceneIDSunrise: sbertypes2.ColorSceneIDDawn,
	yandex2.ColorSceneIDGarland: sbertypes2.ColorSceneIDChristmas,
	yandex2.ColorSceneIDRest:    sbertypes2.ColorSceneIDFito,
}

func nilableBool(v bool) *bool {
	return &v
}

var (
	nilableFalse = nilableBool(false)
	nilableTrue  = nilableBool(true)
)
