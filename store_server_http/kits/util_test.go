package kits

import (
	"bytes"
	"testing"

	"github.com/store_server/logger"
	"github.com/store_server/store_server_http/g"
	"github.com/stretchr/testify/assert"
)

var cfg = "ip_white_list: x.x.x.x|y.y.y.y|z.z.z.z"

func TestCheckIpHasPermission(t *testing.T) {
	data := bytes.NewBufferString(cfg).Bytes()
	err := g.ParseConfig(data)
	assert.NoError(t, err)
	logger.Entry().Info(g.Config().IpWhiteList)

	ok := CheckIpHasPermission("x.x.x.x")
	assert.Equal(t, ok, true)

	ok = CheckIpHasPermission("o.o.o.o")
	assert.Equal(t, ok, false)
	g.Config().IpWhiteList = ""
}
