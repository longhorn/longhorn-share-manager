package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/longhorn-share-manager/pkg/client"
	"github.com/longhorn/longhorn-share-manager/pkg/meta"
)

func VersionCmd() cli.Command {
	return cli.Command{
		Name: "version",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "client-only",
			},
		},
		Action: func(c *cli.Context) {
			if err := version(c.GlobalString("url"), c.Bool("client-only")); err != nil {
				logrus.Fatalf("Error running version command: %v", err)
			}
		},
	}
}

type VersionOutput struct {
	ClientVersion *meta.VersionOutput `json:"clientVersion"`
	ServerVersion *meta.VersionOutput `json:"serverVersion"`
}

func version(url string, clientOnly bool) error {
	clientVersion := meta.GetVersion()
	v := VersionOutput{ClientVersion: &clientVersion}

	if !clientOnly {
		cli := client.NewShareManagerClient(url)
		version, err := cli.VersionGet()
		if err != nil {
			return err
		}
		v.ServerVersion = version
	}

	output, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
