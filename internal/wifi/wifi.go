package wifi

type WifiBackend interface {
	Connect(ssid, password string, cb func(success bool))
	Reset() error
}
