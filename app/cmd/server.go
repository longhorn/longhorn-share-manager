package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/longhorn/types/pkg/generated/smrpc"

	lhtypes "github.com/longhorn/go-common-libs/types"
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

func ServerCmd() *cli.Command {
	return &cli.Command{
		Name: "daemon",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "volume",
				Usage:    "The volume to export via the nfs server",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "data-engine",
				Usage:    "The volume data engine",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "encrypted",
				Usage:    "signals that a volume is encrypted",
				Sources:  cli.EnvVars("ENCRYPTED"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "passphrase",
				Usage:    "contains the encryption passphrase of the secret key provider",
				Sources:  cli.EnvVars("PASSPHRASE"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "cryptokeyprovider",
				Usage:    "the key provider: empty or 'secret' uses --passphrase, 'vault' fetches the key from Vault",
				Sources:  cli.EnvVars("CRYPTOKEYPROVIDER"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "encryption-secret-name",
				Usage:    "the Secret in Longhorn's namespace configuring a key-store provider",
				Sources:  cli.EnvVars("ENCRYPTION_SECRET_NAME"),
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "encryption-initialized",
				Usage:    "the device was formatted with LUKS before, so a blank-reading device is refused",
				Sources:  cli.EnvVars("ENCRYPTION_INITIALIZED"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "cryptokeycipher",
				Usage:    "contains the encryption algorithm in dm-crypt notation",
				Sources:  cli.EnvVars("CRYPTOKEYCIPHER"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "cryptokeyhash",
				Usage:    "contains the hash algorithm for the checksum resilience mode",
				Sources:  cli.EnvVars("CRYPTOKEYHASH"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "cryptokeysize",
				Usage:    "contains the encryption key size",
				Sources:  cli.EnvVars("CRYPTOKEYSIZE"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "cryptopbkdf",
				Usage:    "contains the Password-Based Key Derivation Function (PBKDF) algorithm for LUKS keyslot",
				Sources:  cli.EnvVars("CRYPTOPBKDF"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "cryptopbkdfiterations",
				Usage:    "contains the forced iteration count for the PBKDF algorithm (higher values increase security but reduce performance)",
				Sources:  cli.EnvVars("CRYPTOPBKDFFORCEITERATIONS"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "cryptopbkdfmemory",
				Usage:    "contains the memory cost parameter for the PBKDF algorithm in KiB (as used by cryptsetup)",
				Sources:  cli.EnvVars("CRYPTOPBKDFMEMORY"),
				Required: false,
			},
			&cli.StringFlag{
				Name:     "fs",
				Usage:    "the filesystem to use for the volume",
				Value:    "ext4",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "mount",
				Usage:    "allows for specifying additional mount options",
				Required: false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			vol := volume.Volume{
				Name:                       c.String("volume"),
				DataEngine:                 c.String("data-engine"),
				Encrypted:                  c.Bool("encrypted"),
				Passphrase:                 c.String("passphrase"),
				KeyProvider:                c.String("cryptokeyprovider"),
				EncryptionSecretName:       c.String("encryption-secret-name"),
				Initialized:                c.Bool("encryption-initialized"),
				CryptoKeyCipher:            c.String("cryptokeycipher"),
				CryptoKeyHash:              c.String("cryptokeyhash"),
				CryptoKeySize:              c.String("cryptokeysize"),
				CryptoPBKDF:                c.String("cryptopbkdf"),
				CryptoPBKDFForceIterations: c.String("cryptopbkdfiterations"),
				CryptoPBKDFMemory:          c.String("cryptopbkdfmemory"),
				FsType:                     c.String("fs"),
				MountOptions:               c.StringSlice("mount"),
			}

			if c.Bool("encrypted") && vol.KeyProviderName() == lhtypes.EncryptionKeyProviderSecret && len(vol.Passphrase) == 0 {
				logrus.Fatalf("Error starting share-manager missing passphrase for encrypted volume %v", vol.Name)
			}

			if err := start(vol); err != nil {
				logrus.Fatalf("Error running start command: %v.", err)
			}

			return nil
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
