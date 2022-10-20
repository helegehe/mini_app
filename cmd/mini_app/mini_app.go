package main

import (
	"flag"
	"fmt"
	"github.com/helegehe/mini_app/internal/pkg"
	"github.com/helegehe/mini_app/internal/router"
	"github.com/helegehe/mini_app/tools/logger"
	"github.com/helegehe/mini_app/tools/logger/zaplog"
)

var (
	configFilePath = "../../configs/mini_app.yaml"
	showVersion    = false
)

func init() {
	flag.StringVar(&configFilePath, "c", configFilePath, "config file path")
	flag.BoolVar(&showVersion, "v", showVersion, "show app version")
}
func main() {
	flag.Parse()
	if showVersion{
		fmt.Println("version 1.0")
		return
	}
	pkg.InitConfig(configFilePath)
	conf := pkg.GlobalConfig

	// 开始初始化db、log、路由等
	zapLogger := zaplog.InitLogger("mini_app", conf.Logger)
	logger.SetLogger(zapLogger)
	pkg.InitDB(conf.MongoConfig.MongoRepo,conf.MySqlConfig.MysqlRepo,conf.Logger.Level)

	router.InitRouter(conf.AppConfig.Port)
}
