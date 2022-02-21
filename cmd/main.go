package main

import (
	"context"
	"os"

	"github.com/frain-dev/immune/system"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	err := os.Setenv("TZ", "") // Use UTC by default :)
	if err != nil {
		log.WithError(err).Fatal("failed to set env")
	}

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
	}

	var configFile string
	cmd.PersistentFlags().StringVar(&configFile, "config", "./immune.json", "Configuration file for immune")

	cmd.AddCommand(addRunCommand())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func addRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"serve", "s"},
		Short:   "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			log.Infof("all tests passed")
			return nil
		},
	}
	return cmd
}
