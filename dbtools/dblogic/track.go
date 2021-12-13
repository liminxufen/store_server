package dblogic

import (
	"fmt"
	"sync"
	"time"

	"github.com/store_server/dbtools/driver"
	"github.com/store_server/logger"

	m "github.com/store_server/dbtools/models"
)

//JOOX CMS TRACK相关操作
type TracksDriver struct {
	*driver.CMSDriver
	*BaseDriver
	lock sync.RWMutex
}

func NewTracksDriver(cmsDriver *driver.CMSDriver) *TracksDriver {
	baseDriver := &BaseDriver{cmsDriver, cmsDriver.MusicDB}
	return &TracksDriver{cmsDriver, baseDriver, sync.RWMutex{}}
}

var (
	TkDriver *TracksDriver
)

/* ---------------------------- t_track ------------------------ */
func (td *TracksDriver) ExecRawQuerySql4Track(sql string, page, pagesize int64) ([]*m.Track, int64, error) { //原生query语句
	tracks := make([]*m.Track, 0)
	res, total, err := td.ExecRawQuerySql(sql, page, pagesize, &m.Track{})
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.Track); ok {
			tracks = append(tracks, t)
		}
	}
	//logger.Entry().Errorf("total: %v|res: %v|tracks: %v", total, res, tracks)
	return tracks, total, nil
}

func (td *TracksDriver) InsertOneTrack(track *m.Track) (int64, error) {
	current := m.TimeNormal{time.Now()}
	track.FvalidTime = current
	track.FuploadTime = current
	track.FmodifyTime = current
	track.FlastestModifyTime = current
	_, err := td.InsertWithModel(track)
	return track.FtrackId, err
}

func (td *TracksDriver) InsertTracks(tracks []*m.Track) (int64, error) {
	for _, track := range tracks {
		current := m.TimeNormal{time.Now()}
		track.FvalidTime = current
		track.FuploadTime = current
		track.FmodifyTime = current
		track.FlastestModifyTime = current
	}
	return td.InsertWithModel(tracks)
}

func (td *TracksDriver) UpdateOneTrack(id int64, track *m.Track) (int64, error) {
	current := m.TimeNormal{time.Now()}
	track.FmodifyTime = current
	track.FlastestModifyTime = current
	affected, err := td.UpdateWithModel(&m.Track{}, "Ftrack_id = ?", track, id)
	if err != nil {
		return track.FtrackId, err
	}
	if affected != 1 {
		return track.FtrackId, fmt.Errorf("update track[%d] error, affect %d raw", track.FtrackId, affected)
	}
	return track.FtrackId, nil
}

func (td *TracksDriver) GetOneTrack(id int64) (track *m.Track, err error) {
	track = &m.Track{}
	err = td.MusicDB.First(track, "Ftrack_id=?", id).Error
	return
}

func (td *TracksDriver) GetTracksByIds(ids []int64) ([]*m.Track, int64, error) {
	var tracks []*m.Track
	err := td.MusicDB.Where("Ftrack_id in (?)", ids).Find(&tracks).Error
	if err != nil {
		return nil, 0, err
	}
	total := int64(len(tracks))
	return tracks, total, nil
}

func (td *TracksDriver) GetTracksByCondition(conds map[string]interface{}, page,
	pagesize int64) ([]*m.Track, int64, error) {
	tracks := make([]*m.Track, 0)
	res, total, err := td.QueryWithModel(&m.Track{}, conds, page, pagesize)
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.Track); ok {
			tracks = append(tracks, t)
		}
	}
	return tracks, total, err
}

func (td *TracksDriver) UpdateTracksAttr(ids []int64, conds,
	updatesAttrs map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = td.UpdateWithModel(&m.Track{}, "Ftrack_id in (?)", updatesAttrs, ids)
	} else if len(conds) != 0 {
		affected, err = td.UpdateWithModel(&m.Track{}, conds, updatesAttrs)
	}
	if err != nil {
		return affected, err
	}
	if len(ids) > 0 {
		if affected != int64(len(ids)) {
			return affected, fmt.Errorf("update tracks count[%d] error, affect %d raw", len(ids), affected)
		}
	}
	return affected, nil
}

