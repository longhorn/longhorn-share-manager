package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	coordinationv1 "k8s.io/api/coordination/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coordinationv1client "k8s.io/client-go/kubernetes/typed/coordination/v1"
	"k8s.io/client-go/rest"
	mount "k8s.io/mount-utils"

	"github.com/longhorn/longhorn-share-manager/pkg/crypto"
	"github.com/longhorn/longhorn-share-manager/pkg/server/nfs"
	"github.com/longhorn/longhorn-share-manager/pkg/types"
	"github.com/longhorn/longhorn-share-manager/pkg/volume"

	commonKubernetes "github.com/longhorn/go-common-libs/kubernetes"
)

const waitBetweenChecks = time.Second * 5
const leaseRenewInterval = time.Second * 3
const healthCheckInterval = time.Second * 10
const configPath = "/tmp/vfs.conf"
const defaultNamespace = "longhorn-system" // backward compatibility namespace
const shareManagerPrefix = "share-manager-"

const EnvKeyFastFailover = "FAST_FAILOVER"
const EnvKeyLeaseLifetime = "LEASE_LIFETIME"
const EnvKeyGracePeriod = "GRACE_PERIOD"
const EnvKeyFormatOptions = "FS_FORMAT_OPTIONS"
const defaultLeaseLifetime = 60
const defaultGracePeriod = 90

const (
	UnhealthyErr = "UNHEALTHY: volume with mount path %v is unhealthy"
	ReadOnlyErr  = "READONLY: volume with mount path %v is read only"
)

type ShareManager struct {
	logger logrus.FieldLogger

	volume        volume.Volume
	shareExported bool

	context  context.Context
	shutdown context.CancelFunc

	enableFastFailover bool
	leaseHolder        string
	leaseClient        coordinationv1client.LeasesGetter
	lease              *coordinationv1.Lease

	nfsServer *nfs.Server

	namespace string
	podName   string
}

func NewShareManager(logger logrus.FieldLogger, volume volume.Volume) (*ShareManager, error) {
	m := &ShareManager{
		volume: volume,
		logger: logger.WithField("volume", volume.Name).WithField("encrypted", volume.IsEncrypted()),
	}
	m.context, m.shutdown = context.WithCancel(context.Background())

	m.enableFastFailover = m.getEnvAsBool(EnvKeyFastFailover, false)
	leaseLifetime := m.getEnvAsInt(EnvKeyLeaseLifetime, defaultLeaseLifetime)
	gracePeriod := m.getEnvAsInt(EnvKeyGracePeriod, defaultGracePeriod)

	// get pod namespace from env
	namespace := os.Getenv(types.EnvPodNamespace)
	if namespace == "" {
		m.logger.Warnf("Cannot detect pod namespace, environment variable %v is missing, using default namespace", types.EnvPodNamespace)
		namespace = defaultNamespace
	}

	m.namespace = namespace

	// get pod name from env
	podName := os.Getenv(types.EnvPodName)
	if podName == "" {
		m.logger.Warnf("Cannot detect pod name, environment variable %v is missing, using generated name", types.EnvPodName)
		podName = shareManagerPrefix + m.volume.Name
	}

	m.podName = podName

	if m.enableFastFailover {
		kubeclientset, err := m.NewKubeClient()
		if err != nil {
			m.logger.WithError(err).Error("Failed to make lease client for fast failover")
			return nil, err
		}

		// Use the clientset to get the node name of the share-manager pod
		// and store for use as lease holder.
		pod, err := kubeclientset.CoreV1().Pods(m.namespace).Get(m.context, m.podName, metav1.GetOptions{})
		if err != nil {
			m.logger.WithError(err).Warn("Failed to get share-manager pod specification from API for fast failover")
			return nil, err
		}
		m.leaseHolder = pod.Spec.NodeName
		m.leaseClient = kubeclientset.CoordinationV1()
	}

	nfsServer, err := nfs.NewServer(logger, configPath, types.ExportPath, volume.Name, leaseLifetime, gracePeriod)
	if err != nil {
		return nil, err
	}
	m.nfsServer = nfsServer
	return m, nil
}

