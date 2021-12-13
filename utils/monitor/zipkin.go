package monitor

import (
	"context"
	//zipkin "github.com/openzipkin/zipkin-go"
	"sync"
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
