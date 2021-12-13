package models

import (
	"encoding/json"
	"github.com/store_server/utils/errors"
)

//t_track model
type Track struct {
	FtrackId           int64      `gorm:"column:Ftrack_id;int(11);not null;primary_key" json:"Ftrack_id" form:"Ftrack_id"`
	FtrackName         string     `gorm:"column:Ftrack_name;varchar(255)" json:"Ftrack_name" form:"Ftrack_name"`
	FalbumId           int64      `gorm:"column:Falbum_id;int(11)" json:"Falbum_id" form:"Falbum_id"`
	Ftype              int64      `gorm:"column:Ftype;int(11)" json:"Ftype" form:"Ftype"`
	Flanguage          int64      `gorm:"column:Flanguage;int(11)" json:"Flanguage" form:"Flanguage"`
	Fsinger            int64      `gorm:"column:Fsinger;int(11)" json:"Fsinger" form:"Fsinger"`
	Fmovie             string     `gorm:"column:Fmovie;varchar(255)" json:"Fmovie" form:"Fmovie"`
	Fsize              int64      `gorm:"column:Fsize;int(11)" json:"Fsize" form:"Fsize"`
	Fduration          int64      `gorm:"column:Fduration;int(11)" json:"Fduration" form:"Fduration"`
	FsingerId1         int64      `gorm:"column:Fsinger_id1;int(11)" json:"Fsinger_id1" form:"Fsinger_id1"`
	FsingerId2         int64      `gorm:"column:Fsinger_id2;int(11)" json:"Fsinger_id2" form:"Fsinger_id2"`
	FsingerId3         int64      `gorm:"column:Fsinger_id3;int(11)" json:"Fsinger_id3" form:"Fsinger_id3"`
	FsingerId4         int64      `gorm:"column:Fsinger_id4;int(11)" json:"Fsinger_id4" form:"Fsinger_id4"`
	Fprice1            int64      `gorm:"column:Fprice1;int(11)" json:"Fprice1" form:"Fprice1"`
	Fprice2            int64      `gorm:"column:Fprice2;int(11)" json:"Fprice2" form:"Fprice2"`
	Fprice3            int64      `gorm:"column:Fprice3;int(11)" json:"Fprice3" form:"Fprice3"`
	Fisrc              string     `gorm:"column:Fisrc;varchar(255)" json:"Fisrc" form:"Fisrc"`
	Fattribute1        int64      `gorm:"column:Fattribute_1;int(11)" json:"Fattribute_1" form:"Fattribute_1"`
	Fattribute2        int64      `gorm:"column:Fattribute_2;int(11)" json:"Fattribute_2" form:"Fattribute_2"`
	Fattribute3        int64      `gorm:"column:Fattribute_3;int(11)" json:"Fattribute_3" form:"Fattribute_3"`
	Fattribute4        int64      `gorm:"column:Fattribute_4;int(11)" json:"Fattribute_4" form:"Fattribute_4"`
	Fgenre             int64      `gorm:"column:Fgenre;int(11)" json:"Fgenre" form:"Fgenre"`
	FsingerAll         string     `gorm:"column:Fsinger_all;varchar(255)" json:"Fsinger_all" form:"Fsinger_all"`
	Flocation          int64      `gorm:"column:Flocation;int(11)" json:"Flocation" form:"Flocation"`
	FvalidTime         TimeNormal `gorm:"column:Fvalid_time" json:"Fvalid_time" form:"Fvalid_time"`
	FuploadTime        TimeNormal `gorm:"column:Fupload_time" json:"Fupload_time" form:"Fupload_time"`
	FmodifyTime        TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
	FlastestModifyTime TimeNormal `gorm:"column:Flasttest_modify_time" json:"Flastest_modify_time" form:"Flastest_modify_time"`
	FtrackCId          int64      `gorm:"column:Ftrack_c_id;int(11)" json:"Ftrack_c_id" form:"Ftrack_c_id"`
	Flyric             int64      `gorm:"column:Flyric;tinyint(4)" json:"Flyric" form:"Flyric"`
	FportalLyric       int64      `gorm:"column:Fportal_lyric;tinyint(4)" json:"Fportal_lyric" form:"Fportal_lyric"`
	Fstatus            int64      `gorm:"column:Fstatus;int(11)" json:"Fstatus" form:"Fstatus"`
	FgoSoso            int64      `gorm:"column:Fgo_soso;int(11)" json:"Fgo_soso" form:"Fgo_soso"`
	Fnote              int64      `gorm:"column:Fnote;int(11)" json:"Fnote" form:"Fnote"`
	Fversion           int64      `gorm:"column:Fversion;int(11)" json:"Fversion" form:"Fversion"`
	Fattribute5        int64      `gorm:"column:Fattribute_5;int(11)" json:"Fattribute_5" form:"Fattribute_5"`
	Fattribute6        int64      `gorm:"column:Fattribute_6;int(11)" json:"Fattribute_6" form:"Fattribute_6"`
	FtrackMid          string     `gorm:"column:Ftrack_mid;varchar(255)" json:"Ftrack_mid" form:"Ftrack_mid"`
	FlinkMv            int64      `gorm:"column:Flink_mv;int(11)" json:"Flink_mv" form:"Flink_mv"`
	FmediaId           int64      `gorm:"column:Fmedia_id;int(11)" json:"Fmedia_id" form:"Fmedia_id"`
	FlinkRing          int64      `gorm:"column:Flink_ring;int(11)" json:"Flink_ring" form:"Flink_ring"`
	FgenreIds          Int64s     `gorm:"column:Fgenre_ids;type:json" json:"Fgenre_ids" form:"Fgenre_ids"`
}

