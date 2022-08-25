package requests

import (
	"github.com/thedevsaddam/govalidator"
)

type UserRegister struct {
	Account         string `json:"account" valid:"account"`
	Password        string `json:"password" valid:"password"`
	PasswordConfirm string `json:"password_confirm" valid:"password_confirm"`
	Name            string `json:"name" valid:"name"`
}

func UserRegisterVer(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"name":             []string{"required"},
		"account":          []string{"required"},
		"password":         []string{"required"},
		"password_confirm": []string{"required"},
	}
	messages := govalidator.MapData{}
	errs := validate(data, rules, messages)
	return errs
}

type UserLogin struct {
	Account  string `json:"account" valid:"account"`
	Password string `json:"password" valid:"password"`
}

func UserLoginVer(data interface{}) map[string][]string {
	rules := govalidator.MapData{
		"account":  []string{"required"},
		"password": []string{"required"},
	}
	messages := govalidator.MapData{}
	errs := validate(data, rules, messages)
	return errs
}
