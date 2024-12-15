package utils

import (
	"time"
)

type VersionInfo struct {
	Version    string
	CommitHash string
	BuildTime  string
}

var (
	Version    string
	CommitHash string
	BuildTime  string

	lastRun time.Time
)

func GetVersion() *VersionInfo {
	return &VersionInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildTime:  BuildTime,
	}
}

func SetLastRun(t time.Time) {
	lastRun = t
}

func GetLastRun() time.Time {
	return lastRun
}
