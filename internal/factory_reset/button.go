package factory_reset

import (
	"log/slog"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
)

var logger = slog.With(slog.String("module", "button"))

type InputDevice interface {
	Read() ([]evdev.InputEvent, error)
}

type ResetCallback func()

type ButtonWatcher struct {
	dev     InputDevice
	keyCode uint16
}

func NewButtonWatcher(devPath string, keyCode uint16) *ButtonWatcher {
	dev, err := evdev.Open(devPath)
	if err != nil {
		logger.Error("cannot open button device", "path", devPath, "error", err)
		return nil
	}
	return &ButtonWatcher{dev: dev, keyCode: keyCode}
}

func NewButtonWatcherWithDevice(dev InputDevice, keyCode uint16) *ButtonWatcher {
	return &ButtonWatcher{dev: dev, keyCode: keyCode}
}

func (b *ButtonWatcher) Run(cb func()) {
	if b.dev == nil {
		slog.Error("ButtonWatcher has no device")
		return
	}
	go func() {
		for {
			events, err := b.dev.Read()
			if err != nil {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			for _, e := range events {
				if e.Type == evdev.EV_KEY && e.Code == b.keyCode && e.Value == 1 {
					slog.Info("reset button pressed")
					cb()
				}
			}
		}
	}()
}