func (m *ShareManager) Run() error {
	vol := m.volume
	mountPath := types.GetMountPath(vol.Name)
	devicePath := types.GetVolumeDevicePath(vol.Name, vol.DataEngine, false)

	defer func() {
		// if the server is exiting, try to unmount & teardown device before we terminate the container
		if err := volume.UnmountVolume(mountPath); err != nil {
			m.logger.WithError(err).Error("Failed to unmount volume")
		}

		if err := m.tearDownDevice(vol); err != nil {
			m.logger.WithError(err).Error("Failed to tear down volume")
		}

		m.Shutdown()
	}()

	// Check every waitBetweenChecks for volume attachment. Then run server process once and wait for completion.
	for ; ; time.Sleep(waitBetweenChecks) {
		select {
		case <-m.context.Done():
			m.logger.Info("NFS server is shutting down")
			return nil
		default:
			if !volume.CheckDeviceValid(devicePath) {
				m.logger.Warn("Waiting with nfs server start, volume is not attached")
				break
			}

			devicePath, err := m.setupDevice(vol, devicePath)
			if err != nil {
				return err
			}

			if err := m.MountVolume(vol, devicePath, mountPath); err != nil {
				m.logger.WithError(err).Warn("Failed to mount volume")
				return err
			}

			if err := m.resizeVolume(devicePath, mountPath); err != nil {
				m.logger.WithError(err).Warn("Failed to resize volume after mount")
				return err
			}

			if err := volume.SetPermissions(mountPath, 0777); err != nil {
				m.logger.WithError(err).Error("Failed to set permissions for volume")
				return err
			}

			m.logger.Info("Starting nfs server, volume is ready for export")

			if m.enableFastFailover {
				if err = m.takeLease(); err != nil {
					m.logger.WithError(err).Error("Failed to take lease for fast failovr")
					return err
				}
				go m.runLeaseRenew()
			}

			go m.runHealthCheck()

			if _, err := m.nfsServer.CreateExport(vol.Name); err != nil {
				m.logger.WithError(err).Error("Failed to create nfs export")
				return err
			}

			m.SetShareExported(true)

			// This blocks until server exits
			if err := m.nfsServer.Run(m.context); err != nil {
				m.logger.WithError(err).Error("NFS server exited with error")
			}
			return err
		}
	}
}

// setupDevice will return a path where the device file can be found
// for encrypted volumes, it will try formatting the volume on first use
// then open it and expose a crypto device at the returned path
func (m *ShareManager) setupDevice(vol volume.Volume, devicePath string) (string, error) {
	diskFormat, err := volume.GetDiskFormat(devicePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to determine filesystem format of volume error")
	}
	m.logger.Infof("Volume %v device %v contains filesystem of format %v", vol.Name, devicePath, diskFormat)

	if vol.IsEncrypted() || diskFormat == "crypto_LUKS" {
		if vol.Passphrase == "" {
			return "", fmt.Errorf("missing passphrase for encrypted volume %v", vol.Name)
		}

		// initial setup of longhorn device for crypto
		if diskFormat == "" {
			m.logger.Info("Encrypting new volume before first use")
			if err := crypto.EncryptVolume(devicePath, vol.Passphrase, vol.CryptoKeyCipher, vol.CryptoKeyHash, vol.CryptoKeySize, vol.CryptoPBKDF); err != nil {
				return "", errors.Wrapf(err, "failed to encrypt volume %v", vol.Name)
			}
		}

		cryptoDevice := types.GetVolumeDevicePath(vol.Name, vol.DataEngine, true)
		m.logger.Infof("Volume %s requires crypto device %s", vol.Name, cryptoDevice)
		if err := crypto.OpenVolume(vol.Name, vol.DataEngine, devicePath, vol.Passphrase); err != nil {
			m.logger.WithError(err).Error("Failed to open encrypted volume")
			return "", err
		}

		// update the device path to point to the new crypto device
		return cryptoDevice, nil
	}

	return devicePath, nil
}

func (m *ShareManager) tearDownDevice(vol volume.Volume) error {
	// close any matching crypto device for this volume
	cryptoDevice := types.GetVolumeDevicePath(vol.Name, vol.DataEngine, true)
	if isOpen, err := crypto.IsDeviceOpen(cryptoDevice); err != nil {
		return err
	} else if isOpen {
		m.logger.Infof("Volume %s has active crypto device %s", vol.Name, cryptoDevice)
		if err := crypto.CloseVolume(vol.Name, vol.DataEngine); err != nil {
			return err
		}
		m.logger.Infof("Volume %s closed active crypto device %s", vol.Name, cryptoDevice)
	}

	return nil
}

func (m *ShareManager) MountVolume(vol volume.Volume, devicePath, mountPath string) error {
	fsType := vol.FsType
	mountOptions := vol.MountOptions
	formatOptions := m.getFormatOptions()

	// https://github.com/longhorn/longhorn/issues/2991
	// pre v1.2 we ignored the fsType and always formatted as ext4
	// after v1.2 we include the user specified fsType to be able to
	// mount priorly created volumes we need to switch to the existing fsType
	diskFormat, err := volume.GetDiskFormat(devicePath)
	if err != nil {
		m.logger.WithError(err).Error("Failed to evaluate disk format")
		return err
	}

	// `unknown data, probably partitions` is used when the disk contains a partition table
	if diskFormat != "" && !strings.Contains(diskFormat, "unknown data") && fsType != diskFormat {
		m.logger.Warnf("Disk is already formatted to %v but user requested fs is %v using existing device fs type for mount", diskFormat, fsType)
		fsType = diskFormat
	}

	return volume.MountVolume(devicePath, mountPath, fsType, mountOptions, formatOptions)
}

func (m *ShareManager) getFormatOptions() []string {
	env := os.Getenv(EnvKeyFormatOptions)
	if env == "" {
		return nil
	}

	return strings.Split(env, ":")
}

func (m *ShareManager) resizeVolume(devicePath, mountPath string) error {
	// Note that we don't need 'cryptsetup resize' here.  The crypto 'open' will have done so if necessary.
	if resized, err := volume.ResizeVolume(devicePath, mountPath); err != nil {
		m.logger.WithError(err).Error("Failed to resize filesystem for volume")
		return err
	} else if resized {
		m.logger.Info("Resized filesystem for volume after mount")
	}

	return nil
}

