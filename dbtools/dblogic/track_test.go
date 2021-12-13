package dblogic

import (
	"testing"
	"time"

	"github.com/store_server/dbtools/driver"
	"github.com/stretchr/testify/assert"

	m "github.com/store_server/dbtools/models"
)

var (
	tracksDriver *TracksDriver
)

func tkSetup() {
	driver.DefaultTestScheme.Setup()
	driver.DefaultTestMongoScheme.Setup()
	tracksDriver = &TracksDriver{
		&driver.CMSDriver{
			MusicDB:  driver.DefaultTestScheme.DB(),
			ImportDB: driver.DefaultTestScheme.DB(),
		},
		lock,
	}
}

func tkCleanup() {
	driver.DefaultTestScheme.Cleanup()
	driver.DefaultTestMongoScheme.Cleanup()
	if tracksDriver != nil {
		if tracksDriver.CMSDriver != nil {
			tracksDriver.CMSDriver.Close()
		}
		tracksDriver = nil
	}
}

func genTrackExample() *m.Track {
	track := &m.Track{
		FtrackId:     int64(1),
		FtrackName:   "中国人",
		FalbumId:     int64(1111111),
		Ftype:        -1,
		Flanguage:    0,
		Fsinger:      71,
		Fmovie:       "test_movie",
		Fsize:        100,
		Fduration:    32,
		FsingerId1:   71,
		FsingerId2:   0,
		FsingerId3:   0,
		FsingerId4:   0,
		Fprice1:      10,
		Fprice2:      0,
		Fprice3:      0,
		Fisrc:        "test_isrc",
		Fattribute1:  1,
		Fattribute2:  2,
		Fattribute3:  3,
		Fattribute4:  4,
		Fgenre:       7,
		FsingerAll:   "test_all",
		Flocation:    9,
		FtrackCId:    1,
		Flyric:       0,
		FportalLyric: 0,
		Fstatus:      -50,
		FgoSoso:      11,
		Fnote:        111,
		Fversion:     1,
		Fattribute5:  5,
		Fattribute6:  6,
		FtrackMid:    "1",
		FlinkMv:      10,
		FmediaId:     10,
		FlinkRing:    7,
		FvalidTime:   time.Now(),
		FuploadTime:  time.Now(),
		FmodifyTime:  time.Now(),
	}
	return track
}

func genTrackExtraOsExample() *m.TrackExtraOs {
	trackExtraOs := &m.TrackExtraOs{
		FtrackId:          int64(1),
		Fregion:           int64(1),
		FlocalName:        "xxxxxxxxuuuuuuuuuuuukkkkkkkkkkkk",
		FlocalCopyright:   1,
		FlocalStatus:      1,
		FlocalMovie:       "test_movie",
		FlocalFrom:        1,
		FactionTemplateId: 0,
		FmvId:             1,
		FcopyrightLimit:   1,
		FreplaceId:        1,
		FlocalIsrc:        "test_isrc",
		FlocalLabel:       "test_label",
		Fsupplier:         "test_supplier",
		FallSources:       "test_all_sources",
		FlocalOtherName:   "test_other_name",
		FlocalValidTime:   time.Now(),
	}
	return trackExtraOs
}

func TestInsertTrack(t *testing.T) {
	tkSetup()
	defer tkCleanup()
	track := genTrackExample()
	id, err := tracksDriver.InsertTrack(track)
	assert.NoError(t, err)
	assert.Equal(t, id, int64(1))

	data, err := tracksDriver.GetOneTrack(track.FtrackId)
	assert.NoError(t, err)
	assert.Equal(t, data.Fsize, 100)
	assert.Equal(t, data.FsingerAll, "test_all")
	assert.Equal(t, data.Fstatus, -50)
	assert.Equal(t, data.FtrackName, "中国人")
}

func TestInsertTrackExtraOs(t *testing.T) {
	tkSetup()
	defer tkCleanup()
	trackExtraOs := genTrackExtraOsExample()
	id, err := tracksDriver.InsertTrackExtraOs(trackExtraOs)
	assert.NoError(t, err)
	assert.Equal(t, id, int64(1))
}

func TestUpdateTrack(t *testing.T) {
	tkSetup()
	defer tkCleanup()
	id := int64(1)
	track := genTrackExample()
	track.FalbumId, track.Fsize = int64(233333), int64(200)
	idx, err := tracksDriver.UpdateTrack(id, track)

	assert.NoError(t, err)
	assert.Equal(t, idx, int64(1))
}

func TestGetTrackByCondition(t *testing.T) {
	tkSetup()
	defer tkCleanup()
	track := genTrackExample()
	track.FtrackId, track.FalbumId, track.Fmovie = int64(2), int64(20191108), "test_20191108_movie"

	id, err := tracksDriver.InsertTrack(track)
	assert.NoError(t, err)
	assert.Equal(t, id, int64(2))

	page, pagesize := 1, 10
	cons := map[string]interface{}{"Fisrc": "test_isrc"}
	tracks, total, err := tracksDriver.GetTrackByCondition(cons, page, pagesize)

	assert.NoError(t, err)
	assert.Equal(t, len(tracks), 2)
	assert.Equal(t, total, 2)
}

func TestUpdateTrackAttr(t *testing.T) {
	tkSetup()
	defer tkCleanup()
	track := genTrackExample()
	updateAttrs := map[string]interface{}{"Fstatus": 50}
	err := tracksDriver.UpdateTrackAttr(track, updateAttrs)

	assert.NoError(t, err)

	newTrack, err := tracksDriver.GetOneTrack(track.FtrackId)
	assert.NoError(t, err)
	assert.Equal(t, newTrack.Fstatus, 50)
}
