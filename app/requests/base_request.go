package requests

import (
	"github.com/thedevsaddam/govalidator"
)

type BaseRequest struct {
	Page int `json:"page" form:"page" uri:"page" valid:"page"`
	Size int `json:"size" form:"size" uri:"size" valid:"size"`
}

func BaseRequestVer(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"page": []string{"required"},
		"size": []string{"required"},
	}

	messages := govalidator.MapData{}
	errs := validate(data, rules, messages)
	return errs
}
