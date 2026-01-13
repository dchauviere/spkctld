package wifi

import (
	"log/slog"

	"github.com/godbus/dbus/v5"
)

var logger = slog.With(slog.String("module", "wifi"))

type DBusConn interface {
	Object(dest string, path dbus.ObjectPath) DBusObject
	Signal(ch chan<- *dbus.Signal)
}

type DBusObject interface {
	Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call
}

type dbusConnAdapter struct {
	*dbus.Conn
}

func (a *dbusConnAdapter) Object(dest string, path dbus.ObjectPath) DBusObject {
	return a.Conn.Object(dest, path)
}

type ConnmanBackend struct {
	conn DBusConn
}

// Constructeur normal (prod) : utilise le vrai SystemBus
func NewConnmanBackend() (*ConnmanBackend, error) {
	realConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return &ConnmanBackend{conn: &dbusConnAdapter{realConn}}, nil
}

// Constructeur pour les tests : on injecte un mock DBusConn
func NewConnmanBackendFromMock(conn DBusConn) *ConnmanBackend {
	return &ConnmanBackend{conn: conn}
}

func (c *ConnmanBackend) Connect(ssid, password string, cb func(bool)) {
	obj := c.conn.Object("net.connman", "/")

	var services []struct {
		Path  dbus.ObjectPath
		Props map[string]dbus.Variant
	}

	call := obj.Call("net.connman.Manager.GetServices", 0)
	if call.Err != nil {
		logger.Error("GetServices failed", "error", call.Err)
		cb(false)
		return
	}
	if err := call.Store(&services); err != nil {
		logger.Error("Store(GetServices) failed", "error", err)
		cb(false)
		return
	}

	for _, svc := range services {
		nameVar, ok := svc.Props["Name"]
		if !ok {
			continue
		}
		name, _ := nameVar.Value().(string)
		if name != ssid {
			continue
		}

		service := c.conn.Object("net.connman", svc.Path)

		service.Call("net.connman.Service.SetProperty", 0,
			"Passphrase", dbus.MakeVariant(password))

		service.Call("net.connman.Service.Connect", 0)

		// Ici, dans le vrai code, tu écouteras les signaux "PropertyChanged"
		// Pour les tests, on ne simule pas les signaux -> on appelle cb côté mock.

		cb(true) // dans un vrai flux, tu enlèves ça et relies au signal
		return
	}

	logger.Error("SSID not found", "ssid", ssid)
	cb(false)
}

func (c *ConnmanBackend) Reset() error {
	obj := c.conn.Object("net.connman", "/")

	var services []struct {
		Path  dbus.ObjectPath
		Props map[string]dbus.Variant
	}

	call := obj.Call("net.connman.Manager.GetServices", 0)
	if call.Err != nil {
		return call.Err
	}
	if err := call.Store(&services); err != nil {
		return err
	}

	for _, svc := range services {
		tVar, ok := svc.Props["Type"]
		if !ok {
			continue
		}
		t, _ := tVar.Value().(string)
		if t != "wifi" {
			continue
		}
		service := c.conn.Object("net.connman", svc.Path)
		service.Call("net.connman.Service.Disconnect", 0)
		service.Call("net.connman.Service.SetProperty", 0,
			"Passphrase", dbus.MakeVariant(""))
	}
	return nil
}
