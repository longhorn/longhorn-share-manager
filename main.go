package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/longhorn/longhorn-share-manager/app/cmd"
)

// following variables will be filled by `-ldflags "-X ..."`
var (
	Version   string
	GitCommit string
	BuildDate string
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
			funcName := path.Base(f.Function)
			return funcName, fileName
		},
		FullTimestamp: true,
	})

	a := &cli.Command{
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if c.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return ctx, nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug",
			},
		},
		Commands: []*cli.Command{
			cmd.ServerCmd(),
		},
	}
	if err := a.Run(context.Background(), os.Args); err != nil {
		logrus.Fatal("Error when executing command: ", err)
	}
}
