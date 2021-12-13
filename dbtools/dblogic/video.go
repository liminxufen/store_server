package dblogic

import (
	"fmt"
	"github.com/store_server/dbtools/driver"
	m "github.com/store_server/dbtools/models"
	"github.com/store_server/logger"
	"sync"
	"time"
)

//JOOX CMS VIDEO相关操作
type VideosDriver struct {
	*driver.CMSDriver
	*BaseDriver
	lock sync.RWMutex
}

func NewVideosDriver(cmsDriver *driver.CMSDriver) *VideosDriver {
	baseDriver := &BaseDriver{cmsDriver, cmsDriver.MusicDB}
	return &VideosDriver{cmsDriver, baseDriver, sync.RWMutex{}}
}

var (
	VoDriver *VideosDriver
)

/* ---------------------------- t_video ------------------------ */
func (vod *VideosDriver) ExecRawQuerySql4Video(sql string, page,
	pagesize int64) ([]*m.Video, int64, error) { //原生query语句
	vos := make([]*m.Video, 0)
	res, total, err := vod.ExecRawQuerySql(sql, page, pagesize, &m.Video{})
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.Video); ok {
			vos = append(vos, t)
		}
	}
	return vos, total, nil
}

func (vod *VideosDriver) InsertOneVideo(vo *m.Video) (int64, error) {
	now := m.TimeNormal{time.Now()}
	vo.FcreateTime = now
	vo.FmodifyTime = now
	_, err := vod.InsertWithModel(vo)
	return vo.Fid, err
}

func (vod *VideosDriver) UpdateOneVideo(id int64, vo *m.Video) (int64, error) {
	vo.FmodifyTime = m.TimeNormal{time.Now()}
	affected, err := vod.UpdateWithModel(&m.Video{}, "Fid = ?", vo, id)
	if err != nil {
		return vo.Fid, err
	}
	if affected != 1 {
		return vo.Fid, fmt.Errorf("update video[%d] error, affect %d raw", vo.Fid, affected)
	}
	return vo.Fid, nil
}

func (vod *VideosDriver) GetOneVideo(id int64) (video *m.Video, err error) {
	video = &m.Video{}
	err = vod.MusicDB.First(video, "Fid=?", id).Error
	return
}

func (vod *VideosDriver) GetVideosByIds(ids []int64) ([]*m.Video, int64, error) {
	var videos []*m.Video
	err := vod.MusicDB.Where("Fid in (?)", ids).Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}
	total := int64(len(videos))
	return videos, total, nil
}

func (vod *VideosDriver) GetVideosByCondition(conds map[string]interface{}, page,
	pagesize int64) ([]*m.Video, int64, error) {
	vos := make([]*m.Video, 0)
	res, total, err := vod.QueryWithModel(&m.Video{}, conds, page, pagesize)
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.Video); ok {
			vos = append(vos, t)
		}
	}
	return vos, total, err
}

func (vod *VideosDriver) UpdateVideoAttr(ids []int64, conds,
	updatesAttrs map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = vod.UpdateWithModel(&m.Video{}, "Fid in (?)", updatesAttrs, ids)
	} else if len(conds) != 0 {
		affected, err = vod.UpdateWithModel(&m.Video{}, conds, updatesAttrs)
	}
	if err != nil {
		return affected, err
	}
	if len(ids) > 0 {
		if affected != int64(len(ids)) {
			return affected, fmt.Errorf("update video count[%d] error, affect %d raw", len(ids), affected)
		}
	}
	return affected, nil
}

func (vod *VideosDriver) DeleteVideos(ids []int64, conds map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = vod.DeleteWithModel(&m.Video{}, "Fid in (?)", ids)
	} else if len(conds) != 0 { //删除条件需严格把关
		logger.Entry().Debugf("delete video condition: %v", conds)
		affected, err = vod.DeleteWithModel(&m.Video{}, conds)
	}
	if err != nil {
		return affected, err
	}
	if len(ids) > 0 {
		if affected != int64(len(ids)) {
			return affected, fmt.Errorf("delete video count[%d] error, affect %d raw", len(ids), affected)
		}
	}
	return affected, nil
}

/* ---------------------------- t_video_extra_os ------------------------ */
func (vod *VideosDriver) ExecRawQuerySql4VideoExtraOs(sql string, page,
	pagesize int64) ([]*m.VideoExtraOs, int64, error) { //原生query语句
	vos := make([]*m.VideoExtraOs, 0)
	res, total, err := vod.ExecRawQuerySql(sql, page, pagesize, &m.VideoExtraOs{})
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.VideoExtraOs); ok {
			vos = append(vos, t)
		}
	}
	return vos, total, nil
}

func (vod *VideosDriver) InsertOneVideoExtraOs(vo *m.VideoExtraOs) (int64, error) {
	now := m.TimeNormal{time.Now()}
	vo.FcreateTime = now
	vo.FmodifyTime = now
	_, err := vod.InsertWithModel(vo)
	return vo.FlocalId, err
}

