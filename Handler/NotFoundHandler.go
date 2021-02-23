package Handler

import (
	"github.com/kataras/iris/v12"
)

// 400错误
func NotFoundHandler(ctx iris.Context) {
	_, _ = ctx.Write([]byte("404"))
}
