package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/longhorn-share-manager/pkg/server"
	"github.com/longhorn/longhorn-share-manager/pkg/util"
	"github.com/longhorn/longhorn-share-manager/pkg/volume"
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
			cli.BoolFlag{
				Name:     "encrypted",
				Usage:    "signals that a volume is encrypted",
				EnvVar:   "ENCRYPTED",
				Required: false,
			},
			cli.StringFlag{
				Name:     "passphrase",
				Usage:    "contains the encryption passphrase",
				EnvVar:   "PASSPHRASE",
				Required: false,
			},
			cli.StringFlag{
				Name:     "fs",
				Usage:    "the filesystem to use for the volume",
				Value:    "ext4",
				Required: false,
			},
			cli.StringSliceFlag{
				Name:     "mount",
				Usage:    "allows for specifying additional mount options",
				Required: false,
			},
		},
		Action: func(c *cli.Context) {
			vol := volume.Volume{
				Name:         c.String("volume"),
				Passphrase:   c.String("passphrase"),
				FsType:       c.String("fs"),
				MountOptions: c.StringSlice("mount"),
			}

			if c.Bool("encrypted") && len(vol.Passphrase) == 0 {
				logrus.Fatalf("Error starting share-manager missing passphrase for encrypted volume %v", vol.Name)
			}

			if err := start(vol); err != nil {
				logrus.Fatalf("Error running start command: %v.", err)
			}
		},
	}
}

func start(vol volume.Volume) error {
	logger := util.NewLogger()
	manager, err := server.NewShareManager(logger, vol)
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
