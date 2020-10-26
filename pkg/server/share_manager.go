package server

import (
	"golang.org/x/net/context"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/longhorn/longhorn-share-manager/pkg/server/nfs"
)

const waitBetweenChecks = time.Second * 5
const healthCheckInterval = time.Second * 10

type ShareManager struct {
	logger logrus.FieldLogger

	volume string

	context  context.Context
	shutdown context.CancelFunc

	nfsServer *nfs.Server
}

func NewShareManager(logger logrus.FieldLogger, volume string) (*ShareManager, error) {
	m := &ShareManager{logger: logger, volume: volume}
	m.context, m.shutdown = context.WithCancel(context.Background())

	nfsServer, err := nfs.NewDefaultServer(logger)
	if err != nil {
		return nil, err
	}
	m.nfsServer = nfsServer
	return m, nil
}

func (m *ShareManager) Run() error {
	for ; ; time.Sleep(waitBetweenChecks) {
		select {
		case <-m.context.Done():
			m.logger.Info("nfs server is shutting down")
			return nil
		default:
			log := m.logger.WithField("volume", m.volume)
			if !checkDeviceValid(m.volume) {
				log.Warn("waiting with nfs server start, volume is not attached")
				break
			}

			if err := mountVolume(m.volume); err != nil {
				log.Warn("waiting with nfs server start, failed to mount volume")
				break
			}

			log.Info("starting nfs server, volume is ready for export")

			// This blocks until server exits
			if err := m.nfsServer.Run(m.context); err != nil {
				m.logger.WithError(err).Error("nfs server exited with error")
			}

			// if the server is exiting, try to unmount before we terminate the container
			if err := unmountVolume(m.volume); err != nil {
				m.logger.WithError(err).Error("failed to unmount volume")
			}
		}
	}
}

func (m *ShareManager) runHealthCheck() {
	m.logger.WithField("volume", m.volume).Info("starting health check for volume")
	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	select {
	case <-m.context.Done():
		m.logger.Info("nfs server is shutting down")
		return
	case <-ticker.C:
		if !m.hasHealthyVolume() {
			m.logger.WithField("volume", m.volume).Error("volume health check failed, terminating")
			m.Shutdown()
			return
		}
	}
}

func (m *ShareManager) hasHealthyVolume() bool {
	err := exec.CommandContext(m.context, "ls", "/export").Run()
	return err == nil
}

func (m *ShareManager) Shutdown() {
	m.shutdown()
}
