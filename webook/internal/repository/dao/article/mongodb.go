package article

import (
	"context"
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBArticleDAO struct {
	col     *mongo.Collection
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func NewMongoDBArticleDAO(mongoDB *mongo.Database, node *snowflake.Node) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
		col:     mongoDB.Collection("article"),
		liveCol: mongoDB.Collection("published_article"),
		node:    node,
	}
}

func (dao *MongoDBArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	art.Id = dao.node.Generate().Int64()
	_, err := dao.col.InsertOne(ctx, &art)
	return art.Id, err
}

func (dao *MongoDBArticleDAO) UpdateById(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()

	// 过滤条件
	filter := bson.M{
		"id":        art.Id,
		"author_id": art.AuthorId,
	}

	// update 赋值
	set := bson.D{
		bson.E{
			Key: "$set",
			Value: bson.M{
				"title":   art.Title,
				"content": art.Content,
				"status":  art.Status,
				"utime":   now,
			},
		},
	}

	res, err := dao.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return 0, err
	}

	if res.ModifiedCount == 0 {
		return 0, errors.New("可能是别人的文章")
	}

	return art.Id, nil
}

func (dao *MongoDBArticleDAO) Upsert(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)

	// 写者库
	if id > 0 {
		id, err = dao.UpdateById(ctx, art)
	} else {
		id, err = dao.Insert(ctx, art)
	}

	if err != nil {
		return 0, err
	}

	// 线上库
	art.Id = id
	now := time.Now().UnixMilli()
	art.Utime = now
	filter := bson.D{
		bson.E{
			Key:   "id",
			Value: id,
		},
	}

	set := bson.D{
		bson.E{
			Key:   "$set",
			Value: art,
		},
		bson.E{
			Key: "$setOnInsert",
			Value: bson.D{
				bson.E{
					Key:   "ctime",
					Value: now,
				},
			},
		},
	}

	opt := options.Update().SetUpsert(true)
	_, err = dao.liveCol.UpdateOne(ctx, filter, set, opt)

	return id, err
}

func (dao *MongoDBArticleDAO) UpdateStatus(ctx context.Context, art Article) (int64, error) {
	filter := bson.D{
		bson.E{
			Key:   "id",
			Value: art.Id,
		},
		bson.E{
			Key:   "author_id",
			Value: art.AuthorId,
		},
	}

	set := bson.D{
		bson.E{
			Key: "$set",
			Value: bson.M{
				"status": art.Status,
				"utime":  time.Now().UnixMilli(),
			},
		},
	}

	res, err := dao.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return 0, err
	}

	if res.ModifiedCount == 0 {
		return 0, errors.New("可能是别人的文章")
	}

	// 线上库
	_, err = dao.liveCol.UpdateOne(ctx, filter, set)
	return art.Id, err
}
