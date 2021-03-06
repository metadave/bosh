package action

import (
	"errors"

	bosherr "bosh/errors"
	boshsettings "bosh/settings"
	boshdirs "bosh/settings/directories"
)

type diskMounter interface {
	MountPersistentDisk(volumeID string, mountPoint string) error
}

type mountPoints interface {
	IsMountPoint(string) (bool, error)
}

type MountDiskAction struct {
	settings    boshsettings.Service
	diskMounter diskMounter
	mountPoints mountPoints
	dirProvider boshdirs.DirectoriesProvider
}

func NewMountDisk(
	settings boshsettings.Service,
	diskMounter diskMounter,
	mountPoints mountPoints,
	dirProvider boshdirs.DirectoriesProvider,
) (mountDisk MountDiskAction) {
	mountDisk.settings = settings
	mountDisk.diskMounter = diskMounter
	mountDisk.mountPoints = mountPoints
	mountDisk.dirProvider = dirProvider
	return
}

func (a MountDiskAction) IsAsynchronous() bool {
	return true
}

func (a MountDiskAction) IsPersistent() bool {
	return false
}

func (a MountDiskAction) Run(diskCid string) (value interface{}, err error) {
	err = a.settings.LoadSettings()
	if err != nil {
		err = bosherr.WrapError(err, "Refreshing the settings")
		return
	}

	disksSettings := a.settings.GetDisks()
	devicePath, found := disksSettings.Persistent[diskCid]
	if !found {
		err = bosherr.New("Persistent disk with volume id '%s' could not be found", diskCid)
		return
	}

	mountPoint := a.dirProvider.StoreDir()

	isMountPoint, err := a.mountPoints.IsMountPoint(mountPoint)
	if err != nil {
		err = bosherr.WrapError(err, "Checking mount point")
		return
	}
	if isMountPoint {
		mountPoint = a.dirProvider.StoreMigrationDir()
	}

	err = a.diskMounter.MountPersistentDisk(devicePath, mountPoint)
	if err != nil {
		err = bosherr.WrapError(err, "Mounting persistent disk")
		return
	}

	value = make(map[string]string)
	return
}

func (a MountDiskAction) Resume() (interface{}, error) {
	return nil, errors.New("not supported")
}
