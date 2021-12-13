package models

import (
	"encoding/json"
	"github.com/store_server/utils/errors"
)

//t_video model
type Video struct {
	Fid               int64      `gorm:"column:Fid" json:"Fid" form:"Fid"`
	FregionId         int64      `gorm:"column:Fregion_id" json:"Fregion_id" form:"Fregion_id"`
	Ftitle            string     `gorm:"column:Ftitle" json:"Ftitle" form:"Ftitle"`
	Fstatus           int64      `gorm:"column:Fstatus" json:"Fstatus" form:"Fstatus"`
	FlocalFrom        int64      `gorm:"column:Flocal_from" json:"Flocal_from" form:"Flocal_from"`
	Fsource           int64      `gorm:"column:Fsource" json:"Fsource" form:"Fsource"`
	Fimage            string     `gorm:"column:Fimage" json:"Fimage" form:"Fimage"`
	Fvideo            string     `gorm:"column:Fvideo" json:"Fvideo" form:"Fvideo"`
	Fuuid             string     `gorm:"column:Fuuid" json:"Fuuid" form:"Fuuid"`
	Fmd5              string     `gorm:"column:Fmd5" json:"Fmd5" form:"Fmd5"`
	Fduration         string     `gorm:"column:Fduration" json:"Fduration" form:"Fduration"`
	Fformat           string     `gorm:"column:Fformat" json:"Fformat" form:"Fformat"`
	Fsize             string     `gorm:"column:Fsize" json:"Fsize" form:"Fsize"`
	Fupc              string     `gorm:"column:Fupc" json:"Fupc" form:"Fupc"`
	Fisrc             string     `gorm:"column:Fisrc" json:"Fisrc" form:"Fisrc"`
	Fgrid             string     `gorm:"column:Fgrid" json:"Fgrid" form:"Fgrid"`
	FuploadStatus     int64      `gorm:"column:Fupload_status" json:"Fupload_status" form:"Fupload_status"`
	FcreateTime       TimeNormal `gorm:"column:Fcreate_time" json:"Fcreate_time" form:"Fcreate_time"`
	FmodifyTime       TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
	Fwatermark        int64      `gorm:"column:Fwatermark" json:"Fwatermark" form:"Fwatermark"`
	Fcreator          string     `gorm:"column:Fcreator" json:"Fcreator" form:"Fcreator"`
	FvideoFile        string     `gorm:"column:Fvideo_file" json:"Fvideo_file" form:"Fvideo_file"`
	FlabelModifyTime  TimeNormal `gorm:"column:Flabel_modify_time" json:"Flabel_modify_time" form:"Flabel_modify_time"`
	FimageFile        string     `gorm:"column:Fimage_file" json:"Fimage_file" form:"Fimage_file"`
	FcopyrightSetting string     `gorm:"column:Fcopyright_setting" json:"Fcopyright_setting" form:"Fcopyright_setting"`
	FlanguageId       int64      `gorm:"column:Flanguage_id" json:"Flanguage_id" form:"Flanguage_id"`
	FvideoType        string     `gorm:"column:Fvideo_type" json:"Fvideo_type" form:"Fvideo_type"`
}

func (Video) TableName() string {
	return "t_video"
}

