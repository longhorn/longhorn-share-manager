package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/longhorn-share-manager/pkg/server"
	"github.com/longhorn/longhorn-share-manager/pkg/util"
)

func ServerCmd() cli.Command {
	return cli.Command{
		Name: "daemon",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "volume",
				Usage:    "The volume to export via the nfs server",
				Required: true,
			},
		},
		Action: func(c *cli.Context) {
			volume := c.String("volume")
			if err := start(volume); err != nil {
				logrus.Fatalf("Error running start command: %v.", err)
			}
		},
	}
}

func start(volume string) error {
	logger := util.NewLogger()
	manager, err := server.NewShareManager(logger, volume)
	if err != nil {
		return err
	}

	shutdownCh := make(chan error)
	go func() {
		err := manager.Run()
		shutdownCh <- err
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Infof("share manager received signal %v to exit", sig)
		manager.Shutdown()
	}()

	return <-shutdownCh
}
