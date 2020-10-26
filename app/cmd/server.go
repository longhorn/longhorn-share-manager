package cmd

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/longhorn-share-manager/pkg/health"
	"github.com/longhorn/longhorn-share-manager/pkg/rpc"
	"github.com/longhorn/longhorn-share-manager/pkg/server"
	"github.com/longhorn/longhorn-share-manager/pkg/util"
)

func ServerCmd() cli.Command {
	return cli.Command{
		Name: "daemon",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "listen",
				Value: "localhost:8500",
			},
			cli.StringFlag{
				Name:  "log",
				Value: "/export/share-manager.log",
			},
		},
		Action: func(c *cli.Context) {
			listen := c.String("listen")
			logFile := c.String("log")
			if err := start(logFile, listen); err != nil {
				logrus.Fatalf("Error running start command: %v.", err)
			}
		},
	}
}

func start(logFile string, listen string) error {
	logger, err := util.NewLogger(logFile)
	if err != nil {
		return err
	}

	manager, err := server.NewShareManager(logger, logFile)
	if err != nil {
		return err
	}
	hc := health.NewHealthCheckServer(manager)

	listenAt, err := net.Listen("tcp", listen)
	if err != nil {
		return errors.Wrap(err, "Failed to listen")
	}

	rpcService := grpc.NewServer()
	rpc.RegisterShareManagerServiceServer(rpcService, manager)
	healthpb.RegisterHealthServer(rpcService, hc)
	reflection.Register(rpcService)

	shutdownCh := make(chan error)
	go func() {
		if err := rpcService.Serve(listenAt); err != nil {
			logger.Errorf("Stopping due to %v:", err)
		}
		manager.Shutdown()
		shutdownCh <- err
	}()
	logger.Infof("share manager listening to %v", listen)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Infof("share manager received signal %v to exit", sig)
		rpcService.Stop()
	}()

	// TODO: graceful shutdown
	return <-shutdownCh
}
