package main

import (
	"github.com/frain-dev/immune/recv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func addRecvCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fire",
		Short: "Fire events at a convoy instance",
		Run: func(cmd *cobra.Command, args []string) {
			rc := recv.NewReceiver(a.config)

			rc.Listen()

			log.Infof("Fire command completed successfully")
		},
	}
	return cmd
}
