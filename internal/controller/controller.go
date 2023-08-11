package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/helegehe/mini_app/tools/httper"
	"math"
	"strconv"
	"time"
)

func Test(ctx *gin.Context) {
	t, _ := ctx.GetQuery("t")
	t2, err := parseTime(t)
	httper.HandleResponse(ctx, err, map[string]interface{}{"t": t2.String()})
}

func parseTime(s string) (time.Time, error) {
	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		ns = math.Round(ns*1000) / 1000
		return time.Unix(int64(s), int64(ns*float64(time.Second))).UTC(), nil
	}
	return time.Now(),nil
}
