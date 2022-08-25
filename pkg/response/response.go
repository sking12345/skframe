// Package response 响应处理工具
package response

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"skframe/pkg/logger"
	"skframe/pkg/ws"
)

type Response struct {
	Code   int         `json:"code,omitempty"`
	Msg    interface{} `json:"msg,omitempty"`
	Result interface{} `json:"result,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

const (
	SuccessStatus = 0
	FailStatus    = -1
)

// JSON 响应 200 和 JSON 数据
func Success(ctx interface{}, data interface{}) bool {
	response := Response{Code: SuccessStatus, Msg: "success", Result: data}
	switch reflect.TypeOf(ctx).String() {
	case "*ws.Context":
		if err := ctx.(*ws.Context).Conn.WriteJSON(response); err != nil {
			logger.Warn("ws", zap.Error(err))
			return false
		}
	case "*gin.Context":
		ctx.(*gin.Context).JSON(http.StatusOK, response)
	}
	return true
}

func Fail(ctx interface{}, msg interface{}, err ...interface{}) bool {
	response := Response{Code: FailStatus, Msg: msg, Errors: err}
	switch reflect.TypeOf(ctx).String() {
	case "*ws.Context":
		if err := ctx.(*ws.Context).Conn.WriteJSON(response); err != nil {
			logger.Warn("ws", zap.Error(err))
			return false
		}
	case "*gin.Context":
		ctx.(*gin.Context).JSON(http.StatusOK, response)
	}
	return true
}

func ValidationError(ctx interface{}, errors interface{}, msg string) bool {
	switch reflect.TypeOf(ctx).String() {
	case "*ws.Context":
		if err := ctx.(*ws.Context).Conn.WriteJSON(Response{Code: FailStatus, Msg: msg, Result: errors}); err != nil {
			logger.Warn("ws", zap.Error(err))
			return false
		}
	case "*gin.Context":
		ctx.(*gin.Context).JSON(http.StatusOK, Response{Code: FailStatus, Msg: msg, Result: errors})
	}
	return true

}
