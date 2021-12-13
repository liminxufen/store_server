package models

import (
	"time"
)

/*---------------------------- mongo文档模型定义 ---------------------------*/

/*video resource文档*/
type PublishedAlbum struct {
	Id              int64         `json:"id" bson:"_id"`
	AlbumName       string        `json:"album_name" bson:"album_name"`
	RegionId        int           `json:"region_id" bson:"region_id"`
	AlbumId         int64         `json:"album_id" bson:"album_id"`
	Deleted         int           `json:"deleted" bson:"deleted"`
	Upc             string        `json:"upc" bson:"upc"`
	SingerName      string        `json:"singer_name" bson:"singer_name"`
	SingerId        int64         `json:"singer_id" bson:"singer_id"`
	CompanyName     string        `json:"company_name" bson:"company_name"`
	LocalName       string        `json:"local_name" bson:"local_name"`
	From            int           `json:"from" bson:"from"`
	LocalPublicTime time.Duration `json:"local_public_time" bson:"local_public_time"`
	LocalStartDate  time.Duration `json:"local_start_date" bson:"local_start_date"`
	LocalCopyright  int           `json:"local_copyright" bson:"local_copyright"`
	LocalStatus     int           `json:"local_status" bson:"local_status"`
	CreateTime      time.Time     `json:"create_time" bson:"create_time"`
	ModifyTime      time.Time     `json:"modify_time" bson:"modify_time"`
}

/*video resource文档*/
type ExternalResource struct {
	Id             int64             `json:"id" bson:"_id"`
	Available      string            `json:"available" bson:"available"`
	InternalFileId string            `json:"internal_file_id" bson:"internal_file_id"` //media unique id, inner server use
	ExternalId     string            `json:"external_id" bson:"external_id"`           //biz_uuid, vod cloud use, may be replace by vod unique fileId
	BizName        string            `json:"biz_name" bson:"biz_name"`                 //video or audio
	Deleted        int               `json:"deleted" bson:"deleted"`
	FileUrl        string            `json:"file_url" bson:"file_url"` //static url
	CreateTime     time.Time         `json:"create_time" bson:"create_time"`
	ModifyTime     time.Time         `json:"modify_time" bson:"modify_time"`
	OriginMediaUrl string            `json:"origin_media_url" bson:"origin_media_url"` //源视频 url
	CoverUrl       string            `json:"cover_url" bson:"cover_url"`
	DownloadUrl    map[string]string `json:"download_url" bson:"download_url"`
}

/*音频信息文件*/
type PreviewAudio struct {
	Id          int64                    `json:"id,omitempty" bson:"_id"`
	Uuid        string                   `json:"uuid,omitempty" bson:"uuid"` //biz_uuid, vod cloud use, may be replace by vod unique fileId
	InnerFileId string                   `json:"inner_file_id" bson:"inner_file_id"`
	Rate        string                   `json:"rate,omitempty" bson:"rate"`
	DownloadUrl []map[string]interface{} `json:"download_url,omitempty" bson:"download_url"`
	TrackId     int                      `json:"track_id,omitempty" bson:"track_id"`
	EncodeStat  int                      `json:"encode_stat,omitempty" bson:"encode_stat"`
	SrcPath     string                   `json:"src_path,omitempty" bson:"src_path"`
	DstPath     string                   `json:"dst_path,omitempty" bson:"dst_path"`
	CreateTime  time.Time                `json:"create_time,omitempty" bson:"create_time"`
	ModifyTime  time.Time                `json:"modify_time,omitempty" bson:"modify_time"`
	Deleted     int                      `json:"deleted,omitempty" bson:"deleted"`
}
