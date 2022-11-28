package handlers

import (
	"github.com/gin-gonic/gin"
)

type ResultCode int
type ResultMsg string

const (
	DefaultStatus              ResultCode = -1
	StatusOK                              = 200
	StatgenesisovedPermanently            = 301
	StatusFound                           = 302
	StatusBadRequest                      = 400
	StatusNotFound                        = 404
	StatusInternalServerError             = 500

	RequestParameterError    = 1001
	RequestCheckTokenError   = 1002
	RequestCheckTokenTimeOut = 1003
	RequestTokenNotFound     = 1004
	CreateTokenError         = 1005
	DataConvertError         = 2001
	ParameterConvertError    = 2002

	InitDataBaseError = 3001
	QueryDBError      = 3002
	UserNotFound      = 3003
	AddUserError      = 3004
	DataNotFound      = 3005
	RateLimit         = 4001
)

var statusText = map[ResultCode]ResultMsg{
	DefaultStatus:              ResultMsg(""),
	StatusOK:                   ResultMsg("OK"),
	StatusBadRequest:           ResultMsg("Bad Request"),
	StatgenesisovedPermanently: ResultMsg("Moved Permanently"),
	StatusFound:                ResultMsg("Found"),
	StatusNotFound:             ResultMsg("Not Found"),
	StatusInternalServerError:  ResultMsg("Internal Server Error"),
	RequestParameterError:      ResultMsg("Request Parameter Error"),
	DataConvertError:           ResultMsg("Data Convert Error"),
	RequestCheckTokenError:     ResultMsg("Token Is Not Exists, Please Login"),
	ParameterConvertError:      ResultMsg("Parameter Error, Please Check Parameter"),
	InitDataBaseError:          ResultMsg("Init DataBase Error"),
	QueryDBError:               ResultMsg("Query DataBase Error"),
	RequestCheckTokenTimeOut:   ResultMsg("request check token time out"),
	RequestTokenNotFound:       ResultMsg("request token not found, please login first"),
	UserNotFound:               ResultMsg("user not found"),
	CreateTokenError:           ResultMsg("create token error"),
	AddUserError:               ResultMsg("add user error"),
	DataNotFound:               ResultMsg("data not found"),
	RateLimit:                  ResultMsg("rate limit"),
}

func StatusText(code ResultCode) ResultMsg {
	return statusText[code]
}

type Result struct {
	Code ResultCode `json:"code"`
	Msg  ResultMsg  `json:"msg"`
	Data any        `json:"data"`
}

func (h *Handler) SendSuccess(ctx *gin.Context, arg ...any) {
	var result Result

	h.buildResult(true, &result, arg...)

	ctx.JSON(StatusOK, result)
}

func (h *Handler) SendFailure(ctx *gin.Context, arg ...any) {
	var (
		result Result
	)

	h.buildResult(false, &result, arg...)

	ctx.JSON(StatusOK, result)
}

func (h *Handler) buildResult(ok bool, result *Result, arg ...any) {
	var (
		code, msg bool
	)
	for _, v := range arg {
		switch v.(type) {
		case ResultCode:
			code = true
			result.Code = v.(ResultCode)
		case ResultMsg:
			msg = true
			result.Msg = v.(ResultMsg)
		default:
			result.Data = v
		}
	}

	if code && !msg {
		result.Msg = statusText[result.Code]
	} else {
		if msg && !code {
			if ok {
				result.Code = StatusOK
			} else {
				result.Code = StatusBadRequest
			}
		} else {
			if !msg && !code {
				if ok {
					result.Code = StatusOK
				} else {
					result.Code = StatusBadRequest
				}
				result.Msg = statusText[result.Code]
			}
		}
	}

}
