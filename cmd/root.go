/*
Copyright © 2026 David Chauviere david.chauviere@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dchauviere/spkctld/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootCmd struct {
	cfgFile string
}

func (rc *RootCmd) Command() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands.
	var rootCmd = &cobra.Command{
		Use:   "spkctld",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initConfig(rc.cfgFile)
			logging.Init(viper.GetString("log.level"))
		},
	}
	rc.init(rootCmd)
	rootCmd.AddCommand((&serveCommand{}).Command())
	return rootCmd
}

func (rc *RootCmd) init(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&rc.cfgFile,
		"config",
		"",
		"config file (default is /etc/spkctld.yaml)",
	)
	cmd.PersistentFlags().StringP("log-level", "l", "info", "Log level")
	_ = viper.BindPFlag("log.level", cmd.PersistentFlags().Lookup("log-level"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("spkctld")
	}

	viper.AutomaticEnv()          // read in environment variables that match
	viper.SetEnvPrefix("spkctld") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
