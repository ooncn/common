package Handler

import (
	"github.com/kataras/iris/v12"
)

func (c Context) Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers", ctx.GetHeader("Access-Control-Request-Headers"))
	ctx.Header("Access-Control-Allow-Methods", ctx.GetHeader("Access-Control-Request-Method"))
	ctx.Header("Access-Control-Max-Age", "1800")
	ctx.Header("Allow", "GET, POST, OPTIONS")
	ctx.Header("Content-Type", "application/json;charset=UTF-8")
	if ctx.Method() == "OPTIONS" {
		ctx.WriteString("")
		ctx.StopExecution()
	} else { /*
			fmt.Println(ctx.GetHeader("X-Forwarded-For"))
			fmt.Println(ctx.Path())
			fmt.Println(ctx.Method())*/
		ctx.Next()
	}
}

func (c Context) LocalCors(ctx iris.Context) {
	Origin := ctx.GetHeader("Origin")
	//fmt.Println("LocalCors:"+Origin)
	ctx.Header("Access-Control-Allow-Origin", Origin)
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers", ctx.GetHeader("Access-Control-Request-Headers"))
	ctx.Header("Access-Control-Allow-Methods", ctx.GetHeader("Access-Control-Request-Method"))
	ctx.Header("Access-Control-Max-Age", "1800")
	ctx.Header("Allow", "GET, POST, OPTIONS")
	//ctx.Header("Content-Type", "application/json;charset=UTF-8")
	if ctx.Method() == "OPTIONS" {
		ctx.WriteString("")
		ctx.StopExecution()
	} else {
		ctx.Next()
	}
}
