package Handler

import (
	"github.com/kataras/iris/v12"
	"github.com/ooncn/common/constant"
	"github.com/ooncn/common/obj"
	"github.com/ooncn/common/oredis"
	"github.com/ooncn/common/util"
)

func (c Context) AuthMangerBefore(ctx iris.Context) {
	c.SetCxt(ctx)
	shareInformation := "这是处理程序之间可共享的信息"
	requestPath := ctx.Path()
	println("在主处理程序之前： " + requestPath)
	ctx.Values().Set("info", shareInformation)

	if ctx.Method() != "OPTIONS" {
		re := make(map[string]interface{})
		re["code"] = 500
		token := ctx.GetHeader("token")
		if len(token) == 0 {
			token = ctx.URLParam("token")
		}
		if len(token) != 32 {
			c.RespErrorD("LOGOUT")
			return
		}
		user := obj.VoUserToken{}
		err := oredis.Ser.GetType(constant.REDIS_USER_SESSION+token, &user)
		if err != nil {
			c.RespErrorD("LOGOUT")
			return
		}
		//if *login.EndTime < *util.TimeUtil.TimeGet() {
		//	err := sqlSer.Model(&login).Update(db.UserLogin{Token: "0", EndTime: login.EndTime}).Error
		//	if err != nil {
		//		logs.Error(err)
		//	}
		//	ctx.StatusCode(500)
		//	re["msg"] = "LOGOUT"
		//	ctx.JSON(re)
		//	return
		//}

		if util.IsBlank(user.Id) {
			c.RespErrorD("NOT_USER")
			return
		}
		if user.Type != 100 {
			c.RespErrorD("NOT_USER")
			return
		}
	}
	ctx.Next() // execute the next handler, in this case the main one.
}

// 服务中间件
func (c Context) AuthModularBefore(ctx iris.Context) {
	c.SetCxt(ctx)
	shareInformation := "这是处理程序之间可共享的信息"
	requestPath := ctx.Path()
	println("在主处理程序之前： " + requestPath)
	ctx.Values().Set("info", shareInformation)

	// 请求Redis TOKEN 信息
	if ctx.Method() != "OPTIONS" {
		token := ctx.GetHeader("token")
		if len(token) == 0 {
			token = ctx.URLParam("token")
		}
		if len(token) != 32 {
			c.RespErrorD("LOGOUT")
			return
		}
		userStr, err := oredis.Ser.Get("user:session:" + token)
		if err != nil {
			c.RespErrorD("LOGOUT")
			return
		}
		user := obj.VoUserToken{}
		util.JsonToType(userStr, &user)
		//if *login.EndTime < *util.TimeUtil.TimeGet() {
		//	err := sqlSer.Model(&login).Update(db.UserLogin{Token: "0", EndTime: login.EndTime}).Error
		//	if err != nil {
		//		logs.Error(err)
		//	}
		//	ctx.StatusCode(500)
		//	re["msg"] = "LOGOUT"
		//	ctx.JSON(re)
		//	return
		//}

		if err != nil || util.IsBlank(user.Id) {
			c.RespErrorD("NOT_USER")
			return
		}
		/*if *user.Type != 100 {
			ctx.StatusCode(500)
			re["msg"] = "NO_AUTHORITY"
			ctx.JSON(re)
			return
		}*/
	}
	ctx.Next() // execute the next handler, in this case the main one.
}

func (c Context) AuthApiBefore(ctx iris.Context) {
	c.SetCxt(ctx)
	// 请求Redis TOKEN 信息
	if ctx.Method() != "OPTIONS" {
		token := ctx.GetHeader("token")
		if len(token) == 0 {
			token = ctx.URLParam("token")
		}
		if len(token) != 32 {
			c.RespErrorD("LOGOUT")
			return
		}
		user := obj.VoUserToken{}
		err := oredis.Ser.GetType(constant.REDIS_USER_SESSION+token, &user)
		if err != nil {
			c.RespErrorD("LOGOUT")
			return
		}
		if util.IsBlank(user.Id) {
			c.RespErrorD("NOT_USER")
			return
		}
	}
	ctx.Next() // execute the next handler, in this case the main one.
}
