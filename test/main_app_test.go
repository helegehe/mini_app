package test

import (
	"context"
	"fmt"
	"github.com/helegehe/mini_app/internal/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)
var MongoRepo = "mongodb://admin:123456@127.0.0.1:27017/admin?connectTimeoutMS=10000&authSource=admin"
var MysqlRepo = "qstack:123456@tcp(10.20.74.41:3306)/matrix?charset=utf8mb4&parseTime=True&loc=Local"
func init() {
	pkg.InitDB(MongoRepo,MysqlRepo,"debug")
}

type student struct {
	Name string
	Age int
}
func TestAddMongoDb(t *testing.T) {
	stu := student{Name: "张三",Age: 20}
	coll := pkg.Mongo_Client.Database("local_db").Collection("student")
	addRt ,err := coll.InsertOne(context.TODO(),stu)
	if err != nil{
		fmt.Println(err.Error())
	}else{
		fmt.Println(addRt.InsertedID)
	}
}

func TestAddManyMongoDb(t *testing.T) {
	stu := student{Name: "张三",Age: 21}
	stu1 := student{Name: "李四",Age: 22}
	stus := []interface{}{stu,stu1}
	coll := pkg.Mongo_Client.Database("local_db").Collection("student")
	addRt ,err := coll.InsertMany(context.TODO(),stus)
	if err != nil{
		fmt.Println(err.Error())
	}else{
		fmt.Println(addRt.InsertedIDs)
	}
}

func TestQueryMongoDB(t *testing.T){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	coll := pkg.Mongo_Client.Database("local_db").Collection("student")
	cur ,err := coll.Find(ctx,bson.D{{"name",bson.D{{"$in", bson.A{"李四","王五"}}}}})
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx){
		var rt  bson.D
		err = cur.Decode(&rt)
		if err != nil{
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("result: %v\n", rt)
		fmt.Printf("result.Map(): %v\n", rt.Map()["name"])
	}
}

func TestUpdateMongo(t *testing.T) {

	coll := pkg.Mongo_Client.Database("local_db").Collection("student")
	update := bson.D{{"$set", bson.D{{"name", "王五"}, {"age", 23}}}}
	ur,err := coll.UpdateMany(context.TODO(),bson.D{{"name","张三"},{"age",20}},update)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ur.MatchedCount)
}

func TestDelMongo(t *testing.T) {
	coll :=pkg.Mongo_Client.Database("local_db").Collection("student")
	dr ,err := coll.DeleteMany(context.TODO(),bson.D{{"name","王五"}})
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(dr.DeletedCount)
}
