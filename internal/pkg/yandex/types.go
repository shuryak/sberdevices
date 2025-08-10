// Package yandex is an implementation of the Yandex Smart Home REST operating protocol
// (https://yandex.ru/dev/dialogs/smart-home/doc/reference/resources.html#rest)
package yandex

// DevicesResponse - https://yandex.ru/dev/dialogs/smart-home/doc/reference/get-devices.html#output-structure
type DevicesResponse struct {
	// RequestID - request ID
	RequestID string `json:"request_id"`
	// Payload - object with devices
	Payload *DevicesResponsePayload `json:"payload,omitempty"`
}

// DevicesResponsePayload - https://yandex.ru/dev/dialogs/smart-home/doc/reference/get-devices.html#output-structure
type DevicesResponsePayload struct {
	// UserID -  user ID
	UserID string `json:"user_id"`
	// Devices - array of the user's devices
	Devices []Device `json:"devices"`
}

// Device - https://yandex.ru/dev/dialogs/smart-home/doc/reference/get-devices.html#output-structure
type Device struct {
	// ID - device ID. It must be unique among all the manufacturer's devices
	ID string `json:"id"`
	// Name - device name
	Name string `json:"name,omitempty"`
	// Description - device description
	Description string `json:"description,omitempty"`
	// Room - name of the room where the device is located
	Room string `json:"room,omitempty"`
	// Type - Device type
	Type DeviceType `json:"type,omitempty"`
	// CustomData - an object consisting of a set of "key":"value" pairs with any nesting level, providing additional
	// information about the device. Object content size must not exceed 1024 bytes. Yandex Smart Home saves this object
	// and sends it in Information about the states of the user's devices and Changing the state of devices requests
	CustomData map[string]interface{} `json:"custom_data,omitempty"`
	// Capabilities - array with information about device capabilities
	Capabilities []DeviceCapability `json:"capabilities,omitempty"`
	// Properties - array with information about the device's properties
	Properties map[string]interface{} `json:"properties,omitempty"` // TODO: property struct
	// DeviceInfo - array with information about the device's properties
	DeviceInfo *DeviceInfo `json:"device_info,omitempty"`
	// ActionResult is used in Change device state response
	// (https://yandex.ru/dev/dialogs/smart-home/doc/reference/post-action.html). Result of device state change. This
	// parameter is required if capabilities is missing
	ActionResult *DeviceActionResult `json:"action_result,omitempty"`
	// ErrorCode - an error code. If the field is filled in, the capabilities and properties parameters are ignored.
	ErrorCode ErrorCode `json:"error_code,omitempty"`
	// ErrorMessage - extended human-readable description of a possible error. Available only on the
	// https://dialogs.yandex.ru/developer/ "Testing" tab of the developer console
	ErrorMessage string `json:"error_message,omitempty"`
}

// DeviceCapability - https://yandex.ru/dev/dialogs/smart-home/doc/reference/get-devices.html#output-structure
type DeviceCapability struct {
	// Type - type of capability
	Type DeviceCapabilityType `json:"type"`
	// State - state of the device capability is used in the operations Information about the states of user devices
	// (https://yandex.ru/dev/dialogs/smart-home/doc/reference/post-devices-query.html), Notification about device state
	// change (https://yandex.ru/dev/dialogs/smart-home/doc/reference-alerts/post-skill_id-callback-state.html).
	// Commands to control device capabilities are used in the Change device state operation
	State *DeviceCapabilityState `json:"state,omitempty"`
	// Retrievable - if it's possible to request the state of this device capability. Acceptable values:
	// true (default): a state request is available for the capability; false: a state request is not available for the
	// capability
	Retrievable *bool `json:"retrievable,omitempty"`
	// Reportable - indicates that the notification service reports the capability state change. Acceptable values:
	// true: notification is enabled. The manufacturer notifies Yandex Smart Home of every change in the capability
	// state. false (default): Notification is disabled. The manufacturer doesn't notify Yandex Smart Home of the
	// capability state change
	Reportable bool `json:"reportable,omitempty"`
	// Parameters - the object must contain at least one parameter
	Parameters *DeviceCapabilitiesParameters `json:"parameters,omitempty"`
}

// DeviceInfo - https://yandex.ru/dev/dialogs/smart-home/doc/reference/get-devices.html#output-structure
type DeviceInfo struct {
	// Manufacturer - name of the device manufacturer. It can contain up to 256 characters. This parameter is required
	// in the official skill description
	Manufacturer string `json:"manufacturer"`
	// Model - name of the device model. It can contain up to 256 characters.This parameter is required in the official
	// skill description
	Model string `json:"model"`
	// HWVersion - device hardware version. It can contain up to 256 characters
	HWVersion string `json:"hw_version,omitempty"`
	// SWVersion - device software version. It can contain up to 256 characters
	SWVersion string `json:"sw_version,omitempty"`
}

// DeviceActionResult - result of device state change
// (https://yandex.ru/dev/dialogs/smart-home/doc/reference/post-action.html).
//
// We recommend that you return a parameter for each capability separately. If this is not possible, return the result
// for the entire device. The response to the request is considered successful (status:"DONE") if at least one state of
// the device has changed.
type DeviceActionResult struct {
	Status       ResultStatus `json:"status,omitempty"`
	ErrorCode    ErrorCode    `json:"error_code,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

// DeviceCapabilityState - state of the device capability is used in the operations Information about the states of user
// devices (https://yandex.ru/dev/dialogs/smart-home/doc/reference/post-devices-query.html), Notification about device
// state change (https://yandex.ru/dev/dialogs/smart-home/doc/reference-alerts/post-skill_id-callback-state.html)
type DeviceCapabilityState struct {
	Instance     DeviceInstance      `json:"instance"`
	Value        interface{}         `json:"value,omitempty"`
	ActionResult *DeviceActionResult `json:"action_result,omitempty"`
}

type DeviceCapabilitiesParameters struct {
	Split        *bool                              `json:"split,omitempty"`
	Instance     DeviceInstance                     `json:"instance,omitempty"`
	Unit         Unit                               `json:"unit,omitempty"`
	RandomAccess *bool                              `json:"random_access,omitempty"`
	Range        *DeviceCapabilitiesParametersRange `json:"range,omitempty"`
	ColorModel   ColorModel                         `json:"color_model,omitempty"`
	TemperatureK *DeviceCapabilitiesParametersRange `json:"temperature_k,omitempty"`
	ColorScene   *DeviceColorScene                  `json:"color_scene,omitempty"`
}

type DeviceCapabilitiesParametersRange struct {
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
	Precision float64 `json:"precision"`
}

type DeviceColorScene struct {
	Scenes []DeviceColorSceneItem `json:"scenes"`
}

type DeviceColorSceneItem struct {
	ID ColorSceneID `json:"id"`
}
