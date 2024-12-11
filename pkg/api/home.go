package api

type GetHomeInfoRes struct {
	ProjectCount         int64 `json:"projectCount"`
	HostCount            int64 `json:"hostCount"`
	TaskNeedApproveCount int64 `json:"taskNeedApproveCount"`
	TaskRunningCount     int64 `json:"taskRunningCount"`
}
