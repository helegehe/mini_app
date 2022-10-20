package router

import (
	"github.com/gin-gonic/gin"
	"github.com/helegehe/mini_app/internal/controller"
	"net/http"
)

func InitRouter(port string){
	r := gin.New()

	// 开发阶段 使用gin的Recovery将日志格式化输出到控制台
	r.Use(gin.Recovery())
	// 生产阶段 使用zap输出到固定文件
	r.Use(RecoveryWithZap())
	// todo add  login
	//r.Use(authMiddleware())

	// 将api调用记录全部存放到zap日志中去
	//r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, false))

	// 无需鉴权的模块，用这个路由组注册的路由无需鉴权
	noNeedAuthRouterGroup := r.Group("/mini_app/api/v1")
	RegisterNoNeedAuthRouters(noNeedAuthRouterGroup)

	// 第三方API服务
	_3rdGroup := r.Group("/mini_app/api/v1/3rd")
	Register3rdRouters(_3rdGroup)

	// 需要鉴权的路由
	v1 := r.Group("/mini_app/api/v1")
	v1.GET("test",controller.Test)
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}

// RegisterNoNeedAuthRouters 把零散分布在各处的无需鉴权的路由纳入此处...
func RegisterNoNeedAuthRouters(rg *gin.RouterGroup) {
	// 探活
	rg.GET("/ping", func(ctx *gin.Context) {
		_, _ = Text(ctx.Writer, http.StatusOK, []byte("pong"))
	})
	// todo add login
	//rg.POST("/login", controller.Login)
}

// Register3rdRouters 第三方API服务
func Register3rdRouters(rg *gin.RouterGroup) {
	// todo add path
}

// Text 纯文本格式渲染HTTP Response Body
func Text(w http.ResponseWriter, code int, body []byte) (int, error) {
	w.Header().Set("Content-Category", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	return w.Write(body)
}
