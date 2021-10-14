package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/longhorn-share-manager/app/cmd"
)

// following variables will be filled by `-ldflags "-X ..."`

func main() {
	a := cli.NewApp()

	a.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	a.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
	}
	a.Commands = []cli.Command{
		cmd.ServerCmd(),
	}
	if err := a.Run(os.Args); err != nil {
		logrus.Fatal("Error when executing command: ", err)
	}
}
