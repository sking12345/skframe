package requests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"reflect"
	"skframe/pkg/response"
	"skframe/pkg/ws"
)

type ValidatorFunc func(interface{}) map[string][]string

func RequestData(ctx interface{}, obj interface{}) error {
	ctxName := reflect.TypeOf(ctx).String()
	switch ctxName {
	case "*ws.Context":
		body, _ := json.Marshal(ctx.(*ws.Context).Body)
		if err := json.Unmarshal(body, obj); err != nil {
			return err
		}
	case "*gin.Context":
		if err := ctx.(*gin.Context).ShouldBind(obj); err != nil {
			return err
		}
	}
	return nil
}

func Validate(ctx interface{}, obj interface{}, handler ValidatorFunc, msg string) bool {
	if err := RequestData(ctx, obj); err != nil {
		response.ValidationError(ctx, err, msg)
		return false
	}
	if err := handler(obj); len(err) > 0 {
		response.ValidationError(ctx, err, msg)
		return false
	}
	return true
}

func validate(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	// 配置选项
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}

	// 开始验证
	return govalidator.New(opts).ValidateStruct()
}
