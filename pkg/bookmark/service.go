package bookmark

import (
	"strings"
	"time"

	"github.com/eminetto/clean-architecture-go/pkg/entity"
)

//Service service interface
type Service struct {
	repo Repository
}

//NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

//Store an bookmark
func (s *Service) Store(b *entity.Bookmark) (entity.ID, error) {
	b.ID = entity.NewID()
	b.CreatedAt = time.Now()
	return s.repo.Store(b)
}

//Find a bookmark
func (s *Service) Find(id entity.ID) (*entity.Bookmark, error) {
	return s.repo.Find(id)
}

//Search bookmarks
func (s *Service) Search(query string) ([]*entity.Bookmark, error) {
	return s.repo.Search(strings.ToLower(query))
}

//FindAll bookmarks
func (s *Service) FindAll() ([]*entity.Bookmark, error) {
	return s.repo.FindAll()
}

//Delete a bookmark
func (s *Service) Delete(id entity.ID) error {
	b, err := s.Find(id)
	if err != nil {
		return err
	}
	if b.Favorite {
		return entity.ErrCannotBeDeleted
	}
	return s.repo.Delete(id)
}
