package bookmark

import "github.com/eminetto/clean-architecture-go/pkg/entity"

//Reader interface
type Reader interface {
	Find(id entity.ID) (*entity.Bookmark, error)
	Search(query string) ([]*entity.Bookmark, error)
	FindAll() ([]*entity.Bookmark, error)
}

//Writer bookmark writer
type Writer interface {
	Store(b *entity.Bookmark) (entity.ID, error)
	Delete(id entity.ID) error
}

//Repository repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase use case interface
type UseCase interface {
	Reader
	Writer
}
