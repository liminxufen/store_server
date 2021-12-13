package common

import (
	"github.com/store_server/dbtools/models"
)

/*track rpc 公共参数*/

//create track rpc request
type CreateTrackRpcReq struct {
	CreateDoc *models.Track `json:"track,omitempty"`
}

//create track rpc response
type CreateTrackRpcRsp struct {
	Id int64 `json:"id"`
}

//delete track rpc request
type DeleteTrackRpcReq struct {
	Id int64 `json:"id,omitempty"`
}

//delete track rpc response
type DeleteTrackRpcRsp struct {
	Id      int64 `json:"id,omitempty"`
	Deleted bool  `json:"deleted,omitempty"`
}

//update track rpc request
type UpdateTrackRpcReq struct {
	Id        int64         `json:"id,omitempty"`
	UpdateDoc *models.Track `json:"track,omitempty"`
}

//update track rpc response
type UpdateTrackRpcRsp struct {
	Id      int64 `json:"id,omitempty"`
	Changed bool  `json:"changed,omitempty"`
}

//search track rpc filter info for request
type SearchTrackRpcReq_FilterInfo struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	SingerId int64  `json:"singer_id,omitempty"`
	AlbumId  int64  `json:"album_id,omitempty"`
}

//search track rpc sort info for request
type SearchTrackRpcReq_SortInfo struct {
	Id         int64 `json:"id,omitempty"`
	CreateTime int   `json:"create_time,omitempty"`
	ModifyTime int   `json:"modify_time,omitempty"`
}

//search track rpc request
type SearchTrackRpcReq struct {
	Start  int64                         `json:"start,omitempty"`
	Count  int64                         `json:"count,omitempty"`
	Filter *SearchTrackRpcReq_FilterInfo `json:"filter,omitempty"`
	Sort   []*SearchTrackRpcReq_SortInfo `json:"sort,omitempty"`
}

//search track rpc response
type SearchTrackRpcRsp struct {
	Data *struct {
		Total int64           `json:"total,omitempty"`
		Start int64           `json:"start,omitempty"`
		Count int64           `json:"count,omitempty"`
		List  []*models.Track `json:"list,omitempty"`
	}
}
