package volume

import (
	"fmt"
	"os"

	"k8s.io/kubernetes/pkg/volume/util/hostutil"
	utilexec "k8s.io/utils/exec"
	"k8s.io/utils/mount"
)

type Volume struct {
	Name         string
	Passphrase   string
	FsType       string
	MountOptions []string
}

func (v Volume) IsEncrypted() bool {
	return len(v.Passphrase) > 0
}

func GetDiskFormat(devicePath string) (string, error) {
	mounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: utilexec.New()}
	return mounter.GetDiskFormat(devicePath)
}

func CheckDeviceValid(devicePath string) bool {
	isDevice, err := hostutil.NewHostUtil().PathIsDevice(devicePath)
	return err == nil && isDevice
}

func CheckMountValid(mountPath string) bool {
	notMnt, err := mount.IsNotMountPoint(mount.New(""), mountPath)
	return err == nil && !notMnt
}

func MountVolume(devicePath, mountPath, fsType string, mountOptions []string) error {
	if !CheckDeviceValid(devicePath) {
		return fmt.Errorf("cannot mount device %v to %v invalid device", devicePath, mountPath)
	}

	if CheckMountValid(mountPath) {
		return nil
	}

	mounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: utilexec.New()}

	if exists, err := hostutil.NewHostUtil().PathExists(mountPath); !exists || err != nil {
		if err != nil {
			return err
		}

		if err := makeDir(mountPath); err != nil {
			return err
		}
	}

	return mounter.FormatAndMount(devicePath, mountPath, fsType, mountOptions)
}

func SetPermissions(mountPath string, mode os.FileMode) error {
	if !CheckMountValid(mountPath) {
		return fmt.Errorf("cannot set permissions %v for path %v invalid mount point", mode, mountPath)
	}

	return os.Chmod(mountPath, mode)
}

func UnmountVolume(mountPath string) error {
	mounter := mount.New("")
	return mount.CleanupMountPoint(mountPath, mounter, true)
}

// makeDir creates a new directory.
// If pathname already exists as a directory, no error is returned.
// If pathname already exists as a file, an error is returned.
func makeDir(pathname string) error {
	err := os.MkdirAll(pathname, os.FileMode(0777))
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}
