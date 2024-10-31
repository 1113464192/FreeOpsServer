package controller

import (
	"FreeOps/internal/consts"
	"FreeOps/internal/service"
	"FreeOps/pkg/api"
	"FreeOps/pkg/logger"
	"FreeOps/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

// UpdateGame
// @Tags 游戏服相关
// @title 新增/修改游戏服信息
// @description 新增不用传游戏服ID，修改才传游戏服ID
// @Summary 新增/修改游戏服信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateGameReq true "传新增或者修改游戏服的所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /games [post]
func UpdateGame(c *gin.Context) {
	var (
		gameReq api.UpdateGameReq
		err     error
	)
	if err = c.ShouldBind(&gameReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.GameServiceApp().UpdateGame(&gameReq); err != nil {
		logger.Log().Error("game", "创建/修改游戏服失败", err)
		c.JSON(500, util.ServerErrorResponse("创建/修改游戏服失败", err))
		return
	}
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// GetGames
// @Tags 游戏服相关
// @title 查询游戏服信息
// @description 查询所有/指定条件游戏服的信息
// @Summary 查询游戏服信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetGamesReq true "传对应条件的参数"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /games [get]
func GetGames(c *gin.Context) {
	var params api.GetGamesReq
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	res, err := service.GameServiceApp().GetGames(&params)
	if err != nil {
		logger.Log().Error("game", "查询游戏服信息失败", err)
		c.JSON(500, util.ServerErrorResponse("查询游戏服信息失败", err))
		return
	}

	logger.Log().Info("game", "查询游戏服信息成功")
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
		Data: res,
	})
}

// DeleteGames
// @Tags 游戏服相关
// @title 删除游戏服
// @description 删除指定游戏服
// @Summary 删除游戏服
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param ids body api.IdsReq true "游戏服ID"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /games [delete]
func DeleteGames(c *gin.Context) {
	var param api.IdsReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err := service.GameServiceApp().DeleteGames(param.Ids); err != nil {
		logger.Log().Error("game", "删除游戏服失败", err)
		c.JSON(500, util.ServerErrorResponse("删除游戏服失败", err))
		return
	}

	logger.Log().Info("game", "删除游戏服成功", fmt.Sprintf("ID: %v", param.Ids))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}

// UpdateGameStatus
// @Tags 游戏服相关
// @title 修改游戏服状态
// @description 修改游戏服状态
// @Summary 修改游戏服状态
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.UpdateGameStatusReq true "传修改游戏服状态所需参数"
// @Success 200 {object} api.Response "{"code": "0000", msg: "string", data: "string"}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"code": "", msg: "", data: ""}"
// @Router /games/status [patch]
func UpdateGameStatus(c *gin.Context) {
	var (
		statusReq api.UpdateGameStatusReq
		err       error
	)
	if err = c.ShouldBind(&statusReq); err != nil {
		c.JSON(500, util.BindErrorResponse(err))
		return
	}
	if err = service.GameServiceApp().UpdateGameStatus(&statusReq); err != nil {
		logger.Log().Error("game", "修改游戏服状态失败", err)
		c.JSON(500, util.ServerErrorResponse("修改游戏服状态失败", err))
		return
	}

	logger.Log().Info("game", "修改游戏服状态成功", fmt.Sprintf("ID: %d", statusReq.Id))
	c.JSON(200, api.Response{
		Code: consts.SERVICE_SUCCESS_CODE,
		Msg:  "Success",
	})
}
