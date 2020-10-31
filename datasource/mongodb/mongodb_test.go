package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

const (
	DatabaseName   = "testDB"
	CollectionName = "testColl"
)

func Test_MongoDbTest(t *testing.T) {
	printPink("=================================== begin connect =======================================")
	client, err := BuildMongoDBConnect("mongodb://192.168.9.102:31002")
	if err != nil {
		fmt.Printf("build mongodb connect fail, cause: %v", err)
		return
	}
	printPink("=================================== end connect =======================================")
	defer func() {
		printPink("=================================== begin close session =======================================")
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
		printPink("=================================== end close session =======================================")
	}()

	// collection
	collect := client.Database(DatabaseName).Collection(CollectionName)

	// insert
	insertOneResult, err := collect.InsertOne(context.Background(), buildInsertOneData())
	if err == nil {
		fmt.Printf("插入一条 result: %T %+v, err:%v\n", insertOneResult.InsertedID, insertOneResult, err)
	}

	insertManyResult, err := collect.InsertMany(context.Background(), buildInsertManyData())
	if err == nil {
		fmt.Printf("插入多条 result: %T %+v, err:%v\n", insertManyResult.InsertedIDs, insertManyResult, err)
	}

	// update TODO 这里只能用小写吗? 这也对 Go 太不友好了吧
	p := &Person{Name: "zs"}
	updateResult, err := collect.UpdateOne(context.Background(), bson.D{{"id", 1}}, bson.D{{"$set", p}})
	if err == nil {
		fmt.Printf("更新一条 result:%+v, err:%v\n", updateResult, err)
	}

	// query
	cursor, err := collect.Find(context.Background(), bson.D{{"id", bson.D{{"$lte", 2}}}})
	if err == nil {
		for cursor.Next(context.Background()) {
			tmp := new(Person)
			if err := cursor.Decode(tmp); err == nil {
				fmt.Printf("查询结果 result:%v, err:%v\n", tmp, err)
			}
		}
	}

	// delete 同查询, 条件而已

	err = client.Database(DatabaseName).Drop(context.Background())
	fmt.Printf("删除数据库 %s err:%v\n", DatabaseName, err)
}

type Person struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Birthday    *time.Time `json:"birthday"`
	Descendants []*Person  `json:"descendants"`
}

func buildInsertOneData() interface{} {
	return Person{Id: 0, Name: "user0", Birthday: NowTimePointer()}
}

func buildInsertManyData() []interface{} {
	data := make([]interface{}, 0, 5)
	for i := 1; i <= 5; i++ {
		data = append(data, Person{Id: int64(i), Name: fmt.Sprintf("user%d", i), Birthday: NowTimePointer()})
	}
	return data
}

func NowTimePointer() *time.Time {
	tmp := time.Now()
	return &tmp
}
