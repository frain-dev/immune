package main

import (
	"context"

	"github.com/frain-dev/immune/system"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	cmd := &cobra.Command{
		Use:   "Immune",
		Short: "API Testing tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			sys, err := system.NewSystem(cfgPath)
			if err != nil {
				return err
			}

			err = sys.Clean()
			if err != nil {
				return err
			}

			err = sys.Run(context.Background())
			if err != nil {
				return err
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	var configFile string
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "Configuration file for immune")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
