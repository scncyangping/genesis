package business

import (
	"github.com/gin-gonic/gin"
	"genesis/pkg/core/shunt/adapter/http/handlers"
	"genesis/pkg/core/shunt/application/cqe/cmd"
	"genesis/pkg/core/shunt/application/service"
)

type AuthHandler struct {
	*handlers.Handler
	userSrv service.AuthSrv
}

func NewAuthHandler(handler *handlers.Handler, srv service.AuthSrv) *AuthHandler {
	return &AuthHandler{
		Handler: handler,
		userSrv: srv,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var requestBody cmd.LoginCmd

	if err := ctx.ShouldBind(&requestBody); err != nil {
		h.SendFailure(ctx, handlers.ParameterConvertError, handlers.StatusText(handlers.ParameterConvertError))
		return
	}
	if rp, err := h.userSrv.Login(&requestBody); err != nil {
		h.SendFailure(ctx, err.Error())
	} else {
		h.SendSuccess(ctx, rp)
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var re cmd.RegisterCmd

	if err := ctx.ShouldBind(&re); err != nil {
		h.SendFailure(ctx, handlers.ParameterConvertError, handlers.StatusText(handlers.ParameterConvertError))
		return
	}
	if uId, err := h.userSrv.Register(&re); err != nil {
		h.SendFailure(ctx, err.Error())
	} else {
		h.SendSuccess(ctx, uId)
	}
}