func (video *Video) Encoder() ([]byte, error) {
	if video == nil {
		return nil, errors.New("invalid video pointer")
	}
	s, err := json.Marshal(*video)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (video *Video) Decoder(value []byte) error {
	if video == nil {
		return errors.New("invalid video pointer")
	}
	if err := json.Unmarshal(value, video); err != nil {
		return err
	}
	return nil
}

//t_video_aid model
type VideoAid struct {
	Fid         int64      `gorm:"column:Fid" json:"Fid" form:"Fid"`
	FlocalVId   int64      `gorm:"column:Flocal_v_id" json:"Flocal_v_id" form:"Flocal_v_id"`
	FregionId   int64      `gorm:"column:Fregion_id" json:"Fregion_id" form:"Fregion_id"`
	Ftype       int64      `gorm:"column:Ftype" json:"Ftype" form:"Ftype"`
	FitemId     int64      `gorm:"column:Fitem_id" json:"Fitem_id" form:"Fitem_id"`
	FcreateTime TimeNormal `gorm:"column:Fcreate_time" json:"Fcreate_time" form:"Fcreate_time"`
	FmodifyTime TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
}

func (VideoAid) TableName() string {
	return "t_video_aid"
}

func (videoAd *VideoAid) Encoder() ([]byte, error) {
	if videoAd == nil {
		return nil, errors.New("invalid video aid pointer")
	}
	s, err := json.Marshal(*videoAd)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (videoAd *VideoAid) Decoder(value []byte) error {
	if videoAd == nil {
		return errors.New("invalid video aid pointer")
	}
	if err := json.Unmarshal(value, videoAd); err != nil {
		return err
	}
	return nil
}

//t_video_upload model
type VideoUpload struct {
	Fid         int64      `gorm:"column:Fid" json:"Fid" form:"Fid"`
	FregionId   int64      `gorm:"column:Fregion_id" json:"Fregion_id" form:"Fregion_id"`
	Fvideo      string     `gorm:"column:Fvideo" json:"Fvideo" form:"Fvideo"`
	Fmd5        string     `gorm:"column:Fmd5" json:"Fmd5" form:"Fmd5"`
	FcreateTime TimeNormal `gorm:"column:Fcreate_time" json:"Fcreate_time" form:"Fcreate_time"`
	FmodifyTime TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
}

func (VideoUpload) TableName() string {
	return "t_video_upload"
}

func (videoUd *VideoUpload) Encoder() ([]byte, error) {
	if videoUd == nil {
		return nil, errors.New("invalid video upload pointer")
	}
	s, err := json.Marshal(*videoUd)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (videoUd *VideoUpload) Decoder(value []byte) error {
	if videoUd == nil {
		return errors.New("invalid video upload pointer")
	}
	if err := json.Unmarshal(value, videoUd); err != nil {
		return err
	}
	return nil
}

//t_video_import model
type VideoImport struct {
	Fid         int64      `gorm:"column:Fid" json:"Fid" form:"Fid"`
	FvId        int64      `gorm:"column:Fv_id" json:"Fv_id" form:"Fv_id"`
	FlocalPath  string     `gorm:"column:Flocal_path" json:"Flocal_path" form:"Flocal_path"`
	FsrcPath    string     `gorm:"column:Fsrc_path" json:"Fsrc_path" form:"Fsrc_path"`
	Fstatus     int64      `gorm:"column:Fstatus" json:"Fstatus" form:"Fstatus"`
	FcreateTime TimeNormal `gorm:"column:Fcreate_time" json:"Fcreate_time" form:"Fcreate_time"`
	FmodifyTime TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
}

func (VideoImport) TableName() string {
	return "t_video_import"
}

func (videoIt *VideoImport) Encoder() ([]byte, error) {
	if videoIt == nil {
		return nil, errors.New("invalid video import pointer")
	}
	s, err := json.Marshal(*videoIt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (videoIt *VideoImport) Decoder(value []byte) error {
	if videoIt == nil {
		return errors.New("invalid video import pointer")
	}
	if err := json.Unmarshal(value, videoIt); err != nil {
		return err
	}
	return nil
}

//t_video_extra_os model
type VideoExtraOs struct {
	FlocalId        int64      `gorm:"column:Flocal_id" json:"Flocal_id" form:"Flocal_id"`
	FvId            int64      `gorm:"column:Fv_id" json:"Fv_id" form:"Fv_id"`
	FregionId       int64      `gorm:"column:Fregion_id" json:"Fregion_id" form:"Fregion_id"`
	Ftype           int64      `gorm:"column:Ftype" json:"Ftype" form:"Ftype"`
	Fdesc           string     `gorm:"column:Fdesc" json:"Fdesc" form:"Fdesc"`
	FlocalTitle     string     `gorm:"column:Flocal_title" json:"Flocal_title" form:"Flocal_title"`
	Fstatus         int64      `gorm:"column:Fstatus" json:"Fstatus" form:"Fstatus"`
	FlocalImage     string     `gorm:"column:Flocal_image" json:"Flocal_image" form:"Flocal_image"`
	FtrackList      string     `gorm:"column:Ftrack_list" json:"Ftrack_list" form:"Ftrack_list"`
	FsingerList     string     `gorm:"column:Fsinger_list" json:"Fsinger_list" form:"Fsinger_list"`
	FtagList        string     `gorm:"column:Ftag_list" json:"Ftag_list" form:"Ftag_list"`
	Fvip            int64      `gorm:"column:Fvip" json:"Fvip" form:"Fvip"`
	Fsubscript      int64      `gorm:"column:Fsubscript" json:"Fsubscript" form:"Fsubscript"`
	FcreateTime     TimeNormal `gorm:"column:Fcreate_time" json:"Fcreate_time" form:"Fcreate_time"`
	FpubTime        TimeNormal `gorm:"column:Fpub_time" json:"Fpub_time" form:"Fpub_time"`
	FmodifyTime     TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
	FgifPic         string     `gorm:"column:Fgif_pic" json:"Fgif_pic" form:"Fgif_pic"`
	FviewCount      int64      `gorm:"column:Fview_count" json:"Fview_count" form:"Fview_count"`
	FmatchStatus    string     `gorm:"column:Fmatch_status" json:"Fmatch_status" form:"Fmatch_status"`
	FlocalCopyright int64      `gorm:"column:Flocal_copyright" json:"Flocal_copyright" form:"Flocal_copyright"`
	FlanguageId     int64      `gorm:"column:Flanguage_id" json:"Flanguage_id" form:"Flanguage_id"`
}

func (VideoExtraOs) TableName() string {
	return "t_video_extra_os"
}

func (videoEs *VideoExtraOs) Encoder() ([]byte, error) {
	if videoEs == nil {
		return nil, errors.New("invalid video extra os pointer")
	}
	s, err := json.Marshal(*videoEs)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (videoEs *VideoExtraOs) Decoder(value []byte) error {
	if videoEs == nil {
		return errors.New("invalid video extra os pointer")
	}
	if err := json.Unmarshal(value, videoEs); err != nil {
		return err
	}
	return nil
}

//t_video_singer_track model
type VideoSingerTrack struct {
	Fid         int64      `gorm:"column:Fid" json:"Fid" form:"Fid"`
	FlocalVId   int64      `gorm:"column:Flocal_v_id" json:"Flocal_v_id" form:"Flocal_v_id"`
	FregionId   int64      `gorm:"column:Fregion_id" json:"Fregion_id" form:"Fregion_id"`
	Ftype       int64      `gorm:"column:Ftype" json:"Ftype" form:"Ftype"`
	FsingerId   int64      `gorm:"column:Fsinger_id" json:"Fsinger_id" form:"Fsinger_id"`
	FtrackId    int64      `gorm:"column:Ftrack_id" json:"Ftrack_id" form:"Ftrack_id"`
	Fstatus     int64      `gorm:"column:Fstatus" json:"Fstatus" form:"Fstatus"`
	FcreateTime TimeNormal `gorm:"column:Fcreate_time" json:"Fcreate_time" form:"Fcreate_time"`
	FmodifyTime TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
}

func (VideoSingerTrack) TableName() string {
	return "t_video_singer_track"
}

func (videoSgTk *VideoSingerTrack) Encoder() ([]byte, error) {
	if videoSgTk == nil {
		return nil, errors.New("invalid video singer track pointer")
	}
	s, err := json.Marshal(*videoSgTk)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (videoSgTk *VideoSingerTrack) Decoder(value []byte) error {
	if videoSgTk == nil {
		return errors.New("invalid video singer track pointer")
	}
	if err := json.Unmarshal(value, videoSgTk); err != nil {
		return err
	}
	return nil
}
