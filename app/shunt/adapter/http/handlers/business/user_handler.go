// @Author: YangPing
// @Create: 2023/10/21
// @Description: 用户控制类

package business

import (
	"genesis/app/shunt/adapter/cqe/cmd"
	"genesis/app/shunt/adapter/cqe/query"
	"genesis/app/shunt/adapter/http/handlers"
	"genesis/app/shunt/service/bus"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	srv bus.UserSrv
}

// NEED WIRE

func NewAuthHandler(srv bus.UserSrv) *AuthHandler {
	return &AuthHandler{
		srv: srv,
	}
}

func (h *AuthHandler) AddUser(ctx *gin.Context) {
	var (
		re cmd.UserSaveCmd
	)

	if err := ctx.ShouldBind(&re); err != nil {
		handlers.SendFailure(ctx, handlers.ParameterConvertError, handlers.StatusValidator(re, err))
		return
	}

	if uId, err := h.srv.Add(&re); err != nil {
		handlers.SendFailure(ctx, err.Error())
	} else {
		handlers.SendSuccess(ctx, uId)
	}
}

func (h *AuthHandler) UpdateUser(ctx *gin.Context) {
	var re cmd.UserUpdateCmd

	if err := ctx.ShouldBind(&re); err != nil {
		handlers.SendFailure(ctx, handlers.ParameterConvertError, handlers.StatusValidator(re, err))
		return
	}

	if err := h.srv.Update(&re); err != nil {
		handlers.SendFailure(ctx, err.Error())
	} else {
		handlers.SendSuccess(ctx)
	}
}

func (h *AuthHandler) QueryUser(ctx *gin.Context) {
	var (
		qu query.UserListQuery
	)
	if err := ctx.ShouldBind(&qu); err != nil {
		handlers.SendFailure(ctx, handlers.ParameterConvertError, handlers.StatusValidator(qu, err))
		return
	}
	if res, err := h.srv.Query(&qu); err != nil {
		handlers.SendFailure(ctx, err.Error())
	} else {
		handlers.SendSuccess(ctx, res)
	}
}

func (h *AuthHandler) GetUserById(ctx *gin.Context) {
	pCode := ctx.Param("id")
	if vo, err := h.srv.GetById(pCode); err != nil {
		handlers.SendFailure(ctx, handlers.ResultMsg(err.Error()))
	} else {
		handlers.SendSuccess(ctx, vo)
	}
}

func (h *AuthHandler) DeleteUserById(ctx *gin.Context) {
	if err := h.srv.DeleteById(ctx.Param("id")); err != nil {
		handlers.SendFailure(ctx, handlers.ResultMsg(err.Error()))
	} else {
		handlers.SendSuccess(ctx)
	}
}

func (h *AuthHandler) DeleteUserByIds(ctx *gin.Context) {
	var p struct {
		Ids []string `json:"ids" binding:"required,min=1" msg:"待删除数据不能为空"`
	}

	if err := ctx.ShouldBind(&p); err != nil || len(p.Ids) < 1 {
		handlers.SendFailure(ctx, handlers.ParameterConvertError, handlers.StatusValidator(p, err))
		return
	}

	m := map[string]any{
		"id": p.Ids,
	}
	if err := h.srv.DeleteByMap(m); err != nil {
		handlers.SendFailure(ctx, handlers.ResultMsg(err.Error()))
	} else {
		handlers.SendSuccess(ctx)
	}
}