func (td *TracksDriver) DeleteOneTrack(id int64) (err error) {
	_, err = td.DeleteWithModelID(&m.Track{}, id)
	return err
}

func (td *TracksDriver) DeleteTracks(ids []int64, conds map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = td.DeleteWithModel(&m.Track{}, "Ftrack_id in (?)", ids)
	} else if len(conds) != 0 { //删除条件需严格把关
		logger.Entry().Debugf("delete tracks condition: %v", conds)
		affected, err = td.DeleteWithModel(&m.Track{}, conds)
	}
	if err != nil {
		return affected, err
	}
	if len(ids) > 0 {
		if affected != int64(len(ids)) {
			return affected, fmt.Errorf("delete tracks count[%d] error, affect %d raw", len(ids), affected)
		}
	}
	return affected, nil
}

func (td *TracksDriver) ExportAllTracks(sql string) ([][]interface{}, error) { //导出所有歌曲信息
	return td.ExportAllRecords(sql)
}

/* ---------------------------- t_track_extra_os ------------------------ */
func (td *TracksDriver) ExecRawQuerySql4TrackExtraOs(sql string, page,
	pagesize int64) ([]*m.TrackExtraOs, int64, error) { //原生query语句
	tracks := make([]*m.TrackExtraOs, 0)
	res, total, err := td.ExecRawQuerySql(sql, page, pagesize, &m.TrackExtraOs{})
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.TrackExtraOs); ok {
			tracks = append(tracks, t)
		}
	}
	return tracks, total, nil
}

func (td *TracksDriver) InsertOneTrackExtraOs(track *m.TrackExtraOs) (int64, error) {
	current := m.TimeNormal{time.Now()}
	track.FlocalValidTime = current
	track.FmodifyTime = current
	_, err := td.InsertWithModel(track)
	return track.FtrackId, err
}

func (td *TracksDriver) InsertTrackExtraOs(tracks []*m.TrackExtraOs) (int64, error) {
	for _, track := range tracks {
		current := m.TimeNormal{time.Now()}
		track.FmodifyTime = current
	}
	return td.InsertWithModel(tracks)
}

func (td *TracksDriver) UpdateTrackExtraOsAttr(ids []int64, conds,
	updatesAttrs map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = td.UpdateWithModel(&m.TrackExtraOs{}, "Ftrack_id in (?)", updatesAttrs, ids)
	} else if len(conds) != 0 {
		affected, err = td.UpdateWithModel(&m.TrackExtraOs{}, conds, updatesAttrs)
	}
	if err != nil {
		return affected, err
	}
	if len(ids) > 0 {
		if affected != int64(len(ids)) {
			return affected, fmt.Errorf("update track extra os count[%d] error, affect %d raw", len(ids), affected)
		}
	}
	return affected, nil
}

func (td *TracksDriver) GetOneTrackExtraOs(id, region int64) (track *m.TrackExtraOs, err error) {
	track = &m.TrackExtraOs{}
	err = td.MusicDB.First(track, "Ftrack_id=? and Fregion=?", id, region).Error
	return
}

func (td *TracksDriver) GetTrackExtraOsByCondition(conds map[string]interface{}, page,
	pagesize int64) ([]*m.TrackExtraOs, int64, error) {
	tracks := make([]*m.TrackExtraOs, 0)
	res, total, err := td.QueryWithModel(&m.TrackExtraOs{}, conds, page, pagesize)
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.TrackExtraOs); ok {
			tracks = append(tracks, t)
		}
	}
	return tracks, total, err
}

func (td *TracksDriver) DeleteTrackExtraOs(ids []int64, conds map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = td.DeleteWithModel(&m.TrackExtraOs{}, "Ftrack_id in (?)", ids)
	} else if len(conds) != 0 { //删除条件需严格把关
		logger.Entry().Debugf("delete track extra os condition: %v", conds)
		affected, err = td.DeleteWithModel(&m.TrackExtraOs{}, conds)
	}
	return affected, err
}

/* ---------------------------- track 相关join查询------------------------ */

func (td *TracksDriver) JoinQueryWithRawSql(sql string, page, pagesize int64) ([][]interface{}, error) {
	return td.JoinQueryWithSql(sql, page, pagesize)
}
