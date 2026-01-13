package improv

import (
	"log/slog"

	"github.com/dchauviere/spkctld/internal/wifi"
	"tinygo.org/x/bluetooth"
)

var logger = slog.With(slog.String("module", "ble"))

var (
	adapter   = bluetooth.DefaultAdapter
	localName = "Speaker Improv"
)

type ImprovService struct {
	commandCharacteristic bluetooth.Characteristic
	stateCharacteristic   bluetooth.Characteristic
	errorCharacteristic   bluetooth.Characteristic
	rpcCharacteristic     bluetooth.Characteristic

	backend wifi.WifiBackend
}

func NewImprovService(backend wifi.WifiBackend) *ImprovService {
	return &ImprovService{backend: backend}
}

func (s *ImprovService) Start() error {
	svcUUID, _ := bluetooth.ParseUUID(ServiceUUID)
	commandUUID, _ := bluetooth.ParseUUID(CommandUUID)
	stateUUID, _ := bluetooth.ParseUUID(StateUUID)
	errorUUID, _ := bluetooth.ParseUUID(ErrorUUID)
	rpcUUID, _ := bluetooth.ParseUUID(RpcUUID)

	adv := adapter.DefaultAdvertisement()
	advUUID, _ := bluetooth.ParseUUID(ServiceUUID)
	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    localName,
		ServiceUUIDs: []bluetooth.UUID{advUUID},
	})

	if err := adv.Start(); err != nil {
		return err
	}

	return adapter.AddService(&bluetooth.Service{
		UUID: svcUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &s.commandCharacteristic,
				UUID:   commandUUID,
				Flags:  bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					s.handleCommand(value)
				},
			},
			{
				Handle: &s.stateCharacteristic,
				UUID:   stateUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &s.errorCharacteristic,
				UUID:   errorUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &s.rpcCharacteristic,
				UUID:   rpcUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	})

}

func (s *ImprovService) Stop() {
	// tu compléteras si tu veux unregister proprement
}

func (s *ImprovService) NotifyState(state byte) {
	s.stateCharacteristic.Write([]byte{state})
}

func (s *ImprovService) NotifyError(code byte) {
	s.errorCharacteristic.Write([]byte{code})
}

func (s *ImprovService) NotifyRpc(data []byte) {
	s.rpcCharacteristic.Write(data)
}

func (s *ImprovService) handleCommand(raw []byte) {
	if len(raw) < 2 {
		s.NotifyError(ErrorInvalidRPCPacket)
		return
	}
	opcode := raw[0]

	switch opcode {
	case OpcodeSetWifi:
		ssidLen := int(raw[1])
		if len(raw) < 2+ssidLen {
			s.NotifyError(ErrorInvalidRPCPacket)
			return
		}
		ssid := string(raw[2 : 2+ssidLen])
		pwd := string(raw[2+ssidLen:])

		logger.Info("SET_WIFI received", "ssid", ssid)
		s.NotifyState(StateProvisioning)

		s.backend.Connect(ssid, pwd, func(success bool) {
			if success {
				logger.Info("WiFi connected", "ssid", ssid)
				s.NotifyState(StateProvisioned)
				s.NotifyRpc([]byte("http://device.local/setup"))
			} else {
				logger.Error("WiFi connection failed", "ssid", ssid)
				s.NotifyState(StateAuthorizationRequired)
				s.NotifyError(ErrorConnectFailed)
			}
		})

	default:
		logger.Error("unknown opcode", "opcode", opcode)
		s.NotifyError(ErrorUnknownRPCCommand)
	}
}

func (s *ImprovService) Reset() {
	logger.Info("WiFi reset requested")
	if err := s.backend.Reset(); err != nil {
		logger.Error("WiFi reset failed", "error", err)
	}
	s.NotifyState(StateAuthorizationRequired)
	s.NotifyError(ErrorNone)
}
