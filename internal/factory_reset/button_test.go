package factory_reset_test

import (
	"testing"
	"time"

	"github.com/dchauviere/spkctld/internal/factory_reset"
	evdev "github.com/gvalkov/golang-evdev"
)

type mockDev struct {
	events [][]evdev.InputEvent
	idx    int
}

func (m *mockDev) Read() ([]evdev.InputEvent, error) {
	if m.idx >= len(m.events) {
		time.Sleep(5 * time.Millisecond)
		return []evdev.InputEvent{}, nil
	}
	ev := m.events[m.idx]
	m.idx++
	return ev, nil
}

func TestButtonTriggersReset(t *testing.T) {

	pressed := false
	cb := func() { pressed = true }

	mock := &mockDev{
		events: [][]evdev.InputEvent{
			{
				{Type: evdev.EV_KEY, Code: 0x198, Value: 1},
			},
		},
	}

	watcher := factory_reset.NewButtonWatcherWithDevice(mock, 0x198)
	watcher.Run(cb)

	time.Sleep(20 * time.Millisecond)

	if !pressed {
		t.Fatalf("expected reset callback")
	}
}
