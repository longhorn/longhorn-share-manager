package server

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/volume/util/hostutil"
)

const devPath = "/host/dev/longhorn/"
const exportPath = "/tmp/"

func (m *ShareManager) checkDeviceValid(volume string) bool {
	isDevice, err := hostutil.NewHostUtil().PathIsDevice(filepath.Join(devPath, volume))
	return err == nil && isDevice
}

func (m *ShareManager) checkMountValid(volume string) bool {
	notMnt, err := mount.IsNotMountPoint(mount.New(""), filepath.Join(exportPath, volume))
	return err == nil && !notMnt
}

func (m *ShareManager) checkExportValid(volume string) bool {
	exportID := m.nfsServer.GetExport(volume)
	return exportID != 0
}

func (m *ShareManager) mount(volume string) error {
	if !m.checkDeviceValid(volume) {
		return fmt.Errorf("cannot mount volume %v invalid device", volume)
	}
	if m.checkMountValid(volume) {
		return nil
	}

	devicePath := filepath.Join(devPath, volume)
	mountPath := filepath.Join(exportPath, volume)
	mounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOSExec()}

	if exists, err := hostutil.NewHostUtil().PathExists(mountPath); !exists || err != nil {
		if err != nil {
			return err
		}

		if err := makeDir(mountPath); err != nil {
			return err
		}
	}

	return mounter.FormatAndMount(devicePath, mountPath, "ext4", nil)
}

func (m *ShareManager) unmount(volume string) error {
	mountPath := path.Join(exportPath, volume)
	mounter := mount.New("")
	return mount.CleanupMountPoint(mountPath, mounter, true)
}

// makeDir creates a new directory.
// If pathname already exists as a directory, no error is returned.
// If pathname already exists as a file, an error is returned.
func makeDir(pathname string) error {
	err := os.MkdirAll(pathname, os.FileMode(0755))
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}
