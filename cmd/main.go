package main

import (
	"os"

	"github.com/frain-dev/immune/config"
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

	a := &App{}
	cmd := &cobra.Command{
		Use:               "Immune",
		Short:             "API Testing tool",
		PersistentPreRunE: PreRun(a),
	}

	var configFile string
	cmd.PersistentFlags().StringVar(&configFile, "config", "./immune.json", "Configuration file for immune")

	cmd.AddCommand(addFireCommand(a))
	cmd.AddCommand(addRecvCommand(a))

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

type App struct {
	config *config.Config
}

func PreRun(app *App) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfgPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		sys, err := config.LoadConfig(cfgPath)
		if err != nil {
			return err
		}

		err = sys.Validate()
		if err != nil {
			return err
		}

		app.config = sys
		return nil
	}
}
