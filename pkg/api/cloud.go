package api

type CloudProjectReq struct {
	Name          string `form:"name" json:"name" binding:"required"`                   // 项目名
	CloudPlatform string `form:"cloudPlatform" json:"cloudPlatform" binding:"required"` // 云平台
}

type UpdateCloudProjectReq struct {
	ID uint `form:"id" json:"id"` // 云项目ID
	CloudProjectReq
}

type CreateCloudHostReq struct {
	ProjectId     uint   `form:"projectId" json:"projectId" binding:"required"`         // 项目ID
	CloudPlatform string `form:"cloudPlatform" json:"cloudPlatform" binding:"required"` // 云平台
	HostCount     uint64 `form:"hostCount" json:"hostCount" binding:"required"`         // 需购买服务器数量
}
