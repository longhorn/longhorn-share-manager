package crypto

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	lhns "github.com/longhorn/go-common-libs/ns"
	lhtypes "github.com/longhorn/go-common-libs/types"

	"github.com/longhorn/longhorn-share-manager/pkg/types"
)

// EncryptVolume encrypts provided device with LUKS.
func EncryptVolume(devicePath, passphrase, keyCipher, keyHash, keySize, pbkdf string) error {
	namespaces := []lhtypes.Namespace{lhtypes.NamespaceMnt, lhtypes.NamespaceIpc}
	nsexec, err := lhns.NewNamespaceExecutor(lhtypes.ProcessNone, lhtypes.HostProcDirectory, namespaces)
	if err != nil {
		return err
	}

	logrus.Debugf("Encrypting device %s with LUKS", devicePath)
	if _, err := nsexec.LuksFormat(devicePath, passphrase, keyCipher, keyHash, keySize, pbkdf, lhtypes.LuksTimeout); err != nil {
		return errors.Wrapf(err, "failed to encrypt device %s with LUKS", devicePath)
	}
	return nil
}

// OpenVolume opens volume so that it can be used by the client.
// devicePath is the path of the volume on the host that will be opened for instance '/dev/longhorn/volume1'
func OpenVolume(volume, dataEngine, devicePath, passphrase string) error {
	devPath := types.GetVolumeDevicePath(volume, dataEngine, true)
	if isOpen, _ := IsDeviceOpen(devPath); isOpen {
		logrus.Debugf("Device %s is already opened at %s", devicePath, devPath)
		return nil
	}

	namespaces := []lhtypes.Namespace{lhtypes.NamespaceMnt, lhtypes.NamespaceIpc}
	nsexec, err := lhns.NewNamespaceExecutor(lhtypes.ProcessNone, lhtypes.HostProcDirectory, namespaces)
	if err != nil {
		return err
	}

	encryptedDevName := types.GetEncryptVolumeName(volume, dataEngine)
	logrus.Debugf("Opening device %s with LUKS on %s", devicePath, encryptedDevName)
	_, err = nsexec.LuksOpen(encryptedDevName, devicePath, passphrase, lhtypes.LuksTimeout)
	if err != nil {
		logrus.WithError(err).Warnf("Failed to open LUKS device %s to %s", devicePath, encryptedDevName)
	}
	return err
}

// CloseVolume closes encrypted volume so it can be detached.
func CloseVolume(volume, dataEngine string) error {
	namespaces := []lhtypes.Namespace{lhtypes.NamespaceMnt, lhtypes.NamespaceIpc}
	nsexec, err := lhns.NewNamespaceExecutor(lhtypes.ProcessNone, lhtypes.HostProcDirectory, namespaces)
	if err != nil {
		return err
	}

	deviceName := types.GetEncryptVolumeName(volume, dataEngine)
	logrus.Debugf("Closing LUKS device %s", deviceName)
	_, err = nsexec.LuksClose(deviceName, lhtypes.LuksTimeout)
	return err
}

func ResizeEncryptoDevice(volume, dataEngine, passphrase string) error {
	// devPath is the full path of the encrypted device on the host that will be resized
	devPath := types.GetVolumeDevicePath(volume, dataEngine, true)
	if isOpen, err := IsDeviceOpen(devPath); err != nil {
		return err
	} else if !isOpen {
		return fmt.Errorf("volume %v encrypto device %s is closed for resizing", volume, devPath)
	}

	namespaces := []lhtypes.Namespace{lhtypes.NamespaceMnt, lhtypes.NamespaceIpc}
	nsexec, err := lhns.NewNamespaceExecutor(lhtypes.ProcessNone, lhtypes.HostProcDirectory, namespaces)
	if err != nil {
		return err
	}

	_, err = nsexec.LuksResize(types.GetEncryptVolumeName(volume, dataEngine), passphrase, lhtypes.LuksTimeout)
	return err
}

// IsDeviceOpen determines if encrypted device is already open.
func IsDeviceOpen(device string) (bool, error) {
	_, mappedFile, err := DeviceEncryptionStatus(device)
	return mappedFile != "", err
}

// DeviceEncryptionStatus looks to identify if the passed device is a LUKS mapping
// and if so what the device is and the mapper name as used by LUKS.
// If not, just returns the original device and an empty string.
func DeviceEncryptionStatus(devicePath string) (mappedDevice, mapper string, err error) {
	if !strings.HasPrefix(devicePath, types.MapperDevPath) {
		return devicePath, "", nil
	}

	namespaces := []lhtypes.Namespace{lhtypes.NamespaceMnt, lhtypes.NamespaceIpc}
	nsexec, err := lhns.NewNamespaceExecutor(lhtypes.ProcessNone, lhtypes.HostProcDirectory, namespaces)
	if err != nil {
		return devicePath, "", err
	}

	// Check the mapper device using `cryptsetup status`. Sample output:
	//   /dev/mapper/pvc-e2a1c50a-9409-4afd-9e7d-3c1f5a9afa7f is active and is in use.
	//     type:    LUKS2
	//     cipher:  aes-xts-plain64
	//     keysize: 256 bits
	//     key location: keyring
	//     device:  /dev/sda
	//     sector size:  512
	//     offset:  32768 sectors
	//     size:    110592 sectors
	//     mode:    read/write
	volume := strings.TrimPrefix(devicePath, types.MapperDevPath+"/")
	stdout, err := nsexec.LuksStatus(volume, lhtypes.LuksTimeout)
	if err != nil {
		logrus.WithError(err).Debugf("Device %s is not an active LUKS device", devicePath)
		return devicePath, "", nil
	}
	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 1 {
		return "", "", fmt.Errorf("device encryption status returned no stdout for %s", devicePath)
	}

	// There are two possible cases to an encrypted mapper device:
	// - "/path/to/device is active.": the mapper device is activated by `cryptsetup luksOpen`
	// - "/path/to/device is active and is in use.": the activated device is mounted or exposed (e.g., NFS)
	if !strings.Contains(lines[0], " is active") {
		// Implies this is not a LUKS device
		return devicePath, "", nil
	}
	for i := 1; i < len(lines); i++ {
		kv := strings.SplitN(strings.TrimSpace(lines[i]), ":", 2)
		if len(kv) < 1 {
			return "", "", fmt.Errorf("device encryption status output for %s is badly formatted: %s",
				devicePath, lines[i])
		}
		if strings.Compare(kv[0], "device") == 0 {
			return strings.TrimSpace(kv[1]), volume, nil
		}
	}
	// Identified as LUKS, but failed to identify a mapped device
	return "", "", fmt.Errorf("mapped device not found in path %s", devicePath)
}
