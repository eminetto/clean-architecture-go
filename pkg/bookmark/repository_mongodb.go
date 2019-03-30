package bookmark

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//MongoRepository mongodb repo
type MongoRepository struct {
	client *mongo.Client
	db     string
}

//NewMongoRepository create new repository
func NewMongoRepository(c *mongo.Client, db string) *MongoRepository {
	return &MongoRepository{
		client: c,
		db:     db,
	}
}

//Find a bookmark
func (r *MongoRepository) Find(id entity.ID) (*entity.Bookmark, error) {
	result := entity.Bookmark{}
	coll := r.client.Database(r.db).Collection("bookmark")
	err := coll.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	switch err {
	case nil:
		return &result, nil
	case mongo.ErrNoDocuments:
		return nil, entity.ErrNotFound
	default:
		return nil, err
	}
}

//Store a bookmark
func (r *MongoRepository) Store(b *entity.Bookmark) (entity.ID, error) {
	coll := r.client.Database(r.db).Collection("bookmark")
	_, err := coll.InsertOne(context.TODO(), b)
	if err != nil {
		return entity.NewEmptyID(), err
	}
	// b.ID = insertResult.InsertedID.(entity.ID)
	return b.ID, nil
}

//FindAll bookmarks
func (r *MongoRepository) FindAll() ([]*entity.Bookmark, error) {
	var d []*entity.Bookmark
	coll := r.client.Database(r.db).Collection("bookmark")
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"name": 1})
	cur, err := coll.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var i entity.Bookmark
		err := cur.Decode(&i)
		if err != nil {
			return nil, err
		}
		d = append(d, &i)
	}
	if len(d) == 0 {
		return nil, entity.ErrNotFound
	}
	return d, nil
}

//Search bookmarks
func (r *MongoRepository) Search(query string) ([]*entity.Bookmark, error) {
	var d []*entity.Bookmark
	coll := r.client.Database(r.db).Collection("bookmark")
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"name": 1})
	findOptions.SetLimit(10)

	cur, err := coll.Find(context.TODO(), bson.M{"name": primitive.Regex{Pattern: query, Options: "i"}})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var i entity.Bookmark
		err := cur.Decode(&i)
		if err != nil {
			return nil, err
		}
		d = append(d, &i)
	}
	if len(d) == 0 {
		return nil, entity.ErrNotFound
	}
	return d, nil
}

//Delete a bookmark
func (r *MongoRepository) Delete(id entity.ID) error {
	coll := r.client.Database(r.db).Collection("bookmark")
	_, err := coll.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
