package wifi_test

import (
	"testing"

	"github.com/dchauviere/spkctld/internal/wifi"
	"github.com/godbus/dbus/v5"
)

type mockDBusConn struct {
	services []struct {
		Path  dbus.ObjectPath
		Props map[string]dbus.Variant
	}
	calls []string
}

func (m *mockDBusConn) Object(dest string, path dbus.ObjectPath) wifi.DBusObject {
	return &mockDBusObject{mock: m, path: path}
}

func (m *mockDBusConn) Signal(ch chan<- *dbus.Signal) {}

type mockDBusObject struct {
	mock *mockDBusConn
	path dbus.ObjectPath
}

func (o *mockDBusObject) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	o.mock.calls = append(o.mock.calls, method)

	switch method {
	case "net.connman.Manager.GetServices":
		return &dbus.Call{
			Body: []interface{}{o.mock.services},
		}
	default:
		return &dbus.Call{}
	}
}

func TestConnmanConnectSuccess(t *testing.T) {

	mock := &mockDBusConn{
		services: []struct {
			Path  dbus.ObjectPath
			Props map[string]dbus.Variant
		}{
			{
				Path: "/service/wifi1",
				Props: map[string]dbus.Variant{
					"Name": dbus.MakeVariant("TEST"),
					"Type": dbus.MakeVariant("wifi"),
				},
			},
		},
	}

	backend := wifi.NewConnmanBackendFromMock(mock)

	called := false
	backend.Connect("TEST", "pass", func(ok bool) {
		called = true
		if !ok {
			t.Fatalf("expected success")
		}
	})

	if !called {
		t.Fatalf("callback not called")
	}

	if len(mock.calls) == 0 {
		t.Fatalf("expected DBus calls")
	}
}

func TestConnmanConnectSSIDNotFound(t *testing.T) {

	mock := &mockDBusConn{
		services: []struct {
			Path  dbus.ObjectPath
			Props map[string]dbus.Variant
		}{},
	}

	backend := wifi.NewConnmanBackendFromMock(mock)

	called := false
	backend.Connect("UNKNOWN", "pass", func(ok bool) {
		called = true
		if ok {
			t.Fatalf("expected failure")
		}
	})

	if !called {
		t.Fatalf("callback not called")
	}
}
