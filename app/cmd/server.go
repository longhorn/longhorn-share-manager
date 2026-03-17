package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/types/pkg/generated/smrpc"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/longhorn/longhorn-share-manager/pkg/rpc"
	"github.com/longhorn/longhorn-share-manager/pkg/server"
	"github.com/longhorn/longhorn-share-manager/pkg/types"
	"github.com/longhorn/longhorn-share-manager/pkg/util"
	"github.com/longhorn/longhorn-share-manager/pkg/volume"
)

const (
	listenPort = ":9600"
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
			cli.StringFlag{
				Name:     "data-engine",
				Usage:    "The volume data engine",
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
				Name:     "cryptokeycipher",
				Usage:    "contains the encryption algorithm in dm-crypt notation",
				EnvVar:   "CRYPTOKEYCIPHER",
				Required: false,
			},
			cli.StringFlag{
				Name:     "cryptokeyhash",
				Usage:    "contains the hash algorithm for the checksum resilience mode",
				EnvVar:   "CRYPTOKEYHASH",
				Required: false,
			},
			cli.StringFlag{
				Name:     "cryptokeysize",
				Usage:    "contains the encryption key size",
				EnvVar:   "CRYPTOKEYSIZE",
				Required: false,
			},
			cli.StringFlag{
				Name:     "cryptopbkdf",
				Usage:    "contains the Password-Based Key Derivation Function (PBKDF) algorithm for LUKS keyslot",
				EnvVar:   "CRYPTOPBKDF",
				Required: false,
			},
			cli.StringFlag{
				Name:     "cryptopbkdfiterations",
				Usage:    "contains the forced iteration count for the PBKDF algorithm (higher values increase security but reduce performance)",
				EnvVar:   "CRYPTOPBKDFFORCEITERATIONS",
				Required: false,
			},
			cli.StringFlag{
				Name:     "cryptopbkdfmemory",
				Usage:    "contains the memory cost parameter for the PBKDF algorithm in KiB (as used by cryptsetup)",
				EnvVar:   "CRYPTOPBKDFMEMORY",
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
				Name:                       c.String("volume"),
				DataEngine:                 c.String("data-engine"),
				Passphrase:                 c.String("passphrase"),
				CryptoKeyCipher:            c.String("cryptokeycipher"),
				CryptoKeyHash:              c.String("cryptokeyhash"),
				CryptoKeySize:              c.String("cryptokeysize"),
				CryptoPBKDF:                c.String("cryptopbkdf"),
				CryptoPBKDFForceIterations: c.String("cryptopbkdfiterations"),
				CryptoPBKDFMemory:          c.String("cryptopbkdfmemory"),
				FsType:                     c.String("fs"),
				MountOptions:               c.StringSlice("mount"),
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
	if vol.DataEngine != types.DataEngineTypeV1 && vol.DataEngine != types.DataEngineTypeV2 {
		logger.Errorf("Invalid data engine value: %s", vol.DataEngine)
		return fmt.Errorf("invalid data engine value: %s", vol.DataEngine)
	}

	manager, err := server.NewShareManager(logger, vol)
	if err != nil {
		return err
	}

	shutdownCh := make(chan error)
	defer close(shutdownCh)
	go func() {
		err := manager.Run()
		shutdownCh <- err
	}()

	go func() {
		listen, err := net.Listen("tcp", listenPort)
		if err != nil {
			logrus.WithError(err).Warnf("Failed to listen on port %v", listenPort)
			shutdownCh <- err
			return
		}

		s := grpc.NewServer()
		srv := rpc.NewShareManagerServer(manager)
		smrpc.RegisterShareManagerServiceServer(s, srv)
		healthpb.RegisterHealthServer(s, rpc.NewShareManagerHealthCheckServer(srv))
		reflection.Register(s)

		logrus.Infof("Listening on share manager gRPC server %s", listenPort)
		err = s.Serve(listen)
		logrus.WithError(err).Warnf("Share manager gRPC server at %v is down", listenPort)
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
