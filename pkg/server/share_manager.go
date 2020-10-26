package server

import (
	"fmt"
	"github.com/longhorn/longhorn-share-manager/pkg/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/longhorn/longhorn-share-manager/pkg/rpc"
	"github.com/longhorn/longhorn-share-manager/pkg/server/nfs"
	"github.com/longhorn/longhorn-share-manager/pkg/types"
	"github.com/longhorn/longhorn-share-manager/pkg/util/broadcaster"
)

const monitorInterval = 10 * time.Second
const waitTimeRestart = 1 * time.Second

type ShareManager struct {
	logger logrus.FieldLogger

	broadcaster *broadcaster.Broadcaster
	broadcastCh chan interface{}

	// TODO: replace this update channel with a queue
	shareUpdateCh chan types.Share
	shares        map[string]*types.Share
	lock          sync.RWMutex

	context  context.Context
	shutdown context.CancelFunc

	logFile string

	nfsServer *nfs.Server
}

func NewShareManager(logger logrus.FieldLogger, logFile string) (*ShareManager, error) {
	m := &ShareManager{
		logger: logger,

		broadcaster: &broadcaster.Broadcaster{},
		broadcastCh: make(chan interface{}),

		shareUpdateCh: make(chan types.Share),
		shares:        map[string]*types.Share{},
		lock:          sync.RWMutex{},

		logFile: logFile,
	}
	m.context, m.shutdown = context.WithCancel(context.Background())

	nfsServer, err := nfs.NewDefaultServer(logger)
	if err != nil {
		return nil, err
	}
	m.nfsServer = nfsServer

	// help to kickstart the broadcaster
	// TODO: no idea what the broadcaster does yet and why we are cancelling it
	c, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := m.broadcaster.Subscribe(c, m.broadcastConnector); err != nil {
		return nil, err
	}
	go m.runStreamUpdater()
	go m.runNFSServer()
	go m.runMonitor()
	return m, nil
}

func (m *ShareManager) runNFSServer() {
	for ; ; time.Sleep(waitTimeRestart) {
		select {
		case <-m.context.Done():
			m.logger.Info("nfs server is shutting down")
			return
		default:
			// before we start the nfs server, we need to verify that all devices for known exports are attached.
			// i.e the server crashed and we moved to a new node since there could be outstanding locks,
			// we need to reattach & mount all devices before starting up the server again,
			// this way the server can correctly start in grace mode.
			knownExports := m.nfsServer.GetExports()
			waitWithStart := false
			for vol := range knownExports {
				log := m.logger.WithField("volume", vol)
				if !m.checkDeviceValid(vol) {
					log.Warn("waiting with nfs server start volume is not attached")
					waitWithStart = true
				} else if !m.checkMountValid(vol) {
					log.Warn("waiting with nfs server start volume is not mounted")
					waitWithStart = true
				} else {
					log.Debug("volume is ready for export")
				}
			}

			if waitWithStart {
				break
			}

			// This blocks until server exits (presumably due to an error)
			if err := m.nfsServer.Run(m.context); err != nil {
				m.logger.Errorf("nfs server exited with err: %v", err)
			}

			// mark all shares as error
			updatedShares := map[string]types.Share{}
			func() {
				m.lock.Lock()
				defer m.lock.Unlock()
				for vol, share := range m.shares {
					getLoggerForShare(m.logger, share).Warn("share became invalid since nfs server exited")
					share.State = types.ShareStateError
					share.Error = "nfs server terminated"
					updatedShares[vol] = *share
				}
			}()

			// update monitor with new share state
			for vol := range updatedShares {
				m.shareUpdateCh <- updatedShares[vol]
			}
		}
	}
}

func (m *ShareManager) runStreamUpdater() {
	done := false
	for {
		select {
		case <-m.context.Done():
			m.logger.Info("Share manager is shutting down")
			done = true
			break
		case share := <-m.shareUpdateCh:
			rsp := client.ShareToRPC(&share)
			m.broadcastCh <- interface{}(rsp)
		}
		if done {
			break
		}
	}
}

func (m *ShareManager) runMonitor() {
	timer := time.NewTicker(monitorInterval)
	defer timer.Stop()

	for {
		select {
		case <-m.context.Done():
			m.logger.Infof("Share manager is shutting down")
			return
		case <-timer.C:
			// periodically check all valid shares
			// to make sure that the disks are functioning
			updatedShares := map[string]types.Share{}
			func() {
				m.lock.Lock()
				defer m.lock.Unlock()
				for vol, share := range m.shares {
					if share.State == types.ShareStateReady && !m.checkShareValid(share) {
						getLoggerForShare(m.logger, share).Warn("share became invalid")
						share.State = types.ShareStateError
						share.State = "share is not valid"
						updatedShares[vol] = *share
					}
				}
			}()

			// TODO: we need to add a queue instead of this single update channel
			for vol := range updatedShares {
				m.shareUpdateCh <- updatedShares[vol]
			}
		}
	}
}

