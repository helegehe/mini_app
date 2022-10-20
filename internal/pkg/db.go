package pkg

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var Mongo_Client *mongo.Client
var Mysql_Client *gorm.DB

func InitDB(mongo_url,mysql_url ,logLevel string ){
	if mongo_url != ""{
		initMongoDB(mongo_url,logLevel)
	}

	if mysql_url != ""{
		initMySql(mysql_url)
	}
}
func initMongoDB(url ,logLevel string)  {
	if url == ""{
		return
	}
	// 设置客户端连接配置
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}
	clientOptions := options.Client().ApplyURI(url)
	if logLevel == "debug"{
		clientOptions.SetMonitor(cmdMonitor)
	}

	// 连接到MongoDB
	var err error
	Mongo_Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = Mongo_Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// todo
func initMySql(url string)  {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Mysql_Client = db
	s, err := Mysql_Client.DB()
	s.SetMaxOpenConns(50)
	s.SetMaxIdleConns(10)
	s.SetConnMaxLifetime(time.Hour)
}




