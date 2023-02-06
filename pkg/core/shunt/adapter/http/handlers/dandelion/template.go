package dandelion

import (
	"genesis/pkg/core/shunt/adapter/facade/dandelion_facade"
	"genesis/pkg/core/shunt/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	*handlers.Handler
}

func NewTemplateHandler(handler *handlers.Handler) *TemplateHandler {
	return &TemplateHandler{
		Handler: handler,
	}
}

func (h *TemplateHandler) TemplateGenerate(ctx *gin.Context) {
	var rb dandelion_facade.TemplateGenerateDTO

	if err := ctx.ShouldBind(&rb); err != nil {
		h.SendFailure(ctx, handlers.ParameterConvertError, handlers.StatusText(handlers.ParameterConvertError))
		return
	}
	if err := dandelion_facade.TempalteGenerate(rb); err != nil {
		h.SendFailure(ctx, err.Error())
	} else {
		h.SendSuccess(ctx)
	}
}