func (m *ShareManager) checkShareValid(share *types.Share) bool {
	return m.checkDeviceValid(share.Volume) &&
		m.checkMountValid(share.Volume) &&
		m.checkExportValid(share.Volume)
}

func (m *ShareManager) Shutdown() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.shutdown()
}

func (m *ShareManager) broadcastConnector() (chan interface{}, error) {
	return m.broadcastCh, nil
}

func (m *ShareManager) Subscribe() (<-chan interface{}, error) {
	return m.broadcaster.Subscribe(context.TODO(), m.broadcastConnector)
}

func getLoggerForShare(logger logrus.FieldLogger, share *types.Share) logrus.FieldLogger {
	return logger.WithField("volume", share.Volume).WithField("exportID", share.ExportID)
}

func (m *ShareManager) ShareCreate(ctx context.Context, req *rpc.ShareCreateRequest) (*rpc.ShareResponse, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	volume := req.Volume

	// new share
	share := m.shares[volume]
	if share == nil {
		share = &types.Share{
			Volume: volume,
			State:  types.ShareStatePending,
		}
		m.shares[volume] = share
	}
	log := getLoggerForShare(m.logger, share)

	if !m.checkDeviceValid(volume) {
		share.State = types.ShareStateError
		share.Error = "invalid block device"

		log.Errorf("cannot share volume invalid block device")
		m.shareUpdateCh <- *share
		return client.ShareToRPC(share), nil
	}

	if !m.checkMountValid(volume) {
		if err := m.mount(volume); err != nil {
			share.State = types.ShareStateError
			share.Error = "failed to mount volume"

			log.WithError(err).Error("cannot share volume failed to mount")
			m.shareUpdateCh <- *share
			return client.ShareToRPC(share), nil
		}
	}

	exportID, err := m.nfsServer.CreateExport(volume)
	if err != nil {
		share.State = types.ShareStateError
		share.Error = "failed to export volume"

		log.WithError(err).Error("cannot share volume failed to export")
		m.shareUpdateCh <- *share
		return client.ShareToRPC(share), nil
	}

	// share is ready
	share.ExportID = exportID
	share.State = types.ShareStateReady
	share.Error = ""

	m.shareUpdateCh <- *share
	return client.ShareToRPC(share), nil
}

func (m *ShareManager) ShareDelete(ctx context.Context, req *rpc.ShareDeleteRequest) (*rpc.ShareResponse, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	volume := req.Volume

	// not present share
	share := m.shares[volume]
	if share == nil {
		share = &types.Share{
			Volume: volume,
			State:  types.ShareStateDeleted,
		}
		m.shares[volume] = share
	}

	// TODO: remove from nfs server
	exportID, err := m.nfsServer.DeleteExport(volume)
	if err != nil {
		m.logger.WithField("volume", volume).WithError(err).Error("failed to delete export from nfs server")
		return nil, err
	}

	if err := m.unmount(volume); err != nil {
		m.logger.WithField("volume", volume).WithError(err).Error("failed to unmount volume")
		return nil, err
	}

	// share deleted
	share.ExportID = exportID
	share.State = types.ShareStateDeleted
	share.Error = ""
	return client.ShareToRPC(share), nil
}

func (m *ShareManager) ShareGet(ctx context.Context, req *rpc.ShareGetRequest) (*rpc.ShareResponse, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	share, ok := m.shares[req.Volume]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown share %v", req.Volume)
	}
	return client.ShareToRPC(share), nil
}

func (m *ShareManager) ShareList(ctx context.Context, req *rpc.ShareListRequest) (*rpc.ShareListResponse, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return client.ShareMapToRPC(m.shares), nil
}

func (m *ShareManager) ShareWatch(req *rpc.ShareWatchRequest, srv rpc.ShareManagerService_ShareWatchServer) error {
	responseChan, err := m.Subscribe()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			m.logger.Errorf("share watch errored out: %v", err)
			return
		}
		m.logger.Debug("share watch ended successfully")
	}()

	m.logger.Debugf("streaming share updates")
	for resp := range responseChan {
		share, ok := resp.(*rpc.ShareResponse)
		if !ok {
			return fmt.Errorf("BUG: cannot get ShareResponse from channel")
		}
		if err := srv.Send(share); err != nil {
			return err
		}
	}

	return nil
}