func (vod *VideosDriver) GetVideoExtraOs(id, region int64) ([]*m.VideoExtraOs, int64, error) {
	var videos []*m.VideoExtraOs
	err := vod.MusicDB.Where("Flocal_id=? and Fregion_id=?", id, region).Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}
	total := int64(len(videos))
	return videos, total, nil
}

func (vod *VideosDriver) GetVideoExtraOsByCondition(conds map[string]interface{}, page,
	pagesize int64) ([]*m.VideoExtraOs, int64, error) {
	videos := make([]*m.VideoExtraOs, 0)
	res, total, err := vod.QueryWithModel(&m.VideoExtraOs{}, conds, page, pagesize)
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.VideoExtraOs); ok {
			videos = append(videos, t)
		}
	}
	return videos, total, err
}

func (vod *VideosDriver) DeleteVideoExtraOs(ids []int64, conds map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = vod.DeleteWithModel(&m.VideoExtraOs{}, "Flocal id in (?)", ids)
	} else if len(conds) != 0 { //删除条件需严格把关
		logger.Entry().Debugf("delete video extra os condition: %v", conds)
		affected, err = vod.DeleteWithModel(&m.VideoExtraOs{}, conds)
	}
	if err != nil {
		return affected, err
	}
	if len(ids) > 0 {
		if affected != int64(len(ids)) {
			return affected, fmt.Errorf("delete video extra os count[%d] error, affect %d raw", len(ids), affected)
		}
	}
	return affected, nil
}

/* ---------------------------- t_video_singer_track ------------------------ */
func (vod *VideosDriver) ExecRawQuerySql4VideoSingerTrack(sql string, page,
	pagesize int64) ([]*m.VideoSingerTrack, int64, error) { //原生query语句
	vos := make([]*m.VideoSingerTrack, 0)
	res, total, err := vod.ExecRawQuerySql(sql, page, pagesize, &m.VideoSingerTrack{})
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.VideoSingerTrack); ok {
			vos = append(vos, t)
		}
	}
	return vos, total, nil
}

func (vod *VideosDriver) InsertOneVideoSingerTrack(vo *m.VideoSingerTrack) (int64, error) {
	now := m.TimeNormal{time.Now()}
	vo.FcreateTime = now
	vo.FmodifyTime = now
	_, err := vod.InsertWithModel(vo)
	return vo.Fid, err
}

func (vod *VideosDriver) UpdateOneVideoSingerTrack(id int64, vo *m.VideoSingerTrack) (int64, error) {
	vo.FmodifyTime = m.TimeNormal{time.Now()}
	affected, err := vod.UpdateWithModel(&m.VideoSingerTrack{}, "Fid = ?", vo, id)
	if err != nil {
		return vo.Fid, err
	}
	if affected != 1 {
		return vo.Fid, fmt.Errorf("update video singer track[%d] error, affect %d raw", vo.Fid, affected)
	}
	return vo.Fid, nil
}

func (vod *VideosDriver) GetVideoSingerTrack(id, region int64) (videos []*m.VideoSingerTrack, total int64, err error) {
	videos = []*m.VideoSingerTrack{}
	err = vod.MusicDB.Where("Flocal_v_id=? and Fregion_id=?", id, region).Find(&videos).Error
	if err != nil {
		return
	}
	total = int64(len(videos))
	return
}

func (vod *VideosDriver) GetVideoSingerTrackByCondition(conds map[string]interface{}, page,
	pagesize int64) ([]*m.VideoSingerTrack, int64, error) {
	vos := make([]*m.VideoSingerTrack, 0)
	res, total, err := vod.QueryWithModel(&m.VideoSingerTrack{}, conds, page, pagesize)
	if err != nil {
		return nil, 0, err
	}
	for _, r := range res {
		if t, ok := r.(*m.VideoSingerTrack); ok {
			vos = append(vos, t)
		}
	}
	return vos, total, err
}

func (vod *VideosDriver) UpdateVideoSingerTrackAttr(ids []int64, conds,
	updatesAttrs map[string]interface{}) (affected int64, err error) {
	if len(ids) > 0 {
		affected, err = vod.UpdateWithModel(&m.VideoSingerTrack{}, "Fid in (?)", updatesAttrs, ids)
	} else if len(conds) != 0 {
		affected, err = vod.UpdateWithModel(&m.VideoSingerTrack{}, conds, updatesAttrs)
	}
	if err != nil {
		return affected, err
	}
	if len(ids) > 0 {
		if affected != int64(len(ids)) {
			return affected, fmt.Errorf("update video singer track count[%d] error, affect %d raw", len(ids), affected)
		}
	}
	return affected, nil
}

func (vod *VideosDriver) ExportAllVideoSingerTracks(sql string) ([][]interface{}, error) { //导出所有视频关联歌曲艺人信息
	return vod.ExportAllRecords(sql)
}

/* ---------------------------- video 相关join查询------------------------ */

func (vod *VideosDriver) JoinQueryWithRawSql(sql string, page, pagesize int64) ([][]interface{}, error) {
	return vod.JoinQueryWithSql(sql, page, pagesize)
}
