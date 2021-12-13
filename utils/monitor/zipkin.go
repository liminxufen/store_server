package monitor

import (
	"context"
	"sync"
	//zipkin "github.com/openzipkin/zipkin-go"
)

//zipkin driver
type ZipkinDriver struct {
	ctx context.Context
	sync.RWMutex
}

func NewZipkinDriver(ctx context.Context) (zd *ZipkinDriver, err error) {
	zd = &ZipkinDriver{ctx: ctx}
	return
}
