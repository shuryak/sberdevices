package yandex

type ErrorCode string

const (
	ErrorCodeDoorOpen                  ErrorCode = "DOOR_OPEN"
	ErrorCodeLidOpen                   ErrorCode = "LID_OPEN"
	ErrorCodeRemoteControlDisabled     ErrorCode = "REMOTE_CONTROL_DISABLED"
	ErrorCodeNotEnoughWater            ErrorCode = "NOT_ENOUGH_WATER"
	ErrorCodeLowChargeLevel            ErrorCode = "LOW_CHARGE_LEVEL"
	ErrorCodeContainerFull             ErrorCode = "CONTAINER_FULL"
	ErrorCodeContainerEmpty            ErrorCode = "CONTAINER_EMPTY"
	ErrorCodeDripTrayFull              ErrorCode = "DRIP_TRAY_FULL"
	ErrorCodeDeviceStuck               ErrorCode = "DEVICE_STUCK"
	ErrorCodeDeviceOff                 ErrorCode = "DEVICE_OFF"
	ErrorCodeFirmwareOutOfDate         ErrorCode = "FIRMWARE_OUT_OF_DATE"
	ErrorCodeNotEnoughDetergent        ErrorCode = "NOT_ENOUGH_DETERGENT"
	ErrorCodeHumanInvolvementNeeded    ErrorCode = "HUMAN_INVOLVEMENT_NEEDED"
	ErrorCodeDeviceUnreachable         ErrorCode = "DEVICE_UNREACHABLE"
	ErrorCodeDeviceBusy                ErrorCode = "DEVICE_BUSY"
	ErrorCodeInternalError             ErrorCode = "INTERNAL_ERROR"
	ErrorCodeInvalidAction             ErrorCode = "INVALID_ACTION"
	ErrorCodeInvalidValue              ErrorCode = "INVALID_VALUE"
	ErrorCodeNotSupportedInCurrentMode ErrorCode = "NOT_SUPPORTED_IN_CURRENT_MODE"
	ErrorCodeAccountLinkingError       ErrorCode = "ACCOUNT_LINKING_ERROR"
	ErrorCodeDeviceNotFound            ErrorCode = "DEVICE_NOT_FOUND"
)