func (m *ShareManager) getEnvAsInt(key string, defaultVal int) int {
	env := os.Getenv(key)
	if env == "" {
		m.logger.Warnf("Failed to get expected environment variable, env %v wasn't set, defaulting to %v", key, defaultVal)
		return defaultVal
	}
	value, err := strconv.Atoi(env)
	if err != nil {
		m.logger.Warnf("Failed to convert environment variable, %v, value %v, to an int, defaulting to %v", key, env, defaultVal)
		return defaultVal
	}
	return value
}

func (m *ShareManager) getEnvAsBool(key string, defaultVal bool) bool {
	env := os.Getenv(key)
	if env == "" {
		m.logger.Warnf("Failed to get expected environment variable, env %v wasn't set, defaulting to %v", key, defaultVal)
		return defaultVal
	}
	value, err := strconv.ParseBool(env)
	if err != nil {
		m.logger.Warnf("Failed to convert environment variable, %v, value %v, to an int, defaulting to %v", key, env, defaultVal)
		return defaultVal
	}
	return value
}

func (m *ShareManager) NewKubeClient() (*kubernetes.Clientset, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (m *ShareManager) takeLease() error {
	if m.leaseClient == nil {
		return fmt.Errorf("kubernetes API client is unset")
	}

	lease, err := m.leaseClient.Leases(m.namespace).Get(m.context, m.volume.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	m.lease = lease

	now := time.Now()
	currentHolder := *m.lease.Spec.HolderIdentity
	m.logger.Infof("Updating lease holderIdentity from %v to %v", currentHolder, m.leaseHolder)

	*m.lease.Spec.HolderIdentity = m.leaseHolder
	*m.lease.Spec.LeaseTransitions = *m.lease.Spec.LeaseTransitions + 1
	m.lease.Spec.AcquireTime = &metav1.MicroTime{Time: now}
	m.lease.Spec.RenewTime = &metav1.MicroTime{Time: now}

	lease, err = m.leaseClient.Leases(m.namespace).Update(m.context, m.lease, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	m.lease = lease
	m.logger.Infof("Took lease for volume %v as holder %v", m.volume.Name, m.leaseHolder)
	return nil
}

func (m *ShareManager) renewLease() error {
	m.lease.Spec.RenewTime = &metav1.MicroTime{Time: time.Now()}
	lease, err := m.leaseClient.Leases(m.namespace).Update(m.context, m.lease, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	m.lease = lease
	return nil
}

func (m *ShareManager) runLeaseRenew() {
	m.logger.Infof("Starting lease renewal for volume mounted at: %v", types.GetMountPath(m.volume.Name))
	ticker := time.NewTicker(leaseRenewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.context.Done():
			m.logger.Info("NFS lease renewal is ending")
			return
		case <-ticker.C:
			if err := m.renewLease(); err != nil {
				m.logger.Warn("Failed to renew share-manager lease - expect to be terminated.")
			}
		}
	}
}

func (m *ShareManager) runHealthCheck() {
	m.logger.Infof("Starting health check for volume mounted at: %v", types.GetMountPath(m.volume.Name))
	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.context.Done():
			m.logger.Info("NFS server is shutting down")
			return
		case <-ticker.C:
			if err := m.hasHealthyVolume(); err != nil {
				if strings.Contains(err.Error(), "UNHEALTHY") {
					m.logger.WithError(err).Error("Terminating")
					m.Shutdown()
					return
				} else if strings.Contains(err.Error(), "READONLY") {
					m.logger.WithError(err).Warn("Recovering read only volume")
					if err := m.recoverReadOnlyVolume(); err != nil {
						m.logger.WithError(err).Error("Volume is unable to recover by remounting, terminating")
						m.Shutdown()
						return
					}
				}
			}
		}
	}
}

func (m *ShareManager) hasHealthyVolume() error {
	mountPath := types.GetMountPath(m.volume.Name)
	if err := exec.CommandContext(m.context, "ls", mountPath).Run(); err != nil {
		return fmt.Errorf(UnhealthyErr, mountPath)
	}

	mounter := mount.New("")
	mountPoints, _ := mounter.List()
	for _, mp := range mountPoints {
		if mp.Path == mountPath && commonKubernetes.IsMountPointReadOnly(mp) {
			return fmt.Errorf(ReadOnlyErr, mountPath)
		}
	}
	return nil
}

func (m *ShareManager) recoverReadOnlyVolume() error {
	mountPath := types.GetMountPath(m.volume.Name)

	cmd := exec.CommandContext(m.context, "mount", "-o", "remount,rw", mountPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "remount failed with output: %s", out)
	}

	return nil
}

func (m *ShareManager) GetVolume() volume.Volume {
	return m.volume
}

func (m *ShareManager) SetShareExported(val bool) {
	m.shareExported = val
}

func (m *ShareManager) ShareIsExported() bool {
	return m.shareExported
}

func (m *ShareManager) Shutdown() {
	m.shutdown()
}
