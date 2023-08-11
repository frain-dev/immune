package main

import (
	"context"

	"github.com/frain-dev/immune/fire"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func addFireCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fire",
		Short: "Fire events at a convoy instance",
		Run: func(cmd *cobra.Command, args []string) {
			firer := fire.NewFire(a.config)
			l, err := firer.Start(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			err = l.WriteToFile(a.config.LogFile)
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("Fire command completed successfully")
		},
	}
	return cmd
}
