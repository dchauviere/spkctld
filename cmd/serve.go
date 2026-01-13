package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dchauviere/spkctld/internal/factory_reset"
	"github.com/dchauviere/spkctld/internal/http"
	"github.com/dchauviere/spkctld/internal/improv"
	"github.com/dchauviere/spkctld/internal/wifi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type serveCommand struct{}

func (sc *serveCommand) Command() *cobra.Command {

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "start the improv daemon",
		Long:  `start the improv daemon which allows to configure WiFi over BLE`,
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("starting improv daemon")

			backend, err := wifi.NewConnmanBackend()
			if err != nil {
				slog.Error("connman init failed", "error", err)
			}

			improvService := improv.NewImprovService(backend)

			// HTTP server for next_url
			go http.StartHttpServer()

			// Factory reset
			if viper.GetBool("factory-reset.enabled") {
				slog.Info("factory reset enabled")
				// Bouton reset (adapter devPath + keyCode à ton board)
				button := factory_reset.NewButtonWatcher(viper.GetString("reset-input-device"), 0x74) // 0x198 = exemple KEY_POWER
				button.Run(improvService.Reset)
			} else {
				slog.Info("factory reset disabled")
			}

			if err := improvService.Start(); err != nil {
				slog.Error("failed to start BLE service", "error", err)
			}

			slog.Info("improv BLE active")

			// Graceful shutdown
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig

			slog.Info("shutting down improv daemon")
			improvService.Stop()
		},
	}

	sc.init(serveCmd)
	return serveCmd
}

func (sc *serveCommand) init(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("reset-input-device", "r", "/dev/input/event0", "Reset input device path")
	_ = viper.BindPFlag("reset-input-device", cmd.PersistentFlags().Lookup("reset-input-device"))

	cmd.Flags().Bool("factory-reset-enable", false, "Factory reset enabled")
	_ = viper.BindPFlag("factory-reset.enabled", cmd.Flags().Lookup("factory-reset-enable"))
}
