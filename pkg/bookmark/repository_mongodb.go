package bookmark

import (
	"fmt"
	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"github.com/juju/mgosession"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

//MongoRepository mongodb repo
type MongoRepository struct {
	pool *mgosession.Pool
	db   string
}

//NewMongoRepository create new repository
func NewMongoRepository(p *mgosession.Pool, db string) *MongoRepository {
	return &MongoRepository{
		pool: p,
		db:   db,
	}
}

//Find a bookmark
func (r *MongoRepository) Find(id entity.ID) (*entity.Bookmark, error) {
	result := entity.Bookmark{}
	session := r.pool.Session(nil)
	coll := session.DB(r.db).C("bookmark")
	err := coll.Find(bson.M{"_id": id}).One(&result)
	switch err {
	case nil:
		return &result, nil
	case mgo.ErrNotFound:
		return nil, entity.ErrNotFound
	default:
		return nil, err
	}
}

//Store a bookmark
func (r *MongoRepository) Store(b *entity.Bookmark) (entity.ID, error) {
	session := r.pool.Session(nil)
	coll := session.DB(r.db).C("bookmark")
	err := coll.Insert(b)
	if err != nil {
		return entity.ID(fmt.Sprint(0)), err
	}
	return b.ID, nil
}

//FindAll bookmarks
func (r *MongoRepository) FindAll() ([]*entity.Bookmark, error) {
	var d []*entity.Bookmark
	session := r.pool.Session(nil)
	coll := session.DB(r.db).C("bookmark")
	err := coll.Find(nil).Sort("name").All(&d)
	switch err {
	case nil:
		return d, nil
	case mgo.ErrNotFound:
		return nil, entity.ErrNotFound
	default:
		return nil, err
	}
}

//Search bookmarks
func (r *MongoRepository) Search(query string) ([]*entity.Bookmark, error) {
	var d []*entity.Bookmark
	session := r.pool.Session(nil)
	coll := session.DB(r.db).C("bookmark")
	err := coll.Find(bson.M{"name": &bson.RegEx{Pattern: query, Options: "i"}}).Limit(10).Sort("name").All(&d)
	switch err {
	case nil:
		return d, nil
	case mgo.ErrNotFound:
		return nil, entity.ErrNotFound
	default:
		return nil, err
	}
}

//Delete a bookmark
func (r *MongoRepository) Delete(id entity.ID) error {
	session := r.pool.Session(nil)
	coll := session.DB(r.db).C("bookmark")
	return coll.Remove(bson.M{"_id": id})
}
