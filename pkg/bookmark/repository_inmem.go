package bookmark

import (
	"strings"

	"github.com/eminetto/clean-architecture-go/pkg/entity"
)

//IRepo in memory repo
type IRepo struct {
	m map[string]*entity.Bookmark
}

//NewInmemRepository create new repository
func NewInmemRepository() *IRepo {
	var m = map[string]*entity.Bookmark{}
	return &IRepo{
		m: m,
	}
}

//Store a Bookmark
func (r *IRepo) Store(a *entity.Bookmark) (entity.ID, error) {
	r.m[a.ID.String()] = a
	return a.ID, nil
}

//Find a Bookmark
func (r *IRepo) Find(id entity.ID) (*entity.Bookmark, error) {
	if r.m[id.String()] == nil {
		return nil, entity.ErrNotFound
	}
	return r.m[id.String()], nil
}

//Search Bookmarks
func (r *IRepo) Search(query string) ([]*entity.Bookmark, error) {
	var d []*entity.Bookmark
	for _, j := range r.m {
		if strings.Contains(strings.ToLower(j.Name), query) {
			d = append(d, j)
		}
	}
	if len(d) == 0 {
		return nil, entity.ErrNotFound
	}

	return d, nil
}

//FindAll Bookmarks
func (r *IRepo) FindAll() ([]*entity.Bookmark, error) {
	var d []*entity.Bookmark
	for _, j := range r.m {
		d = append(d, j)
	}
	return d, nil
}

//Delete a Bookmark
func (r *IRepo) Delete(id entity.ID) error {
	if r.m[id.String()] == nil {
		return entity.ErrNotFound
	}
	r.m[id.String()] = nil
	return nil
}
