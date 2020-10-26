package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/longhorn-share-manager/pkg/client"
	"github.com/longhorn/longhorn-share-manager/pkg/util"
)

func ShareCmd() cli.Command {
	return cli.Command{
		Name: "share",
		Subcommands: []cli.Command{
			ShareCreateCmd(),
			ShareDeleteCmd(),
			ShareGetCmd(),
			ShareListCmd(),
		},
	}
}

func ShareCreateCmd() cli.Command {
	return cli.Command{
		Name: "create",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "volume",
			},
		},
		Action: func(c *cli.Context) {
			if err := createShare(c.GlobalString("url"), c.String("volume")); err != nil {
				logrus.Fatalf("Error running share create command: %v", err)
			}
		},
	}
}

func createShare(url, volume string) error {
	c := client.NewShareManagerClient(url)
	share, err := c.ShareCreate(volume)
	if err != nil {
		return fmt.Errorf("failed to create share: %v", err)
	}
	return util.PrintJSON(share)
}

func ShareDeleteCmd() cli.Command {
	return cli.Command{
		Name: "delete",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "volume",
			},
		},
		Action: func(c *cli.Context) {
			if err := deleteShare(c.GlobalString("url"), c.String("volume")); err != nil {
				logrus.Fatalf("Error running share delete command: %v", err)
			}
		},
	}
}

func deleteShare(url, volume string) error {
	c := client.NewShareManagerClient(url)
	share, err := c.ShareDelete(volume)
	if err != nil {
		return fmt.Errorf("failed to delete share: %v", err)
	}
	return util.PrintJSON(share)
}

func ShareGetCmd() cli.Command {
	return cli.Command{
		Name: "get",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "volume",
			},
		},
		Action: func(c *cli.Context) {
			if err := getShare(c.GlobalString("url"), c.String("volume")); err != nil {
				logrus.Fatalf("Error running share get command: %v", err)
			}
		},
	}
}

func getShare(url, volume string) error {
	c := client.NewShareManagerClient(url)
	process, err := c.ShareGet(volume)
	if err != nil {
		return fmt.Errorf("failed to get share: %v", err)
	}
	return util.PrintJSON(process)
}

func ShareListCmd() cli.Command {
	return cli.Command{
		Name:      "list",
		ShortName: "ls",
		Action: func(c *cli.Context) {
			if err := listShares(c.GlobalString("url")); err != nil {
				logrus.Fatalf("Error running share list command: %v", err)
			}
		},
	}
}

func listShares(url string) error {
	c := client.NewShareManagerClient(url)
	shares, err := c.ShareList()
	if err != nil {
		return fmt.Errorf("failed to list shares: %v", err)
	}
	return util.PrintJSON(shares)
}
