// @Author: YangPing
// @Create: 2023/10/21
// @Description: 用户请求返回配置

package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type ResultCode int
type ResultMsg string

const (
	DefaultStatus             ResultCode = -1
	StatusOK                             = 200
	StatusFound                          = 302
	StatusBadRequest                     = 400
	StatusNotFound                       = 404
	StatusInternalServerError            = 500

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
	DefaultStatus:             ResultMsg(""),
	StatusOK:                  ResultMsg("OK"),
	StatusBadRequest:          ResultMsg("Bad Request"),
	StatusFound:               ResultMsg("Found"),
	StatusNotFound:            ResultMsg("Not Found"),
	StatusInternalServerError: ResultMsg("Internal Server Error"),
	RequestParameterError:     ResultMsg("Request Parameter Error"),
	DataConvertError:          ResultMsg("Data Convert Error"),
	RequestCheckTokenError:    ResultMsg("Token Is Not Exists, Please Login"),
	ParameterConvertError:     ResultMsg("Parameter Error, Please Check Parameter"),
	InitDataBaseError:         ResultMsg("Init DataBase Error"),
	QueryDBError:              ResultMsg("Query DataBase Error"),
	RequestCheckTokenTimeOut:  ResultMsg("request check token time out"),
	RequestTokenNotFound:      ResultMsg("request token not found, please login first"),
	UserNotFound:              ResultMsg("user not found"),
	CreateTokenError:          ResultMsg("create token error"),
	AddUserError:              ResultMsg("add user error"),
	DataNotFound:              ResultMsg("data not found"),
	RateLimit:                 ResultMsg("rate limit"),
}

func StatusText(code ResultCode) ResultMsg {
	return statusText[code]
}

type Result struct {
	Code ResultCode `json:"code"`
	Msg  ResultMsg  `json:"msg"`
	Data any        `json:"data"`
}

func SendSuccess(ctx *gin.Context, arg ...any) {
	var result Result

	buildResult(true, &result, arg...)

	ctx.JSON(StatusOK, result)
}

func SendFailure(ctx *gin.Context, arg ...any) {
	var (
		result Result
	)

	buildResult(false, &result, arg...)

	ctx.JSON(StatusOK, result)
}

func buildResult(ok bool, result *Result, arg ...any) {
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

func StatusValidator(obj any, rawErr error) ResultMsg {
	return ResultMsg(ValidateErr(obj, rawErr).Error())
}

// ValidateErr
// Name string `json:name binding:"notempty" msg:"name must is not null"`
func ValidateErr(obj any, rawErr error) error {
	validationErrs, ok := rawErr.(validator.ValidationErrors)
	if !ok {
		return rawErr
	}
	var errString []string
	for _, validationErr := range validationErrs {
		field, ok := reflect.TypeOf(obj).FieldByName(validationErr.Field())
		if ok {
			if e := field.Tag.Get("msg"); e != "" {
				//errString = append(errString, fmt.Sprintf("%s: %s", validationErr.Namespace(), e))
				errString = append(errString, e)
				continue
			}
		}
		errString = append(errString, validationErr.Error())
	}
	return errors.New(strings.Join(errString, ", "))
}
