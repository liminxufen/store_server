package monitor

import (
	"context"
	"sync"

	client "github.com/influxdata/influxdb-client-go/v2"
)

//influx driver
type InfluxDriver struct {
	sync.RWMutex
	ctx      context.Context
	client   client.Client
	database string
}

func NewInfluxDriver(ctx context.Context, address string, opts ...string) (id *InfluxDriver, err error) {
	id = &InfluxDriver{ctx: ctx}
	//conf := client.Options{}
	if len(opts) > 1 {
	}
	if len(opts) > 2 {
		id.database = opts[2]
	}
	return
}

func (id *InfluxDriver) SetDatabase(database string) {
	if database == "" {
		return
	}
	id.database = database
}
