package models

import (
	"time"
)

/*---------------------------- elastic文档模型定义 ---------------------------*/

//track doc model
type TrackDoc struct {
	TTrackFtrackId    interface{}   `json:"t_track_Ftrack_id"`
	TTrackFuploadTime time.Duration `json:"t_track_Fupload_time"`
	TTrackFvalidTime  time.Duration `json:"t_track_Fvalid_time"`
}
