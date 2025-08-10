package sbertypes

import "time"

type DevicesResponse struct {
	Result     DeviceItems `json:"result"`
	Pagination Pagination  `json:"pagination"`
}

type DeviceResponse struct {
	Result DeviceItem `json:"result"`
}

type DeviceItems []DeviceItem

func (di DeviceItems) ToMap() map[string]DeviceItem {
	res := make(map[string]DeviceItem)

	for i := range di {
		res[di[i].ID] = di[i]
	}

	return res
}

type DeviceItem struct {
	ID             string            `json:"id"`
	DeviceTypeID   string            `json:"device_type_id"`
	Name           *DeviceName       `json:"name,omitempty"`
	DeviceInfo     *DeviceInfo       `json:"device_info,omitempty"`
	Attributes     []DeviceAttribute `json:"attributes"`
	ReportedState  DeviceStates      `json:"reported_state"`
	DesiredState   []DeviceState     `json:"desired_state"`
	Commands       []DeviceCommand   `json:"commands"`
	SerialNumber   string            `json:"serial_number"`
	ExternalID     string            `json:"external_id"`
	Images         *DeviceImages     `json:"images,omitempty"`
	Categories     []string          `json:"categories"`
	GroupIDs       []string          `json:"group_ids"`
	DeviceTypeName string            `json:"device_type_name"`
	HWVersion      string            `json:"hw_version"`
	SWVersion      string            `json:"sw_version"`
	FullCategories []DeviceCategory  `json:"full_categories"`

	// TODO: owner_info, etc
}

type DeviceStates []DeviceState

func (ds DeviceStates) ToMap() map[StateCommand]DeviceState {
	res := make(map[StateCommand]DeviceState)

	for i := range ds {
		res[ds[i].Key] = ds[i]
	}

	return res
}

type DeviceName struct {
	Name        string   `json:"name"`
	DefaultName string   `json:"defaultName,omitempty"`
	Nickname    []string `json:"nickname,omitempty"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	HWVersion    string `json:"hw_version,omitempty"`
	SWVersion    string `json:"sw_version,omitempty"`
	Description  string `json:"description,omitempty"`
	ProductID    string `json:"product_id,omitempty"`
	Partner      string `json:"partner,omitempty"`
	SWVersionInt string `json:"sw_version_int,omitempty"`
}

type DeviceAttribute struct {
	Key         StateCommand                `json:"key"`
	Type        SberDataType                `json:"type"`
	IntValues   *DeviceAttributeIntValues   `json:"int_values,omitempty"`
	EnumValues  *DeviceAttributeEnumValues  `json:"enum_values,omitempty"`
	ColorValues *DeviceAttributeColorValues `json:"color_values,omitempty"`
	Name        string                      `json:"name,omitempty"`
	IsVisible   bool                        `json:"is_visible"`
	// TODO: float_value and etc
}

type DeviceAttributeIntValues struct {
	Range DeviceAttributeIntValuesRange `json:"range"`
	Unit  string                        `json:"unit"`
}

type DeviceAttributeIntValuesRange struct {
	Min  int `json:"min"`
	Max  int `json:"max"`
	Step int `json:"step"`
}

type DeviceAttributeEnumValues struct {
	Values []string `json:"values"`
}

type DeviceAttributeColorValues struct {
	Hue        DeviceAttributeIntValuesRange `json:"h"` // TODO: maybe separate data type?
	Saturation DeviceAttributeIntValuesRange `json:"s"`
	Value      DeviceAttributeIntValuesRange `json:"v"`
}

type DeviceState struct {
	Key          StateCommand           `json:"key"`
	Type         SberDataType           `json:"type"`
	FloatValue   float64                `json:"float_value,omitempty"`
	IntegerValue string                 `json:"integer_value,omitempty"`
	StringValue  string                 `json:"string_value,omitempty"`
	BoolValue    bool                   `json:"bool_value,omitempty"`
	EnumValue    string                 `json:"enum_value,omitempty"`
	ColorValue   *DeviceStateColorValue `json:"color_value,omitempty"`
	LastSync     time.Time              `json:"last_sync,omitempty"`
}

type DeviceStateColorValue struct {
	Hue        int `json:"h"`
	Saturation int `json:"s"`
	Value      int `json:"v"`
}

type DeviceCommand struct {
	Key         StateCommand   `json:"key"`
	StateFields []StateCommand `json:"state_fields"` // TODO: separate state and key?
	// TODO: exceptions
}

type DeviceImages struct {
	LauncherExtraLargeOff  string `json:"launcher_extra_large_off,omitempty"`
	LauncherExtraLargeOn   string `json:"launcher_extra_large_on,omitempty"`
	LauncherLargePromoOff  string `json:"launcher_large_promo_off,omitempty"`
	LauncherLargePromoOn   string `json:"launcher_large_promo_on,omitempty"`
	LauncherSmallBoxOff    string `json:"launcher_small_box_off,omitempty"`
	LauncherSmallBoxOn     string `json:"launcher_small_box_on,omitempty"`
	LauncherSmallPortalOff string `json:"launcher_small_portal_off,omitempty"`
	LauncherSmallPortalOn  string `json:"launcher_small_portal_on,omitempty"`
	ListOff                string `json:"list_off,omitempty"`
	ListOn                 string `json:"list_on,omitempty"`
	Photo                  string `json:"photo,omitempty"`
}

type DeviceCategory struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	DefaultName  string `json:"default_name"`
	ImageSetType string `json:"image_set_type"`
	SortWeight   int    `json:"sort_weight"`
	// TODO: etc
}

type ColorSceneID string

const (
	ColorSceneIDCandle    ColorSceneID = "candle"
	ColorSceneIDArctic    ColorSceneID = "arctic"
	ColorSceneIDRomantic  ColorSceneID = "romantic"
	ColorSceneIDSunset    ColorSceneID = "sunset"
	ColorSceneIDDawn      ColorSceneID = "dawn"
	ColorSceneIDChristmas ColorSceneID = "christmas"
	ColorSceneIDFito      ColorSceneID = "fito"
)
