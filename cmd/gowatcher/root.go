package main

import (
	"fmt"
	"os"

	"github.com/damonchen/gowatcher/pkg/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/damonchen/gowatcher/pkg/config"
	"github.com/damonchen/gowatcher/pkg/version"
)

var (
	cfgFile string

	showVersion bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "watcher.yaml", "config file of gowatcher")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of gowatcher")
}

var rootCmd = &cobra.Command{
	Use:   "gowatcher",
	Short: "gowatcher watch file and run command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}

		// Do not show command usage here.
		err := run(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func run(cfgFile string) error {
	cfg, err := parseConfig(cfgFile)
	if err != nil {
		return err
	}

	return startService(cfg)
}

func startService(cfg *config.Config) error {
	svc := service.NewService(cfg)
	return svc.Run()
}

func parseConfig(cfgFile string) (*config.Config, error) {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.gowatcher")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
