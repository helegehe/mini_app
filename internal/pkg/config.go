package pkg

import (
	"github.com/spf13/viper"
	"time"
)

// GlobalConfig 全局配置索引
var GlobalConfig *Conf

// Conf 配置文件映射
type Conf struct {
	MongoConfig *MongoConfig
	MySqlConfig *MySqlConfig
	Logger *Logger
}

type MongoConfig struct {
	MongoRepo string `json:"mongo_repo"`
}
type MySqlConfig struct {
	MysqlRepo string `json:"mysql_repo"`
}
type AppConfig struct {
	Version string `json:"version"`
}
type Logger struct {
	Level        string        // 日志打印级别
	Path         string        // 日志存放路径
	MaxAge       time.Duration // 最大存放时间
	RotationTime time.Duration // 日志分割时间
}
func InitConfig(path string){
	configVip := viper.New()
	configVip.SetConfigFile(path)

	// 读取配置
	if err := configVip.ReadInConfig(); err != nil {
		panic(err)
	}

	// 配置映射到结构体
	GlobalConfig = &Conf{}
	if err := configVip.Unmarshal(GlobalConfig); err != nil {
		panic(err)
	}
}
