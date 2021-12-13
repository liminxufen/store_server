package dblogic

import (
	"testing"

	"github.com/store_server/dbtools/driver"
	m "github.com/store_server/dbtools/models"
	"github.com/stretchr/testify/assert"
)

var (
	videosDriver *VideosDriver
)

func voSetup() {
	driver.DefaultTestScheme.Setup()
	driver.DefaultTestMongoScheme.Setup()
	videosDriver = &VideosDriver{
		&driver.CMSDriver{
			MusicDB:  driver.DefaultTestScheme.DB(),
			ImportDB: driver.DefaultTestScheme.DB(),
		},
		lock,
	}
}

func voCleanup() {
	driver.DefaultTestScheme.Cleanup()
	driver.DefaultTestMongoScheme.Cleanup()
	if videosDriver != nil {
		if videosDriver.CMSDriver != nil {
			videosDriver.CMSDriver.Close()
		}
		videosDriver = nil
	}
}

func genVideoExample() *m.Video {
	video := &m.Video{
		FregionId:     int64(1),
		Ftitle:        "test",
		Fstatus:       int64(0),
		FlocalFrom:    int64(1),
		Fsource:       int64(100),
		Fimage:        "test_image",
		Fvideo:        "test_video",
		Fuuid:         "test_uuid",
		Fmd5:          "test_md5",
		Fduration:     "3s",
		Fformat:       ".mp4",
		Fsize:         "128k",
		Fupc:          "test_upc",
		Fisrc:         "test_isrc",
		Fgrid:         "test_grid",
		FuploadStatus: int64(2),
		Fwatermark:    int64(1),
		Fcreator:      "erichli",
	}
	return video
}

func TestInsertVideo(t *testing.T) {
	voSetup()
	defer voCleanup()
	video := genVideoExample()
	id, err := videosDriver.InsertVideo(video)
	assert.NoError(t, err)
	assert.Equal(t, id, int64(1))

	data, err := videosDriver.GetOneVideo(video.Fid)
	assert.NoError(t, err)
	assert.Equal(t, data.Fsize, "128k")
	assert.Equal(t, data.Fcreator, "erichli")
	assert.Equal(t, data.Fuuid, "test_uuid")
	assert.Equal(t, data.FuploadStatus, int(2))
}