func (Track) TableName() string {
	return "t_track"
}

func (track *Track) Encoder() ([]byte, error) {
	if track == nil {
		return nil, errors.New("invalid track pointer")
	}
	s, err := json.Marshal(*track)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (track *Track) Decoder(value []byte) error {
	if track == nil {
		return errors.New("invalid track pointer")
	}
	if err := json.Unmarshal(value, track); err != nil {
		return err
	}
	return nil
}

//t_track_extra_os model
type TrackExtraOs struct {
	FtrackId          int64      `gorm:"column:Ftrack_id;int(11);not null;primary_key" json:"Ftrack_id" form:"Ftrack_id"`
	Fregion           int64      `gorm:"column:Fregion;tinyint(4)" json:"Fregion" form:"Fregion"`
	FlocalName        string     `gorm:"column:Flocal_name;varchar(255)" json:"Flocal_name" form:"Flocal_name"`
	FlocalCopyright   int64      `gorm:"column:Flocal_copyright;int(11)" json:"Flocal_copyright" form:"Flocal_copyright"`
	FlocalValidTime   TimeNormal `gorm:"column:Flocal_valid_time" json:"Flocal_valid_time" form:"Flocal_valid_time"`
	FlocalStatus      int64      `gorm:"column:Flocal_status;tinyint(4)" json:"Flocal_status" form:"Flocal_status"`
	FlocalMovie       string     `gorm:"column:Flocal_movie;varchar(255)" json:"Flocal_movie" form:"Flocal_movie"`
	FlocalFrom        int64      `gorm:"column:Flocal_from;tinyint(4)" json:"Flocal_from" form:"Flocal_from"`
	FactionTemplateId int64      `gorm:"column:Faction_template_id;int(11)" json:"Faction_template_id" form:"Faction_template_id"`
	FmvId             int64      `gorm:"column:Fmv_id;int(11)" json:"Fmv_id" form:"Fmv_id"`
	FcopyrightLimit   int64      `gorm:"column:Fcopyright_limit;int(11)" json:"Fcopyright_limit" form:"Fcopyright_limit"`
	FreplaceId        int64      `gorm:"column:Freplace_id;int(11)" json:"Freplace_id" form:"Freplace_id"`
	FlocalIsrc        string     `gorm:"column:Flocal_isrc;varchar(255)" json:"Flocal_isrc" form:"Flocal_isrc"`
	FlocalLabel       string     `gorm:"column:Flocal_label;varchar(255)" json:"Flocal_label" form:"Flocal_label"`
	Fsupplier         string     `gorm:"column:Fsupplier;varchar(255)" json:"Fsupplier" form:"Fsupplier"`
	FallSources       string     `gorm:"column:Fall_sources;varchar(255)" json:"Fall_sources" form:"Fall_sources"`
	FlocalOtherName   string     `gorm:"column:Flocal_other_name;varchar(255)" json:"Flocal_other_name" form:"Flocal_other_name"`
	FmodifyTime       TimeNormal `gorm:"column:Fmodify_time" json:"Fmodify_time" form:"Fmodify_time"`
}

func (TrackExtraOs) TableName() string {
	return "t_track_extra_os"
}

func (track *TrackExtraOs) Encoder() ([]byte, error) {
	if track == nil {
		return nil, errors.New("invalid track extra os pointer")
	}
	s, err := json.Marshal(*track)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (track *TrackExtraOs) Decoder(value []byte) error {
	if track == nil {
		return errors.New("invalid track extra os pointer")
	}
	if err := json.Unmarshal(value, track); err != nil {
		return err
	}
	return nil
}
