package controlles

import (
	"skframe/app/models"
	"skframe/app/models/user"
	"skframe/app/models/user_friend"
	"skframe/app/requests"
	"skframe/pkg/password"
	"skframe/pkg/response"
	"skframe/pkg/try"
)

type UserController struct {
	baseController
}

//func (*TestController) CreatToken(ctx *ControllerContext) {
//	//token := jwt.NewJWT().IssueToken(123, "xxxx")
//}

func (*UserController) Register(ctx interface{}) {
	request := requests.UserRegister{}
	if ok := requests.Validate(ctx, &request, requests.UserRegisterVer, "数据验证错误"); ok == false {
		return
	}
	try.NewTry(func() {
		mdl := user.User{
			Name:     request.Name,
			Account:  request.Account,
			Password: password.HashPassword(request.Password),
		}
		mdl.Create(nil)
		response.Success(ctx, nil)
	}).Catch(func(err interface{}) {

	}).Run()
}

func (*UserController) Login(ctx interface{}) {
	request := requests.UserLogin{}
	if ok := requests.Validate(ctx, &request, requests.UserLoginVer, "数据验证错误"); ok == false {
		return
	}
	try.NewTry(func() {
		mdl := user.First(models.SqlOpt{
			Field: "*",
			Where: map[string]interface{}{"account": request.Account},
		}, nil, false)
		if mdl.ID <= 0 {
			panic("账号不存在")
		}
		if password.ComparePasswords(mdl.Password, request.Password) == false {
			panic("密码错误")
		}
		friendList := user_friend.Find(models.SqlOpt{
			Field: "friend_id,alias_name",
			Where: map[string]interface{}{"user_id": mdl.ID},
		}, nil, false)

		data := map[string]interface{}{
			"name":    mdl.Name,
			"id":      mdl.ID,
			"friends": friendList,
		}
		response.Success(ctx, data)
	}).Catch(func(err interface{}) {
		response.Fail(ctx, err)
	}).Run()
}

func (*UserController) AddFriend(ctx interface{}) {

}
