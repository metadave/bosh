package platform

import (
	boshstats "bosh/platform/stats"
	boshsettings "bosh/settings"
)

type dummyPlatform struct {
}

func newDummyPlatform() (p dummyPlatform) {
	return
}

func (p dummyPlatform) SetupSsh(publicKey, username string) (err error) {
	return
}

func (p dummyPlatform) SetupDhcp(networks boshsettings.Networks) (err error) {
	return
}

func (p dummyPlatform) SetupEphemeralDiskWithPath(devicePath, mountPoint string) (err error) {
	return
}

func (p dummyPlatform) GetStatsCollector() (collector boshstats.StatsCollector) {
	return boshstats.NewDummyStatsCollector()
}