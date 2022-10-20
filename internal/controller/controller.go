package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/helegehe/mini_app/tools/httper"
)

func Test(ctx *gin.Context) {
	// todo cal service
	httper.HandleResponse(ctx, nil, map[string]interface{}{"id": "1"})
}
