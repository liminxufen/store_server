package driver

import (
	"context"
	//"testing"
	//"github.com/stretchr/testify/assert"
)

var (
	cmsDriver *CMSDriver
)

func setup() {
	DefaultTestScheme.Setup()
	DefaultTestMongoScheme.Setup()
	ct := context.Background()
	cmsDriver = &CMSDriver{
		ctx:     ct,
		MusicDB: DefaultTestScheme.DB(),
	}
}

func cleanup() {
	DefaultTestScheme.Cleanup()
	DefaultTestMongoScheme.Cleanup()
	if cmsDriver != nil {
		cmsDriver.MusicDB.Close()
	}
	cmsDriver = nil
}
