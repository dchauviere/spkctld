package improv

const (
	ServiceUUID = "00467768-6228-2272-4663-277478268000"
	CommandUUID = "00467768-6228-2272-4663-277478268001"
	StateUUID   = "00467768-6228-2272-4663-277478268002"
	ErrorUUID   = "00467768-6228-2272-4663-277478268003"
	RpcUUID     = "00467768-6228-2272-4663-277478268004"
)

const (
	CapabilityIdentify   = 0x01
	CapabilityDeviceInfo = 0x02
	CapabilityWifiScan   = 0x04
	CapabilityHostname   = 0x08
)

const (
	StateAuthorizationRequired = 0x01
	StateAuthorized            = 0x02
	StateProvisioning          = 0x03
	StateProvisioned           = 0x04
)

const (
	ErrorNone              = 0x00
	ErrorInvalidRPCPacket  = 0x01
	ErrorUnknownRPCCommand = 0x02
	ErrorConnectFailed     = 0x03
	ErrorNotAuthorized     = 0x04
	ErrorBadHostname       = 0x05
	ErrorUnknown           = 0xFF
)

const (
	OpcodeSetWifi        = 0x01
	OpcodeIdentify       = 0x02
	OpcodeDeviceInfo     = 0x03
	OpcodeScanWifi       = 0x04
	OpcodeGetSetHostname = 0x05
)
