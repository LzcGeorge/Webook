package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://root:mongoroot@localhost:27017/").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
	// 操作 client
	col := client.Database("webook").
		Collection("articles")

	id, err := col.InsertOne(ctx, Article{
		Id:       1,
		Title:    "我的标题",
		Content:  "我的内容",
		AuthorId: 123,
	})
	assert.NoError(t, err)
	fmt.Printf("插入ID: %v\n", id)

}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	// 更新时间
	Utime int64 `bson:"utime,omitempty"`
}
