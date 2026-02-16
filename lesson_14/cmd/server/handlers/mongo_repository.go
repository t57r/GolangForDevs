package handlers

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

// --- Document related methods ---

func (mr *MongoRepository) PutDocument(ctx context.Context, collection string, doc map[string]any) error {
	_, err := mr.db.Collection(collection).InsertOne(ctx, bson.M(doc))
	return err
}

func (mr *MongoRepository) GetDocument(ctx context.Context, collection string, filter map[string]any) (map[string]any, error) {
	var out bson.M
	err := mr.db.Collection(collection).FindOne(ctx, bson.M(filter)).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return map[string]any{"ok": true, "found": false}, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "found": true, "document": out}, nil
}

func (mr *MongoRepository) ListDocuments(ctx context.Context, collection string, filter map[string]any, limit, skip int64) ([]map[string]any, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cur, err := mr.db.Collection(collection).Find(ctx, bson.M(filter), opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []map[string]any
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (mr *MongoRepository) DeleteDocument(ctx context.Context, collection string, filter map[string]any) (int64, error) {
	res, err := mr.db.Collection(collection).DeleteOne(ctx, bson.M(filter))
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// --- Collection related methods ---

func (mr *MongoRepository) CreateCollection(ctx context.Context, name string) error {
	return mr.db.CreateCollection(ctx, name)
}

func (mr *MongoRepository) ListCollections(ctx context.Context) ([]string, error) {
	return mr.db.ListCollectionNames(ctx, bson.M{})
}

func (mr *MongoRepository) DeleteCollection(ctx context.Context, name string) error {
	return mr.db.Collection(name).Drop(ctx)
}

// --- Index related methods ---

func (mr *MongoRepository) CreateIndex(ctx context.Context, collection string, keys map[string]int, unique bool, name string) (string, error) {
	bsonKeys := bson.D{}
	for k, v := range keys {
		if v != 1 && v != -1 {
			return "", errors.New("index direction must be 1 or -1")
		}
		bsonKeys = append(bsonKeys, bson.E{Key: k, Value: v})
	}

	model := mongo.IndexModel{
		Keys: bsonKeys,
		Options: options.Index().
			SetUnique(unique),
	}
	if name != "" {
		model.Options.SetName(name)
	}

	return mr.db.Collection(collection).Indexes().CreateOne(ctx, model)
}

func (mr *MongoRepository) DeleteIndex(ctx context.Context, collection string, name string) error {
	_, err := mr.db.Collection(collection).Indexes().DropOne(ctx, name)
	return err
}
