package main

import (
	"github.com/frain-dev/immune/recv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func addRecvCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recv",
		Short: "Fire events at a convoy instance",
		Run: func(cmd *cobra.Command, args []string) {
			rc := recv.NewReceiver(a.config)

			l := rc.Listen()

			err := l.WriteToFile(a.config.RecvLogFile)
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("Recv command completed successfully")
		},
	}
	return cmd
}
